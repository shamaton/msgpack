package decoding

import (
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_asFloat32(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) (float32, error) {
		return d.asFloat32
	}
	testcases := AsXXXTestCases[float32]{
		{
			Name:     "error",
			Error:    io.EOF,
			MethodAs: method,
		},
		{
			Name:      "ok",
			Data:      []byte{def.Float32, 63, 128, 0, 0},
			Expected:  float32(1),
			ReadCount: 2,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asFloat32WithCode(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (float32, error) {
		return d.asFloat32WithCode
	}
	testcases := AsXXXTestCases[float32]{
		{
			Name:             "Float32.error",
			Code:             def.Float32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Float32.ok",
			Code:             def.Float32,
			Data:             []byte{63, 128, 0, 0},
			Expected:         float32(1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint8.error",
			Code:             def.Uint8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint8.ok",
			Code:             def.Uint8,
			Data:             []byte{1},
			Expected:         float32(1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int8.error",
			Code:             def.Int8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int8.ok",
			Code:             def.Int8,
			Data:             []byte{0xff},
			Expected:         float32(-1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Nil.ok",
			Code:             def.Nil,
			Expected:         float32(0),
			ReadCount:        0,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Str8,
			Error:            def.ErrCanNotDecode,
			MethodAsWithCode: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asFloat64(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) (float64, error) {
		return d.asFloat64
	}
	testcases := AsXXXTestCases[float64]{
		{
			Name:     "error",
			Error:    io.EOF,
			MethodAs: method,
		},
		{
			Name:      "ok",
			Data:      []byte{def.Float64, 63, 240, 0, 0, 0, 0, 0, 0},
			Expected:  float64(1),
			ReadCount: 2,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asFloat64WithCode(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (float64, error) {
		return d.asFloat64WithCode
	}
	testcases := AsXXXTestCases[float64]{
		{
			Name:             "Float64.error",
			Code:             def.Float64,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Float64.ok",
			Code:             def.Float64,
			Data:             []byte{63, 240, 0, 0, 0, 0, 0, 0},
			Expected:         float64(1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Float32.error",
			Code:             def.Float32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Float32.ok",
			Code:             def.Float32,
			Data:             []byte{63, 128, 0, 0},
			Expected:         float64(1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint8.error",
			Code:             def.Uint8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint8.ok",
			Code:             def.Uint8,
			Data:             []byte{1},
			Expected:         float64(1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int8.error",
			Code:             def.Int8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int8.ok",
			Code:             def.Int8,
			Data:             []byte{0xff},
			Expected:         float64(-1),
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Nil.ok",
			Code:             def.Nil,
			Expected:         float64(0),
			ReadCount:        0,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Str8,
			Error:            def.ErrCanNotDecode,
			MethodAsWithCode: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}
