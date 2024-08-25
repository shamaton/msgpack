package decoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_sliceLength(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (int, int, error) {
		return d.sliceLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "FixArray",
			Data:     []byte{def.FixArray + 3},
			Expected: 3,
			MethodAs: method,
		},
		{
			Name:     "Array16.error",
			Data:     []byte{def.Array16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array16.ok",
			Data:     []byte{def.Array16, 0xff, 0xff},
			Expected: math.MaxUint16,
			MethodAs: method,
		},
		{
			Name:     "Array32.error",
			Data:     []byte{def.Array32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array32.ok",
			Data:     []byte{def.Array32, 0xff, 0xff, 0xff, 0xff},
			Expected: math.MaxUint32,
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{def.Nil},
			Error:    def.ErrCanNotDecode,
			MethodAs: method,
		},
	}
	testcases.Run(t)
}

func Test_asFixedSlice_Int(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 3},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]int)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int{3})
}

func Test_asFixedSlice_Int8(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 4},
				Expected:       true,
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
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 5},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]int16)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int16{5})
}

func Test_asFixedSlice_Int32(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 6},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]int32)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int32{6})
}

func Test_asFixedSlice_Int64(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 7},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]int64)
	run(t, v1)
	tu.EqualSlice(t, *v1, []int64{7})
}

func Test_asFixedSlice_Uint(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 5},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]uint)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint{5})
}

func Test_asFixedSlice_Uint8(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 6},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]uint8)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint8{6})
}

func Test_asFixedSlice_Uint16(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 7},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]uint16)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint16{7})
}

func Test_asFixedSlice_Uint32(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 8},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]uint32)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint32{8})
}

func Test_asFixedSlice_Uint64(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.PositiveFixIntMin + 9},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]uint64)
	run(t, v1)
	tu.EqualSlice(t, *v1, []uint64{9})
}

func Test_asFixedSlice_Float32(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.Float32, 63, 128, 0, 0},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]float32)
	run(t, v1)
	tu.EqualSlice(t, *v1, []float32{1})
}

func Test_asFixedSlice_Float64(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.Float64, 63, 240, 0, 0, 0, 0, 0, 0},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]float64)
	run(t, v1)
	tu.EqualSlice(t, *v1, []float64{1})
}

func Test_asFixedSlice_String(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 2)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.FixStr + 1, 'a', def.FixStr + 1, 'b'},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]string)
	run(t, v1)
	tu.EqualSlice(t, *v1, []string{"a", "b"})
}

func Test_asFixedSlice_Bool(t *testing.T) {
	run := func(t *testing.T, v any) {
		method := func(d *decoder) (int, bool, error) {
			rv := reflect.ValueOf(v)
			return d.asFixedSlice(rv.Elem(), 0, 1)
		}

		testcases := AsXXXTestCases[bool]{
			{
				Name:           "error",
				Data:           []byte{},
				Error:          def.ErrTooShortBytes,
				MethodAsCustom: method,
			},
			{
				Name:           "ok",
				Data:           []byte{def.True},
				Expected:       true,
				MethodAsCustom: method,
			},
		}
		testcases.Run(t)
	}

	v1 := new([]bool)
	run(t, v1)
	tu.EqualSlice(t, *v1, []bool{true})
}