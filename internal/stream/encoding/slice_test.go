package encoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_writeSliceLength(t *testing.T) {
	method := func(e *encoder) func(int) error {
		return e.writeSliceLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:         "FixArray.error",
			Value:        5,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:       "FixArray.ok",
			Value:      5,
			Expected:   []byte{def.FixArray + 5},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Array16.error.def",
			Value:        32,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Array16.error.value",
			Value:        32,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Array16},
			Method:       method,
		},
		{
			Name:       "Array16.ok",
			Value:      32,
			Expected:   []byte{def.Array16, 0x00, 0x20},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Array32.error.def",
			Value:        65536,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Array32.error.value",
			Value:        65536,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Array32},
			Method:       method,
		},
		{
			Name:       "Array32.ok",
			Value:      65536,
			Expected:   []byte{def.Array32, 0x00, 0x01, 0x00, 0x00},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}

func Test_writeFixedSlice(t *testing.T) {
	method := func(e *encoder) func(reflect.Value) (bool, error) {
		return e.writeFixedSlice
	}

	t.Run("[]int", func(t *testing.T) {
		value := []int{-1, -2, -3}
		testcases := AsXXXTestCases[[]int]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{0xff, 0xfe, 0xfd},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]uint", func(t *testing.T) {
		value := []uint{1, 2, 3}
		testcases := AsXXXTestCases[[]uint]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{0x01, 0x02, 0x03},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]string", func(t *testing.T) {
		value := []string{"a", "b", "c"}
		testcases := AsXXXTestCases[[]string]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{0xa1, 0x61, 0xa1, 0x62, 0xa1, 0x63},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]float32", func(t *testing.T) {
		value := []float32{4, 5, 6}
		testcases := AsXXXTestCases[[]float32]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.Float32, 0x40, 0x80, 0x00, 0x00, def.Float32, 0x40, 0xa0, 0x00, 0x00, def.Float32, 0x40, 0xc0, 0x00, 0x00},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]float64", func(t *testing.T) {
		value := []float64{4, 5, 6}
		testcases := AsXXXTestCases[[]float64]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Float64, 0x40, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.Float64, 0x40, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.Float64, 0x40, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]bool", func(t *testing.T) {
		value := []bool{true, false, true}
		testcases := AsXXXTestCases[[]bool]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.True, def.False, def.True},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]int8", func(t *testing.T) {
		value := []int8{math.MinInt8, math.MinInt8 + 1, math.MinInt8 + 2}
		testcases := AsXXXTestCases[[]int8]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int8, 0x80,
					def.Int8, 0x81,
					def.Int8, 0x82,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]int16", func(t *testing.T) {
		value := []int16{math.MinInt16, math.MinInt16 + 1, math.MinInt16 + 2}
		testcases := AsXXXTestCases[[]int16]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int16, 0x80, 0x00,
					def.Int16, 0x80, 0x01,
					def.Int16, 0x80, 0x02,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]int32", func(t *testing.T) {
		value := []int32{math.MinInt32, math.MinInt32 + 1, math.MinInt32 + 2}
		testcases := AsXXXTestCases[[]int32]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int32, 0x80, 0x00, 0x00, 0x00,
					def.Int32, 0x80, 0x00, 0x00, 0x01,
					def.Int32, 0x80, 0x00, 0x00, 0x02,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]int64", func(t *testing.T) {
		value := []int64{math.MinInt64, math.MinInt64 + 1, math.MinInt64 + 2}
		testcases := AsXXXTestCases[[]int64]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]uint8", func(t *testing.T) {
		value := []uint8{math.MaxUint8, math.MaxUint8 - 1, math.MaxUint8 - 2}
		testcases := AsXXXTestCases[[]uint8]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint8, 0xff,
					def.Uint8, 0xfe,
					def.Uint8, 0xfd,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]uint16", func(t *testing.T) {
		value := []uint16{math.MaxUint16, math.MaxUint16 - 1, math.MaxUint16 - 2}
		testcases := AsXXXTestCases[[]uint16]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint16, 0xff, 0xff,
					def.Uint16, 0xff, 0xfe,
					def.Uint16, 0xff, 0xfd,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]uint32", func(t *testing.T) {
		value := []uint32{math.MaxUint32, math.MaxUint32 - 1, math.MaxUint32 - 2}
		testcases := AsXXXTestCases[[]uint32]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint32, 0xff, 0xff, 0xff, 0xff,
					def.Uint32, 0xff, 0xff, 0xff, 0xfe,
					def.Uint32, 0xff, 0xff, 0xff, 0xfd,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("[]uint64", func(t *testing.T) {
		value := []uint64{math.MaxUint64, math.MaxUint64 - 1, math.MaxUint64 - 2}
		testcases := AsXXXTestCases[[]uint64]{
			{
				Name:           "error",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
					def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfd,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
}
