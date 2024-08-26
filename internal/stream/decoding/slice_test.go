package decoding

import (
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_sliceLength(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (int, error) {
		return d.sliceLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:             "FixArray",
			Code:             def.FixArray + 3,
			Expected:         3,
			MethodAsWithCode: method,
		},
		{
			Name:             "Array16.error",
			Code:             def.Array16,
			Data:             []byte{},
			ReadCount:        0,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Array16.ok",
			Code:             def.Array16,
			Data:             []byte{0xff, 0xff},
			Expected:         math.MaxUint16,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Array32.error",
			Code:             def.Array32,
			Data:             []byte{},
			ReadCount:        0,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Array32.ok",
			Code:             def.Array32,
			Data:             []byte{0xff, 0xff, 0xff, 0xff},
			Expected:         math.MaxUint32,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Nil,
			Error:            def.ErrCanNotDecode,
			MethodAsWithCode: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asFixedSlice_Int(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 3},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]int)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int{3})
}

func Test_asFixedSlice_Int8(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 4},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]int8)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int8{4})
}

func Test_asFixedSlice_Int16(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 5},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]int16)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int16{5})
}

func Test_asFixedSlice_Int32(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 6},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]int32)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int32{6})
}

func Test_asFixedSlice_Int64(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 7},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]int64)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int64{7})
}

func Test_asFixedSlice_Uint(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 5},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]uint)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint{5})
}

func Test_asFixedSlice_Uint8(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 6},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]uint8)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint8{6})
}

func Test_asFixedSlice_Uint16(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 7},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]uint16)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint16{7})
}

func Test_asFixedSlice_Uint32(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 8},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]uint32)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint32{8})
}

func Test_asFixedSlice_Uint64(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 9},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]uint64)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint64{9})
}

func Test_asFixedSlice_Float32(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.Float32, 63, 128, 0, 0},
				Expected:       true,
				ReadCount:      2,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]float32)
	run(t, v1)
	tu.EqualSlice(t, *v1, []float32{1})
}

func Test_asFixedSlice_Float64(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.Float64, 63, 240, 0, 0, 0, 0, 0, 0},
				Expected:       true,
				ReadCount:      2,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]float64)
	run(t, v1)
	tu.EqualSlice(t, *v1, []float64{1})
}

func Test_asFixedSlice_String(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 2)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.FixStr + 1, 'a', def.FixStr + 1, 'b'},
				Expected:       true,
				ReadCount:      4,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]string)
	run(t, v1)
	tu.EqualSlice(t, *v1, []string{"a", "b"})
}

func Test_asFixedSlice_Bool(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          io.EOF,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.True},
				Expected:       true,
				ReadCount:      1,
				MethodAsCustom: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	v1 := new([]bool)
	run(t, v1)
	tu.EqualSlice(t, *v1, []bool{true})
}
