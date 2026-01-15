package encoding

import (
	"math"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_asInt(t *testing.T) {
	method := func(e *encoder) func(int64) error {
		return e.writeInt
	}
	testcases := AsXXXTestCases[int64]{
		{
			Name:         "Uint.error",
			Value:        math.MaxInt32,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:       "Uint.ok",
			Value:      math.MaxInt32,
			Expected:   []byte{def.Uint32, 0x7f, 0xff, 0xff, 0xff},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "NegativeFix.error",
			Value:        -1,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:       "NegativeFix.ok",
			Value:      -1,
			Expected:   []byte{0xff},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Int8.error.def",
			Value:        math.MinInt8,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Int8.error.value",
			Value:        math.MinInt8,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Int8},
			Method:       method,
		},
		{
			Name:       "Int8.ok",
			Value:      math.MinInt8,
			Expected:   []byte{def.Int8, 0x80},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Int16.error.def",
			Value:        math.MinInt16,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Int16.error.value",
			Value:        math.MinInt16,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Int16},
			Method:       method,
		},
		{
			Name:       "Int16.ok",
			Value:      math.MinInt16,
			Expected:   []byte{def.Int16, 0x80, 0x00},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Int32.error.def",
			Value:        math.MinInt32,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Int32.error.value",
			Value:        math.MinInt32,
			BufferSize:   5,
			PreWriteSize: 1,
			Contains:     []byte{def.Int32},
			Method:       method,
		},
		{
			Name:       "Int32.ok",
			Value:      math.MinInt32,
			Expected:   []byte{def.Int32, 0x80, 0x00, 0x00, 0x00},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Int64.error.def",
			Value:        math.MinInt64,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Int64.error.value",
			Value:        math.MinInt64,
			BufferSize:   9,
			PreWriteSize: 1,
			Contains:     []byte{def.Int64},
			Method:       method,
		},
		{
			Name:       "Int64.ok",
			Value:      math.MinInt64,
			Expected:   []byte{def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
