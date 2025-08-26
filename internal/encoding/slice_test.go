package encoding

import (
	"reflect"
	"testing"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_FixedSlice(t *testing.T) {
	testcases := []struct {
		value any
		size  int
	}{
		{
			value: []int{-1},
			size:  2,
		},
		{
			value: []uint{1},
			size:  2,
		},
		{
			value: []string{"a"},
			size:  3,
		},
		{
			value: []float32{1.23},
			size:  6,
		},
		{
			value: []float64{1.23},
			size:  10,
		},
		{
			value: []bool{true},
			size:  2,
		},
		{
			value: []int8{1},
			size:  2,
		},
		{
			value: []int16{1},
			size:  2,
		},
		{
			value: []int32{1},
			size:  2,
		},
		{
			value: []int64{1},
			size:  2,
		},
		{
			value: []uint8{1},
			size:  2,
		},
		{
			value: []uint16{1},
			size:  2,
		},
		{
			value: []uint32{1},
			size:  2,
		},
		{
			value: []uint64{1},
			size:  2,
		},
	}
	for _, tc := range testcases {
		rv := reflect.ValueOf(tc.value)
		t.Run(rv.Type().String(), func(t *testing.T) {
			e := encoder{}
			size, b := e.calcFixedSlice(rv)
			tu.Equal(t, b, true)
			tu.Equal(t, size, tc.size)

			e.d = make([]byte, size)
			result, b := e.writeFixedSlice(rv, 0)
			tu.Equal(t, b, true)
			tu.Equal(t, result, size)
		})
	}
}
