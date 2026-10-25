package msgpack_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3"
	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

const extDestCode int8 = 98

type extDestValue int8

type extDestEncoder struct{ ext.EncoderCommon }

func (*extDestEncoder) Code() int8                              { return extDestCode }
func (*extDestEncoder) Type() reflect.Type                      { return reflect.TypeOf(extDestValue(0)) }
func (*extDestEncoder) CalcByteSize(reflect.Value) (int, error) { return 3, nil }
func (e *extDestEncoder) WriteToBytes(rv reflect.Value, offset int, data *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, data)
	offset = e.SetByte1Int(int(extDestCode), offset, data)
	return e.SetByte1Int64(rv.Int(), offset, data)
}

type extDestDecoder struct{ ext.DecoderCommon }

func (*extDestDecoder) Code() int8 { return extDestCode }
func (*extDestDecoder) IsType(offset int, d *[]byte) bool {
	return (*d)[offset] == def.Fixext1 && int8((*d)[offset+1]) == extDestCode
}
func (*extDestDecoder) AsValue(offset int, _ reflect.Kind, d *[]byte) (interface{}, int, error) {
	return extDestValue(int8((*d)[offset+2])), offset + 3, nil
}

type extDestStreamEncoder struct{}

func (*extDestStreamEncoder) Code() int8         { return extDestCode }
func (*extDestStreamEncoder) Type() reflect.Type { return reflect.TypeOf(extDestValue(0)) }
func (*extDestStreamEncoder) Write(w ext.StreamWriter, rv reflect.Value) error {
	if err := w.WriteByte1Int(def.Fixext1); err != nil {
		return err
	}
	if err := w.WriteByte1Int(int(extDestCode)); err != nil {
		return err
	}
	return w.WriteByte1Int64(rv.Int())
}

type extDestStreamDecoder struct{}

func (*extDestStreamDecoder) Code() int8 { return extDestCode }
func (*extDestStreamDecoder) IsType(code byte, innerType int8, _ int) bool {
	return code == def.Fixext1 && innerType == extDestCode
}
func (*extDestStreamDecoder) ToValue(_ byte, data []byte, _ reflect.Kind) (interface{}, error) {
	switch data[0] {
	case extDestPayloadNil:
		return nil, nil
	case extDestPayloadPtr:
		v := extDestValue(7)
		return &v, nil
	}
	return extDestValue(int8(data[0])), nil
}

// Payloads that make extDestStreamDecoder return values other than
// extDestValue, to exercise how they are stored.
const (
	extDestPayloadNil byte = 0xf0
	extDestPayloadPtr byte = 0xf1
)

