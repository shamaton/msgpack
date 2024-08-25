package encoding

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/internal/common"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

const dummyByte = 0xc1

type TestWriter struct {
	WrittenBytes []byte
}

var _ io.Writer = (*TestWriter)(nil)

var ErrTestWriter = errors.New("expected written error")

func (w *TestWriter) Write(p []byte) (n int, err error) {
	w.WrittenBytes = append(w.WrittenBytes, p...)
	if bytes.Contains(p, []byte{dummyByte}) {
		return 0, ErrTestWriter
	}
	return len(p), nil
}

func NewTestWriter() *TestWriter {
	return &TestWriter{}
}

type AsXXXTestCase[T any] struct {
	Name            string
	Value           T
	Expected        []byte
	Contains        []byte
	BufferSize      int
	PreWriteSize    int
	Error           error
	AsArray         bool
	Method          func(*encoder) func(T) error
	MethodForFixed  func(*encoder) func(reflect.Value) (bool, error)
	MethodForStruct func(*encoder) func(reflect.Value) error
}

type AsXXXTestCases[T any] []AsXXXTestCase[T]

func (tcs AsXXXTestCases[T]) Run(t *testing.T) {
	for _, tc := range tcs {
		tc.Run(t)
	}
}

func (tc *AsXXXTestCase[T]) Run(t *testing.T) {
	t.Helper()

	if tc.Method == nil && tc.MethodForFixed == nil && tc.MethodForStruct == nil {
		t.Fatal("must set either Method or MethodForFixed or MethodForStruct")
	}

	method := func(e *encoder) error {
		if tc.Method != nil {
			return tc.Method(e)(tc.Value)
		}
		if tc.MethodForFixed != nil {
			_, err := tc.MethodForFixed(e)(reflect.ValueOf(tc.Value))
			return err
		}
		if tc.MethodForStruct != nil {
			return tc.MethodForStruct(e)(reflect.ValueOf(tc.Value))
		}
		panic("unreachable")
	}

	t.Run(tc.Name, func(t *testing.T) {
		w := NewTestWriter()
		e := encoder{
			w:       w,
			buf:     common.GetBuffer(),
			Common:  common.Common{},
			asArray: tc.AsArray,
		}

		if tc.BufferSize < tc.PreWriteSize {
			t.Fatal("buffer size must be greater than pre write size")
		}

		e.buf.Data = make([]byte, tc.BufferSize)
		if tc.PreWriteSize > 0 {
			for i := 0; i < tc.PreWriteSize; i++ {
				_ = e.buf.Write(e.w, dummyByte)
			}
		}

		err := method(&e)
		_ = e.buf.Flush(w)
		common.PutBuffer(e.buf)

		if tc.PreWriteSize > 0 {
			tu.IsError(t, err, ErrTestWriter)
			if !bytes.Contains(w.WrittenBytes, tc.Contains) {
				t.Fatalf("[% 02x] does not contain in [% 02x]", tc.Contains, w.WrittenBytes)
			}
			return
		}
		if tc.Error != nil {
			tu.IsError(t, err, tc.Error)
			return
		}

		tu.NoError(t, err)
		tu.EqualSlice(t, w.WrittenBytes, tc.Expected)
	})
}

func TestEncode(t *testing.T) {
	v := 1
	vv := &v

	w := NewTestWriter()
	err := Encode(w, &vv, false)
	tu.NoError(t, err)

	tu.EqualSlice(t, w.WrittenBytes, []byte{def.PositiveFixIntMin + 1})
}

