package decoding

import (
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_asComplex64(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (complex64, error) {
		return d.asComplex64
	}
	testcases := AsXXXTestCases[complex64]{
		{
			Name:             "Fixext8.error.type",
			Code:             def.Fixext8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext8.error.r",
			Code:             def.Fixext8,
			Data:             []byte{byte(def.ComplexTypeCode())},
			Error:            io.EOF,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext8.error.i",
			Code:             def.Fixext8,
			Data:             []byte{byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:            io.EOF,
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext8.ok",
			Code:             def.Fixext8,
			Data:             []byte{byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
			Expected:         complex(1, 1),
			ReadCount:        3,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext16.error.type",
			Code:             def.Fixext16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext16.error.r",
			Code:             def.Fixext16,
			Data:             []byte{byte(def.ComplexTypeCode())},
			Error:            io.EOF,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext16.error.i",
			Code:             def.Fixext16,
			Data:             []byte{byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:            io.EOF,
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name: "Fixext16.ok",
			Code: def.Fixext16,
			Data: []byte{byte(def.ComplexTypeCode()),
				63, 240, 0, 0, 0, 0, 0, 0, 63, 240, 0, 0, 0, 0, 0, 0},
			Expected:         complex(1, 1),
			ReadCount:        3,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Nil,
			IsTemplateError:  true,
			MethodAsWithCode: method,
		},
	}
	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asComplex128(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (complex128, error) {
		return d.asComplex128
	}
	testcases := AsXXXTestCases[complex128]{
		{
			Name:             "Fixext8.error.type",
			Code:             def.Fixext8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext8.error.r",
			Code:             def.Fixext8,
			Data:             []byte{byte(def.ComplexTypeCode())},
			Error:            io.EOF,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext8.error.i",
			Code:             def.Fixext8,
			Data:             []byte{byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:            io.EOF,
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext8.ok",
			Code:             def.Fixext8,
			Data:             []byte{byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
			Expected:         complex(1, 1),
			ReadCount:        3,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext16.error.type",
			Code:             def.Fixext16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext16.error.r",
			Code:             def.Fixext16,
			Data:             []byte{byte(def.ComplexTypeCode())},
			Error:            io.EOF,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Fixext16.error.i",
			Code:             def.Fixext16,
			Data:             []byte{byte(def.ComplexTypeCode()), 0, 0, 0, 1},
			Error:            io.EOF,
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name: "Fixext16.ok",
			Code: def.Fixext16,
			Data: []byte{byte(def.ComplexTypeCode()),
				63, 240, 0, 0, 0, 0, 0, 0, 63, 240, 0, 0, 0, 0, 0, 0},
			Expected:         complex(1, 1),
			ReadCount:        3,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Nil,
			IsTemplateError:  true,
			MethodAsWithCode: method,
		},
	}
	for _, tc := range testcases {
		tc.Run(t)
	}
}
