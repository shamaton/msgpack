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
			Error:            ErrCanNotDecode,
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