func Test_create(t *testing.T) {
	method := func(e *encoder) func(reflect.Value) error {
		return e.create
	}

	t.Run("uint8", func(t *testing.T) {
		value := uint8(1)
		testcases := AsXXXTestCases[uint8]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.PositiveFixIntMin + 1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("uint16", func(t *testing.T) {
		value := uint16(1)
		testcases := AsXXXTestCases[uint16]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.PositiveFixIntMin + 1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("uint32", func(t *testing.T) {
		value := uint32(1)
		testcases := AsXXXTestCases[uint32]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.PositiveFixIntMin + 1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("uint64", func(t *testing.T) {
		value := uint64(1)
		testcases := AsXXXTestCases[uint64]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.PositiveFixIntMin + 1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("uint", func(t *testing.T) {
		value := uint(1)
		testcases := AsXXXTestCases[uint]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.PositiveFixIntMin + 1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("int8", func(t *testing.T) {
		value := int8(-1)
		testcases := AsXXXTestCases[int8]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{0xff},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("int16", func(t *testing.T) {
		value := int16(-1)
		testcases := AsXXXTestCases[int16]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{0xff},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("int32", func(t *testing.T) {
		value := int32(-1)
		testcases := AsXXXTestCases[int32]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{0xff},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("int64", func(t *testing.T) {
		value := int64(-1)
		testcases := AsXXXTestCases[int64]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{0xff},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("int", func(t *testing.T) {
		value := int(-1)
		testcases := AsXXXTestCases[int]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{0xff},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("float32", func(t *testing.T) {
		value := float32(1)
		testcases := AsXXXTestCases[float32]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.Float32, 0x3f, 0x80, 0x00, 0x00},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("float64", func(t *testing.T) {
		value := float64(1)
		testcases := AsXXXTestCases[float64]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.Float64, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("bool", func(t *testing.T) {
		value := true
		testcases := AsXXXTestCases[bool]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.True},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("string", func(t *testing.T) {
		value := "a"
		testcases := AsXXXTestCases[string]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.FixStr + 1, 0x61},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("complex64", func(t *testing.T) {
		value := complex64(complex(1, 2))
		testcases := AsXXXTestCases[complex64]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Fixext8, byte(def.ComplexTypeCode()),
					0x3f, 0x80, 0x00, 0x00,
					0x40, 0x00, 0x00, 0x00,
				},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("complex128", func(t *testing.T) {
		value := complex128(complex(1, 2))
		testcases := AsXXXTestCases[complex128]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Fixext16, byte(def.ComplexTypeCode()),
					0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})

	type st struct {
		A int
	}

	t.Run("slice", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var value []int
			testcases := AsXXXTestCases[[]int]{
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.Nil},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("bin", func(t *testing.T) {
			value := []byte{1, 2, 3}
			testcases := AsXXXTestCases[[]byte]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      3,
					PreWriteSize:    1,
					Contains:        []byte{def.Bin8, 0x03},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.Bin8, 0x03, 0x01, 0x02, 0x03},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("fixed", func(t *testing.T) {
			value := []int{1, 2, 3}
			testcases := AsXXXTestCases[[]int]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixArray + 3},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixArray + 3, 0x01, 0x02, 0x03},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("slice_slice", func(t *testing.T) {
			value := [][]int{{1}}
			testcases := AsXXXTestCases[[][]int]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixArray + 1},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixArray + 1, def.FixArray + 1, 0x01},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("struct", func(t *testing.T) {
			value := []st{{1}}
			testcases := AsXXXTestCases[[]st]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixArray + 1},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixArray + 1, def.FixMap + 1, def.FixStr + 1, 0x41, 0x01},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("bin", func(t *testing.T) {
			value := [3]byte{1, 2, 3}
			testcases := AsXXXTestCases[[3]byte]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      3,
					PreWriteSize:    1,
					Contains:        []byte{def.Bin8, 0x03},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.Bin8, 0x03, 0x01, 0x02, 0x03},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("fixed", func(t *testing.T) {
			value := [3]int{1, 2, 3}
			testcases := AsXXXTestCases[[3]int]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixArray + 3},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixArray + 3, 0x01, 0x02, 0x03},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("slice_slice", func(t *testing.T) {
			value := [1][]int{{1}}
			testcases := AsXXXTestCases[[1][]int]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixArray + 1},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixArray + 1, def.FixArray + 1, 0x01},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("struct", func(t *testing.T) {
			value := [1]st{{1}}
			testcases := AsXXXTestCases[[1]st]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixArray + 1},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixArray + 1, def.FixMap + 1, def.FixStr + 1, 0x41, 0x01},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var value map[string]int
			testcases := AsXXXTestCases[map[string]int]{
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.Nil},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("fixed", func(t *testing.T) {
			value := map[string]int{"a": 1}
			testcases := AsXXXTestCases[map[string]int]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixMap + 1},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixMap + 1, def.FixStr + 1, 'a', 1},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		t.Run("struct", func(t *testing.T) {
			value := map[string]st{"a": {1}}
			testcases := AsXXXTestCases[map[string]st]{
				{
					Name:            "error.length",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "error.key",
					Value:           value,
					BufferSize:      2,
					PreWriteSize:    1,
					Contains:        []byte{def.FixMap + 1},
					MethodForStruct: method,
				},
				{
					Name:            "error.value",
					Value:           value,
					BufferSize:      4,
					PreWriteSize:    1,
					Contains:        []byte{def.FixMap + 1, def.FixStr + 1, 'a'},
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{def.FixMap + 1, def.FixStr + 1, 'a', def.FixMap + 1, def.FixStr + 1, 'A', 1},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
	})

	t.Run("struct", func(t *testing.T) {
		value := st{1}
		testcases := AsXXXTestCases[st]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.FixMap + 1, def.FixStr + 1, 'A', 1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("pointer", func(t *testing.T) {
		var v *int
		vv := &v
		value := &vv
		t.Run("nil", func(t *testing.T) {
			testcases := AsXXXTestCases[***int]{
				{
					Name:            "nil",
					Value:           value,
					Expected:        []byte{def.Nil},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
		u := 1
		v = &u
		vv = &v
		value = &vv
		t.Run("int", func(t *testing.T) {
			testcases := AsXXXTestCases[***int]{
				{
					Name:            "error",
					Value:           value,
					BufferSize:      1,
					PreWriteSize:    1,
					MethodForStruct: method,
				},
				{
					Name:            "ok",
					Value:           value,
					Expected:        []byte{1},
					BufferSize:      1,
					MethodForStruct: method,
				},
			}
			testcases.Run(t)
		})
	})
	t.Run("any", func(t *testing.T) {
		value := any(1)
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{1},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("invalid", func(t *testing.T) {
		var value error
		testcases := AsXXXTestCases[error]{
			{
				Name:            "ok",
				Value:           value,
				Expected:        []byte{def.Nil},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("invalid", func(t *testing.T) {
		var value func()
		testcases := AsXXXTestCases[func()]{
			{
				Name:            "ok",
				Value:           value,
				BufferSize:      1,
				Error:           def.ErrUnsupportedType,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
}
