package encoding

import (
	"math"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_asUint(t *testing.T) {
	method := func(e *encoder) func(uint64) error {
		return e.writeUint
	}
	testcases := AsXXXTestCases[uint64]{
		{
			Name:         "PositiveFix.error",
			Value:        math.MaxInt8,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:       "PositiveFix.ok",
			Value:      math.MaxInt8,
			Expected:   []byte{0x7f},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Uint8.error.def",
			Value:        math.MaxUint8,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Uint8.error.value",
			Value:        math.MaxUint8,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Uint8},
			Method:       method,
		},
		{
			Name:       "Uint8.ok",
			Value:      math.MaxUint8,
			Expected:   []byte{def.Uint8, 0xff},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Uint16.error.def",
			Value:        math.MaxUint16,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Uint16.error.value",
			Value:        math.MaxUint16,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Uint16},
			Method:       method,
		},
		{
			Name:       "Uint16.ok",
			Value:      math.MaxUint16,
			Expected:   []byte{def.Uint16, 0xff, 0xff},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Uint32.error.def",
			Value:        math.MaxUint32,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Uint32.error.value",
			Value:        math.MaxUint32,
			BufferSize:   5,
			PreWriteSize: 1,
			Contains:     []byte{def.Uint32},
			Method:       method,
		},
		{
			Name:       "Uint32.ok",
			Value:      math.MaxUint32,
			Expected:   []byte{def.Uint32, 0xff, 0xff, 0xff, 0xff},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Uint64.error.def",
			Value:        math.MaxUint64,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Uint64.error.value",
			Value:        math.MaxUint64,
			BufferSize:   9,
			PreWriteSize: 1,
			Contains:     []byte{def.Uint64},
			Method:       method,
		},
		{
			Name:       "Uint64.ok",
			Value:      math.MaxUint64,
			Expected:   []byte{def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
