package encoding

import (
	"reflect"
	"sync"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

type issue55Status int
type issue55OtherStatus int
type issue55Statuses []issue55Status
type issue55Bytes []byte

type issue55Holder struct {
	Status issue55Status
	Value  any
}

type issue55ExtEncoder struct {
	ext.EncoderCommon
	typ       reflect.Type
	code      int8
	marker    byte
	sawNil    bool
	writeCall int
}

func (e *issue55ExtEncoder) Code() int8         { return e.code }
func (e *issue55ExtEncoder) Type() reflect.Type { return e.typ }
func (e *issue55ExtEncoder) CalcByteSize(value reflect.Value) (int, error) {
	if isNilValue(value) {
		e.sawNil = true
	}
	return 3, nil
}
func (e *issue55ExtEncoder) WriteToBytes(value reflect.Value, offset int, data *[]byte) int {
	if isNilValue(value) {
		e.sawNil = true
	}
	e.writeCall++
	offset = e.SetByte1Int(def.Fixext1, offset, data)
	offset = e.SetByte1Int(int(e.code), offset, data)
	return e.SetByte1Uint64(uint64(e.marker), offset, data)
}

func isNilValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func issue55ExtBytes(code int8, marker byte) []byte {
	return []byte{def.Fixext1, byte(code), marker}
}

func addIssue55Encoder(t *testing.T, typ reflect.Type, code int8, marker byte) *issue55ExtEncoder {
	t.Helper()
	encoder := &issue55ExtEncoder{typ: typ, code: code, marker: marker}
	AddExtEncoder(encoder)
	t.Cleanup(func() { RemoveExtEncoder(encoder) })
	return encoder
}

func TestDefinedTypeExtEncoderSelection(t *testing.T) {
	statusType := reflect.TypeOf(issue55Status(0))
	statusEncoder := addIssue55Encoder(t, statusType, 21, 0xa1)
	extStatus := issue55ExtBytes(statusEncoder.code, statusEncoder.marker)

	t.Run("root", func(t *testing.T) {
		encoded, err := Encode(issue55Status(1), false)
		tu.NoError(t, err)
		tu.EqualSlice(t, encoded, extStatus)
	})

	t.Run("struct field", func(t *testing.T) {
		encoded, err := Encode(issue55Holder{Status: 1}, true)
		tu.NoError(t, err)
		expected := append([]byte{def.FixArray + 2}, extStatus...)
		expected = append(expected, def.Nil)
		tu.EqualSlice(t, encoded, expected)
	})

	t.Run("interface field", func(t *testing.T) {
		encoded, err := Encode(issue55Holder{Value: issue55Status(1)}, true)
		tu.NoError(t, err)
		expected := append([]byte{def.FixArray + 2}, extStatus...)
		expected = append(expected, extStatus...)
		tu.EqualSlice(t, encoded, expected)
	})

	t.Run("map key and value", func(t *testing.T) {
		encoded, err := Encode(map[issue55Status]issue55Status{1: 2}, false)
		tu.NoError(t, err)
		expected := append([]byte{def.FixMap + 1}, extStatus...)
		expected = append(expected, extStatus...)
		tu.EqualSlice(t, encoded, expected)
	})

	t.Run("slice elements", func(t *testing.T) {
		encoded, err := Encode([]issue55Status{1, 2}, false)
		tu.NoError(t, err)
		expected := append([]byte{def.FixArray + 2}, extStatus...)
		expected = append(expected, extStatus...)
		tu.EqualSlice(t, encoded, expected)
	})
}

func TestDefinedTypeWholeContainerExtTakesPrecedence(t *testing.T) {
	statusEncoder := addIssue55Encoder(t, reflect.TypeOf(issue55Status(0)), 22, 0xb1)
	statusesEncoder := addIssue55Encoder(t, reflect.TypeOf(issue55Statuses{}), 23, 0xb2)

	encoded, err := Encode(issue55Statuses{1, 2}, false)
	tu.NoError(t, err)
	tu.EqualSlice(t, encoded, issue55ExtBytes(statusesEncoder.code, statusesEncoder.marker))
	tu.Equal(t, statusEncoder.writeCall, 0)
}

func TestDefinedTypeExtBypassesContainerFastPaths(t *testing.T) {
	t.Run("defined byte slice", func(t *testing.T) {
		encoder := addIssue55Encoder(t, reflect.TypeOf(issue55Bytes{}), 24, 0xc1)
		encoded, err := Encode(issue55Bytes{1, 2}, false)
		tu.NoError(t, err)
		tu.EqualSlice(t, encoded, issue55ExtBytes(encoder.code, encoder.marker))
	})

	t.Run("byte elements", func(t *testing.T) {
		encoder := addIssue55Encoder(t, reflect.TypeOf(byte(0)), 25, 0xc2)
		extByte := issue55ExtBytes(encoder.code, encoder.marker)
		for _, value := range []any{[]byte{1, 2}, [2]byte{1, 2}} {
			encoded, err := Encode(value, false)
			tu.NoError(t, err)
			expected := append([]byte{def.FixArray + 2}, extByte...)
			expected = append(expected, extByte...)
			tu.EqualSlice(t, encoded, expected)
		}
	})

	t.Run("fixed slice", func(t *testing.T) {
		encoder := addIssue55Encoder(t, reflect.TypeOf(int(0)), 26, 0xc3)
		extInt := issue55ExtBytes(encoder.code, encoder.marker)
		encoded, err := Encode([]int{1, 2}, false)
		tu.NoError(t, err)
		expected := append([]byte{def.FixArray + 2}, extInt...)
		expected = append(expected, extInt...)
		tu.EqualSlice(t, encoded, expected)
	})

	t.Run("fixed map", func(t *testing.T) {
		stringEncoder := addIssue55Encoder(t, reflect.TypeOf(""), 27, 0xc4)
		intEncoder := addIssue55Encoder(t, reflect.TypeOf(int(0)), 28, 0xc5)
		encoded, err := Encode(map[string]int{"key": 1}, false)
		tu.NoError(t, err)
		expected := append([]byte{def.FixMap + 1}, issue55ExtBytes(stringEncoder.code, stringEncoder.marker)...)
		expected = append(expected, issue55ExtBytes(intEncoder.code, intEncoder.marker)...)
		tu.EqualSlice(t, encoded, expected)
	})
}

func TestDefinedTypeExtNilAndRemovalSemantics(t *testing.T) {
	t.Run("nil registered values", func(t *testing.T) {
		sliceEncoder := addIssue55Encoder(t, reflect.TypeOf(issue55Statuses(nil)), 29, 0xd1)
		pointerEncoder := addIssue55Encoder(t, reflect.TypeOf((*issue55Status)(nil)), 30, 0xd2)

		encoded, err := Encode(issue55Statuses(nil), false)
		tu.NoError(t, err)
		tu.EqualSlice(t, encoded, issue55ExtBytes(sliceEncoder.code, sliceEncoder.marker))
		tu.Equal(t, sliceEncoder.sawNil, true)

		encoded, err = Encode((*issue55Status)(nil), false)
		tu.NoError(t, err)
		tu.EqualSlice(t, encoded, issue55ExtBytes(pointerEncoder.code, pointerEncoder.marker))
		tu.Equal(t, pointerEncoder.sawNil, true)
	})

	t.Run("remove restores built-in", func(t *testing.T) {
		encoder := &issue55ExtEncoder{typ: reflect.TypeOf(issue55Status(0)), code: 31, marker: 0xd3}
		AddExtEncoder(encoder)
		RemoveExtEncoder(encoder)
		encoded, err := Encode(issue55Status(1), false)
		tu.NoError(t, err)
		tu.EqualSlice(t, encoded, []byte{1})
	})

	t.Run("first registration wins", func(t *testing.T) {
		first := addIssue55Encoder(t, reflect.TypeOf(issue55Status(0)), 32, 0xd4)
		second := &issue55ExtEncoder{typ: first.typ, code: 33, marker: 0xd5}
		AddExtEncoder(second)
		encoded, err := Encode(issue55Status(1), false)
		tu.NoError(t, err)
		tu.EqualSlice(t, encoded, issue55ExtBytes(first.code, first.marker))
	})
}

func TestExtEncoderRegistryCopyOnWrite(t *testing.T) {
	typ := reflect.TypeOf(issue55Status(0))
	before := currentExtEncoderRegistry.Load()
	encoder := &issue55ExtEncoder{typ: typ, code: 34, marker: 0xe1}

	AddExtEncoder(encoder)
	afterAdd := currentExtEncoderRegistry.Load()
	if before == afterAdd {
		t.Fatal("registration must publish a new registry snapshot")
	}
	if before.byKind[reflect.Int] != nil {
		t.Fatal("published old snapshot was mutated")
	}
	if afterAdd.byKind[reflect.Int][typ] != encoder {
		t.Fatal("new snapshot does not contain encoder")
	}
	if afterAdd.customCount != before.customCount+1 || !afterAdd.hasCustom[reflect.Int] {
		t.Fatal("new snapshot does not enable the custom kind fast-path flag")
	}

	RemoveExtEncoder(encoder)
	afterRemove := currentExtEncoderRegistry.Load()
	if afterAdd.byKind[reflect.Int][typ] != encoder {
		t.Fatal("published added snapshot was mutated")
	}
	if afterRemove.byKind[reflect.Int] != nil {
		t.Fatal("last removal must restore nil kind bucket")
	}
	if afterRemove.customCount != before.customCount || afterRemove.hasCustom[reflect.Int] {
		t.Fatal("last removal must restore the no-custom fast path")
	}
}

func TestExtEncoderRegistryKeepsKindFlagAfterPartialRemoval(t *testing.T) {
	first := &issue55ExtEncoder{typ: reflect.TypeOf(issue55Status(0)), code: 37, marker: 0xf3}
	second := &issue55ExtEncoder{typ: reflect.TypeOf(issue55OtherStatus(0)), code: 38, marker: 0xf4}
	AddExtEncoder(first)
	AddExtEncoder(second)
	t.Cleanup(func() {
		RemoveExtEncoder(first)
		RemoveExtEncoder(second)
	})

	RemoveExtEncoder(first)
	registry := currentExtEncoderRegistry.Load()
	if !registry.hasCustom[reflect.Int] {
		t.Fatal("partial removal must keep the custom kind flag")
	}
	if registry.byKind[reflect.Int][second.typ] != second {
		t.Fatal("partial removal lost the remaining encoder")
	}
	encoded, err := Encode(issue55OtherStatus(1), false)
	tu.NoError(t, err)
	tu.EqualSlice(t, encoded, issue55ExtBytes(second.code, second.marker))
}

type issue55BlockingEncoder struct {
	*issue55ExtEncoder
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (e *issue55BlockingEncoder) CalcByteSize(value reflect.Value) (int, error) {
	e.once.Do(func() {
		close(e.started)
		<-e.release
	})
	return e.issue55ExtEncoder.CalcByteSize(value)
}

func TestEncodeUsesSingleExtRegistrySnapshot(t *testing.T) {
	typ := reflect.TypeOf(issue55Status(0))
	first := &issue55BlockingEncoder{
		issue55ExtEncoder: &issue55ExtEncoder{typ: typ, code: 35, marker: 0xf1},
		started:           make(chan struct{}),
		release:           make(chan struct{}),
	}
	second := &issue55ExtEncoder{typ: typ, code: 36, marker: 0xf2}
	AddExtEncoder(first)

	type result struct {
		data []byte
		err  error
	}
	resultCh := make(chan result, 1)
	go func() {
		data, err := Encode(issue55Status(1), false)
		resultCh <- result{data: data, err: err}
	}()

	<-first.started
	RemoveExtEncoder(first)
	AddExtEncoder(second)
	close(first.release)

	encoded := <-resultCh
	tu.NoError(t, encoded.err)
	tu.EqualSlice(t, encoded.data, issue55ExtBytes(first.code, first.marker))

	next, err := Encode(issue55Status(1), false)
	tu.NoError(t, err)
	tu.EqualSlice(t, next, issue55ExtBytes(second.code, second.marker))
	RemoveExtEncoder(second)
}
