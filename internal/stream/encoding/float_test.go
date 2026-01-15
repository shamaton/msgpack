package encoding

import (
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_writeFloat32(t *testing.T) {
	method := func(e *encoder) func(float64) error {
		return e.writeFloat32
	}
	v := 1.23
	testcases := AsXXXTestCases[float64]{
		{
			Name:         "error.def",
			Value:        v,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "error.value",
			Value:        v,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Float32},
			Method:       method,
		},
		{
			Name:  "ok",
			Value: v,
			Expected: []byte{
				def.Float32,
				0x3f, 0x9d, 0x70, 0xa4,
			},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}

func Test_writeFloat64(t *testing.T) {
	method := func(e *encoder) func(float64) error {
		return e.writeFloat64
	}
	v := 1.23
	testcases := AsXXXTestCases[float64]{
		{
			Name:         "error.def",
			Value:        v,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "error.value",
			Value:        v,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Float64},
			Method:       method,
		},
		{
			Name:  "ok",
			Value: v,
			Expected: []byte{
				def.Float64,
				0x3f, 0xf3, 0xae, 0x14, 0x7a, 0xe1, 0x47, 0xae,
			},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
