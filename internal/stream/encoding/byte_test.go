package encoding

import (
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_writeByteSliceLength(t *testing.T) {
	method := func(e *encoder) func(int) error {
		return e.writeByteSliceLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:         "Bin8.error.def",
			Value:        5,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Bin8.error.value",
			Value:        5,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Bin8},
			Method:       method,
		},
		{
			Name:       "Bin8.ok",
			Value:      5,
			Expected:   []byte{def.Bin8, 0x05},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Bin16.error.def",
			Value:        256,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Bin16.error.value",
			Value:        256,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Bin16},
			Method:       method,
		},
		{
			Name:       "Bin16.ok",
			Value:      256,
			Expected:   []byte{def.Bin16, 0x01, 0x00},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Bin32.error.def",
			Value:        65536,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Bin32.error.value",
			Value:        65536,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Bin32},
			Method:       method,
		},
		{
			Name:       "Bin32.ok",
			Value:      65536,
			Expected:   []byte{def.Bin32, 0x00, 0x01, 0x00, 0x00},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
