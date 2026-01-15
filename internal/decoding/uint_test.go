package decoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_asUint(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (uint64, int, error) {
		return d.asUint
	}
	testcases := AsXXXTestCases[uint64]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "PositiveFixNum.ok",
			Data:     []byte{def.PositiveFixIntMin + 1},
			Expected: uint64(1),
			MethodAs: method,
		},
		{
			Name:     "Uint8.error",
			Data:     []byte{def.Uint8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint8.ok",
			Data:     []byte{def.Uint8, 1},
			Expected: uint64(1),
			MethodAs: method,
		},
		{
			Name:     "Uint16.error",
			Data:     []byte{def.Uint16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint16.ok",
			Data:     []byte{def.Uint16, 0, 1},
			Expected: uint64(1),
			MethodAs: method,
		},
		{
			Name:     "Uint32.error",
			Data:     []byte{def.Uint32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint32.ok",
			Data:     []byte{def.Uint32, 0, 0, 0, 1},
			Expected: uint64(1),
			MethodAs: method,
		},
		{
			Name:     "Uint64.error",
			Data:     []byte{def.Uint64},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint64.ok",
			Data:     []byte{def.Uint64, 0, 0, 0, 0, 0, 0, 0, 1},
			Expected: uint64(1),
			MethodAs: method,
		},
		{
			Name:     "NegativeFixNum.ok",
			Data:     []byte{0xff},
			Expected: uint64(math.MaxUint64),
			MethodAs: method,
		},
		{
			Name:     "Int8.error",
			Data:     []byte{def.Int8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int8.ok",
			Data:     []byte{def.Int8, 0xff},
			Expected: uint64(math.MaxUint64),
			MethodAs: method,
		},
		{
			Name:     "Int16.error",
			Data:     []byte{def.Int16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int16.ok",
			Data:     []byte{def.Int16, 0xff, 0xff},
			Expected: uint64(math.MaxUint64),
			MethodAs: method,
		},
		{
			Name:     "Int32.error",
			Data:     []byte{def.Int32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int32.ok",
			Data:     []byte{def.Int32, 0xff, 0xff, 0xff, 0xff},
			Expected: uint64(math.MaxUint64),
			MethodAs: method,
		},
		{
			Name:     "Int64.error",
			Data:     []byte{def.Int64},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int64.ok",
			Data:     []byte{def.Int64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			Expected: uint64(math.MaxUint64),
			MethodAs: method,
		},
		{
			Name:     "Nil.ok",
			Data:     []byte{def.Nil},
			Expected: uint64(0),
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{def.Str8},
			Error:    def.ErrCanNotDecode,
			MethodAs: method,
		},
	}
	testcases.Run(t)
}
