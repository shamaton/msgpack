package decoding

import (
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_asInt(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) (int64, error) {
		return d.asInt
	}
	testcases := AsXXXTestCases[int64]{
		{
			Name:      "error",
			Data:      []byte{},
			Error:     io.EOF,
			ReadCount: 1,
			MethodAs:  method,
		},
		{
			Name:      "ok",
			Data:      []byte{def.Int8, 1},
			Expected:  1,
			ReadCount: 2,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asIntWithCode(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (int64, error) {
		return d.asIntWithCode
	}
	testcases := AsXXXTestCases[int64]{
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
			Expected:         1,
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
			Data:             []byte{1},
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint16.error",
			Code:             def.Uint16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint16.ok",
			Code:             def.Uint16,
			Data:             []byte{0, 1},
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int16.error",
			Code:             def.Int16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int16.ok",
			Code:             def.Int16,
			Data:             []byte{0, 1},
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint32.error",
			Code:             def.Uint32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint32.ok",
			Code:             def.Uint32,
			Data:             []byte{0, 0, 0, 1},
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int32.error",
			Code:             def.Int32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int32.ok",
			Code:             def.Int32,
			Data:             []byte{0, 0, 0, 1},
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint64.error",
			Code:             def.Uint64,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint64.ok",
			Code:             def.Uint64,
			Data:             []byte{0, 0, 0, 0, 0, 0, 0, 1},
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int64.error",
			Code:             def.Int64,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int64.ok",
			Code:             def.Int64,
			Data:             []byte{0, 0, 0, 0, 0, 0, 0, 1},
			Expected:         1,
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
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
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
			Expected:         1,
			ReadCount:        1,
			MethodAsWithCode: method,
		},
		{
			Name:             "Nil",
			Code:             def.Nil,
			Expected:         0,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Array16,
			Error:            ErrCanNotDecode,
			MethodAsWithCode: method,
		},
	}
	for _, tc := range testcases {
		tc.Run(t)
	}
}
