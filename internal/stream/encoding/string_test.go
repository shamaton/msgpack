package encoding

import (
	"math"
	"strings"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_writeString(t *testing.T) {
	method := func(e *encoder) func(string) error {
		return e.writeString
	}
	str8 := strings.Repeat("a", math.MaxUint8)
	str16 := strings.Repeat("a", math.MaxUint16)
	str32 := strings.Repeat("a", math.MaxUint16+1)
	testcases := AsXXXTestCases[string]{
		{
			Name:         "Str8.error.def",
			Value:        str8,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Str8.error.length",
			Value:        str8,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Str8},
			Method:       method,
		},
		{
			Name:         "Str8.error.string",
			Value:        str8,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Str8, 0xff},
			Method:       method,
		},
		{
			Name:  "Str8.ok",
			Value: str8,
			Expected: append(
				[]byte{def.Str8, 0xff},
				[]byte(str8)...,
			),
			BufferSize: 1,
			Method:     method,
		}, {
			Name:         "Str16.error.def",
			Value:        str16,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Str16.error.length",
			Value:        str16,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Str16},
			Method:       method,
		},
		{
			Name:         "Str16.error.string",
			Value:        str16,
			BufferSize:   4,
			PreWriteSize: 1,
			Contains:     []byte{def.Str16, 0xff, 0xff},
			Method:       method,
		},
		{
			Name:  "Str16.ok",
			Value: str16,
			Expected: append(
				[]byte{def.Str16, 0xff, 0xff},
				[]byte(str16)...,
			),
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Str32.error.def",
			Value:        str32,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Str32.error.length",
			Value:        str32,
			BufferSize:   2,
			PreWriteSize: 1,
			Contains:     []byte{def.Str32},
			Method:       method,
		},
		{
			Name:         "Str32.error.string",
			Value:        str32,
			BufferSize:   6,
			PreWriteSize: 1,
			Contains:     []byte{def.Str32, 0x00, 0x01, 0x00, 0x00},
			Method:       method,
		},
		{
			Name:  "Str32.ok",
			Value: str32,
			Expected: append(
				[]byte{def.Str32, 0x00, 0x01, 0x00, 0x00},
				[]byte(str32)...,
			),
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
