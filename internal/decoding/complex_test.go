package decoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_asComplex64(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (complex64, int, error) {
		return d.asComplex64
	}
	testcases := AsXXXTestCases[complex64]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.error.type",
			Data:     []byte{def.Fixext8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.error.r",
			Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode())},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.error.i",
			Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.ok",
			Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
			Expected: complex(1, 1),
			MethodAs: method,
		},
		{
			Name:     "Fixext16.error.type",
			Data:     []byte{def.Fixext16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext16.error.r",
			Data:     []byte{def.Fixext16, byte(def.ComplexTypeCode())},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext16.error.i",
			Data:     []byte{def.Fixext16, byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name: "Fixext16.ok",
			Data: []byte{def.Fixext16, byte(def.ComplexTypeCode()),
				63, 240, 0, 0, 0, 0, 0, 0, 63, 240, 0, 0, 0, 0, 0, 0},
			Expected: complex(1, 1),
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{def.Nil},
			Error:    def.ErrCanNotDecode,
			MethodAs: method,
		},
	}
	testcases.Run(t)
}

func Test_asComplex128(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (complex128, int, error) {
		return d.asComplex128
	}
	testcases := AsXXXTestCases[complex128]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.error.type",
			Data:     []byte{def.Fixext8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.error.r",
			Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode())},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.error.i",
			Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext8.ok",
			Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
			Expected: complex(1, 1),
			MethodAs: method,
		},
		{
			Name:     "Fixext16.error.type",
			Data:     []byte{def.Fixext16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext16.error.r",
			Data:     []byte{def.Fixext16, byte(def.ComplexTypeCode())},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Fixext16.error.i",
			Data:     []byte{def.Fixext16, byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name: "Fixext16.ok",
			Data: []byte{def.Fixext16, byte(def.ComplexTypeCode()),
				63, 240, 0, 0, 0, 0, 0, 0, 63, 240, 0, 0, 0, 0, 0, 0},
			Expected: complex(1, 1),
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{def.Nil},
			Error:    def.ErrCanNotDecode,
			MethodAs: method,
		},
	}
	testcases.Run(t)
}
