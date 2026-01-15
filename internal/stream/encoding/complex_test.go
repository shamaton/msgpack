package encoding

import (
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_writeComplex64(t *testing.T) {
	method := func(e *encoder) func(complex64) error {
		return e.writeComplex64
	}
	v := complex64(complex(1, 2))
	testcases := AsXXXTestCases[complex64]{
		{
			Name:         "error.def",
			Value:        v,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "error.code",
			Value:        v,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Fixext8},
			Method:       method,
		},
		{
			Name:         "error.real",
			Value:        v,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Fixext8, byte(def.ComplexTypeCode())},
			Method:       method,
		},
		{
			Name:         "error.imag",
			Value:        v,
			BufferSize:   7,
			PreWriteSize: 1,
			Contains: []byte{
				def.Fixext8, byte(def.ComplexTypeCode()),
				0x3f, 0x80, 0x00, 0x00,
			},
			Method: method,
		},
		{
			Name:  "ok",
			Value: v,
			Expected: []byte{
				def.Fixext8, byte(def.ComplexTypeCode()),
				0x3f, 0x80, 0x00, 0x00,
				0x40, 0x00, 0x00, 0x00,
			},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}

func Test_writeComplex128(t *testing.T) {
	method := func(e *encoder) func(complex128) error {
		return e.writeComplex128
	}
	v := complex128(complex(1, 2))
	testcases := AsXXXTestCases[complex128]{
		{
			Name:         "error.def",
			Value:        v,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "error.code",
			Value:        v,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Fixext16},
			Method:       method,
		},
		{
			Name:         "error.real",
			Value:        v,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Fixext16, byte(def.ComplexTypeCode())},
			Method:       method,
		},
		{
			Name:         "error.imag",
			Value:        v,
			BufferSize:   11,
			PreWriteSize: 1,
			Contains: []byte{
				def.Fixext16, byte(def.ComplexTypeCode()),
				0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			Method: method,
		},
		{
			Name:  "ok",
			Value: v,
			Expected: []byte{
				def.Fixext16, byte(def.ComplexTypeCode()),
				0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