// Registering a custom stream ext decoder used to make the stream decoder read
// an ext frame twice for pointer and interface destinations, failing with EOF.
// The buffer decoder is the reference for every destination.
func TestDecodeExtIntoPointerAndInterface(t *testing.T) {
	enc, dec := &extDestEncoder{}, &extDestDecoder{}
	if err := msgpack.AddExtCoder(enc, dec); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = msgpack.RemoveExtCoder(enc, dec) })
	senc, sdec := &extDestStreamEncoder{}, &extDestStreamDecoder{}
	if err := msgpack.AddExtStreamCoder(senc, sdec); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = msgpack.RemoveExtStreamCoder(senc, sdec) })

	custom := []byte{def.Fixext1, byte(extDestCode), 0x07}
	tm, err := msgpack.Marshal(time.Unix(1700000000, 123).UTC())
	if err != nil {
		t.Fatal(err)
	}

	fixmap1 := func(key string, v []byte) []byte {
		return append([]byte{0x81, 0xa0 | byte(len(key))}, append([]byte(key), v...)...)
	}
	fixarray1 := func(v []byte) []byte { return append([]byte{0x91}, v...) }

	c64, err := msgpack.Marshal(complex64(1 + 2i))
	if err != nil {
		t.Fatal(err)
	}
	c128, err := msgpack.Marshal(complex128(3 + 4i))
	if err != nil {
		t.Fatal(err)
	}

	type customPtrHolder struct{ V *extDestValue }
	type complexHolder struct {
		C complex128
		P *complex64
	}
	type timePtrHolder struct{ T *time.Time }

	cases := []struct {
		name string
		data []byte
		newV func() interface{}
	}{
		{"custom/value", custom, func() interface{} { return new(extDestValue) }},
		{"custom/pointer", custom, func() interface{} { return new(*extDestValue) }},
		{"custom/interface", custom, func() interface{} { return new(interface{}) }},
		{"custom/mapStringIface", fixmap1("a", custom), func() interface{} { return new(map[string]interface{}) }},
		{"custom/mapIfaceIface", fixmap1("a", custom), func() interface{} { return new(map[interface{}]interface{}) }},
		{"custom/sliceIface", fixarray1(custom), func() interface{} { return new([]interface{}) }},
		{"custom/structPtrField", fixmap1("V", custom), func() interface{} { return new(customPtrHolder) }},
		{"time/value", tm, func() interface{} { return new(time.Time) }},
		{"time/pointer", tm, func() interface{} { return new(*time.Time) }},
		{"time/interface", tm, func() interface{} { return new(interface{}) }},
		{"time/mapStringIface", fixmap1("t", tm), func() interface{} { return new(map[string]interface{}) }},
		{"time/sliceIface", fixarray1(tm), func() interface{} { return new([]interface{}) }},
		{"time/structPtrField", fixmap1("T", tm), func() interface{} { return new(timePtrHolder) }},
		{"complex64/value", c64, func() interface{} { return new(complex64) }},
		{"complex64/pointer", c64, func() interface{} { return new(*complex64) }},
		{"complex128/value", c128, func() interface{} { return new(complex128) }},
		{"complex128/sliceOfPointer", fixarray1(c128), func() interface{} { return new([]*complex128) }},
		{"complex/structFields", append(append([]byte{0x82, 0xa1, 'C'}, c128...), append([]byte{0xa1, 'P'}, c64...)...), func() interface{} { return new(complexHolder) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			want := c.newV()
			if err := msgpack.Unmarshal(c.data, want); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			got := c.newV()
			err := mustNotPanic(t, func() error { return msgpack.UnmarshalRead(bytes.NewReader(c.data), got) })
			if err != nil {
				t.Fatalf("UnmarshalRead: %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("UnmarshalRead = %#v, Unmarshal = %#v", reflect.ValueOf(got).Elem(), reflect.ValueOf(want).Elem())
			}
		})
	}

	t.Run("custom/interfaceHoldingPointer", func(t *testing.T) {
		var dst extDestValue
		var v interface{} = &dst
		if err := msgpack.UnmarshalRead(bytes.NewReader(custom), &v); err != nil {
			t.Fatal(err)
		}
		if dst != 7 {
			t.Fatalf("got %d, want 7", dst)
		}
	})

	t.Run("custom/interfaceHoldingNilPointer", func(t *testing.T) {
		var v interface{} = (*extDestValue)(nil)
		if err := msgpack.UnmarshalRead(bytes.NewReader(custom), &v); err != nil {
			t.Fatal(err)
		}
		if v != extDestValue(7) {
			t.Fatalf("got %#v, want extDestValue(7)", v)
		}
	})

	mismatched := map[string]func() interface{}{
		"string":          func() interface{} { return new(string) },
		"pointer":         func() interface{} { return new(*string) },
		"nonEmptyIface":   func() interface{} { return new(fmt.Stringer) },
		"pointerInStruct": func() interface{} { return new(struct{ V *string }) },
	}
	for name, newV := range mismatched {
		t.Run("custom/mismatched/"+name, func(t *testing.T) {
			data := custom
			if name == "pointerInStruct" {
				data = fixmap1("V", custom)
			}
			for _, df := range decodeFuncs {
				// Non-empty interface destinations are not supported by the
				// kind-based decoding of either decoder; only the stream ext
				// path is checked here.
				if name == "nonEmptyIface" && df.name != "UnmarshalRead" {
					continue
				}
				err := mustNotPanic(t, func() error { return df.fn(data, newV()) })
				if err == nil || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("%s: unexpected error: %v", df.name, err)
				}
			}
		})
	}

	nilValue := []byte{def.Fixext1, byte(extDestCode), extDestPayloadNil}
	t.Run("custom/nilValue/interface", func(t *testing.T) {
		var v interface{}
		if err := msgpack.UnmarshalRead(bytes.NewReader(nilValue), &v); err != nil {
			t.Fatal(err)
		}
		if v != nil {
			t.Fatalf("got %#v, want nil", v)
		}
	})
	t.Run("custom/nilValue/value", func(t *testing.T) {
		var v extDestValue
		err := mustNotPanic(t, func() error { return msgpack.UnmarshalRead(bytes.NewReader(nilValue), &v) })
		if err == nil || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// The pointer held by the interface is not addressable, so it can not be
	// replaced by the decoded pointer.
	t.Run("custom/pointerValue/interfaceHoldingPointer", func(t *testing.T) {
		data := []byte{def.Fixext1, byte(extDestCode), extDestPayloadPtr}
		var dst extDestValue
		var v interface{} = &dst
		err := mustNotPanic(t, func() error { return msgpack.UnmarshalRead(bytes.NewReader(data), &v) })
		if err == nil || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
