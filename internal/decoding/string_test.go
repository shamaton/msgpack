package decoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_stringByteLength(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (int, int, error) {
		return d.stringByteLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "FixStr.ok",
			Data:     []byte{def.FixStr + 1},
			Expected: 1,
			MethodAs: method,
		},
		{
			Name:     "Str8.error",
			Data:     []byte{def.Str8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str8.ok",
			Data:     []byte{def.Str8, 0xff},
			Expected: math.MaxUint8,
			MethodAs: method,
		},
		{
			Name:     "Str16.error",
			Data:     []byte{def.Str16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str16.ok",
			Data:     []byte{def.Str16, 0xff, 0xff},
			Expected: math.MaxUint16,
			MethodAs: method,
		},
		{
			Name:     "Str32.error",
			Data:     []byte{def.Str32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str32.ok",
			Data:     []byte{def.Str32, 0xff, 0xff, 0xff, 0xff},
			Expected: math.MaxUint32,

			MethodAs: method,
		},
		{
			Name:     "Nil",
			Data:     []byte{def.Nil},
			Expected: 0,
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{def.Array16},
			Error:    def.ErrCanNotDecode,
			MethodAs: method,
		},
	}
	testcases.Run(t)
}

func Test_asString(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (string, int, error) {
		return d.asString
	}
	testcases := AsXXXTestCases[string]{
		{
			Name:     "error.string",
			Data:     []byte{def.FixStr + 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "ok",
			Data:     []byte{def.FixStr + 1, 'a'},
			Expected: "a",
			MethodAs: method,
		},
	}
	testcases.Run(t)
}

func Test_asStringByte(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) ([]byte, int, error) {
		return d.asStringByte
	}
	testcases := AsXXXTestCases[[]byte]{
		{
			Name:     "error",
			Data:     []byte{def.FixStr + 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "ok",
			Data:     []byte{def.FixStr + 1, 'a'},
			Expected: []byte{'a'},
			MethodAs: method,
		},
	}
	testcases.Run(t)
}
