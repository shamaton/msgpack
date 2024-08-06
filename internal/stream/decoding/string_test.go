package decoding

import (
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_stringByteLength(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (int, error) {
		return d.stringByteLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:             "FixStr",
			Code:             def.FixStr + 1,
			Expected:         1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Str8",
			Code:             def.Str8,
			Data:             []byte{0xff},
			Expected:         math.MaxUint8,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Str16",
			Code:             def.Str16,
			Data:             []byte{0xff, 0xff},
			Expected:         math.MaxUint16,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Str32",
			Code:             def.Str32,
			Data:             []byte{0xff, 0xff, 0xff, 0xff},
			Expected:         math.MaxUint32,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Nil",
			Code:             def.Nil,
			Expected:         0,
			MethodAsWithCode: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asString(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) (string, error) {
		return d.asString
	}
	testcases := AsXXXTestCases[string]{
		{
			Name:      "String.error",
			Data:      []byte{def.FixStr + 1},
			Error:     io.EOF,
			ReadCount: 1,
			MethodAs:  method,
		},
		{
			Name:      "String.ok",
			Data:      []byte{def.FixStr + 1, 'a'},
			Expected:  "a",
			ReadCount: 2,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asStringByte(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) ([]byte, error) {
		return d.asStringByte
	}
	testcases := AsXXXTestCases[[]byte]{
		{
			Name:      "error",
			Data:      []byte{def.FixStr + 1},
			Error:     io.EOF,
			ReadCount: 1,
			MethodAs:  method,
		},
		{
			Name:      "ok",
			Data:      []byte{def.FixStr + 1, 'a'},
			Expected:  []byte{'a'},
			ReadCount: 2,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}
