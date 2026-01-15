package decoding

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

func Test_asInterface(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) (any, error) {
		return d.asInterface
	}
	testcases := AsXXXTestCases[any]{
		{
			Name:      "error",
			Data:      []byte{},
			Error:     io.EOF,
			ReadCount: 0,
			MethodAs:  method,
		},
		{
			Name:      "ok",
			Data:      []byte{def.Nil},
			Expected:  nil,
			ReadCount: 1,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asInterfaceWithCode(t *testing.T) {
	dec := testExt2StreamDecoder{}
	AddExtDecoder(&dec)
	defer RemoveExtDecoder(&dec)

	method := func(d *decoder) func(byte, reflect.Kind) (any, error) {
		return d.asInterfaceWithCode
	}
	testcases := AsXXXTestCases[any]{
		{
			Name:             "Uint8.error",
			Code:             def.Uint8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint16.error",
			Code:             def.Uint16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint32.error",
			Code:             def.Uint32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Uint64.error",
			Code:             def.Uint64,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int8.error",
			Code:             def.Int8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int16.error",
			Code:             def.Int16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int32.error",
			Code:             def.Int32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Int64.error",
			Code:             def.Int64,
			Error:            io.EOF,
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
			Expected:         float32(1),
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
			Name:             "Str.error",
			Code:             def.Str8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin.error",
			Code:             def.Bin8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Array.error.length",
			Code:             def.Array16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Array.error.set",
			Code:             def.Array16,
			Data:             []byte{0, 1},
			ReadCount:        1,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Map.error.length",
			Code:             def.Map16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Map.error.set.key_code",
			Code:             def.Map16,
			Data:             []byte{0, 1},
			ReadCount:        1,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Map.error.set.can.slice",
			Code:             def.Map16,
			Data:             []byte{0, 1, def.Array16},
			ReadCount:        2,
			Error:            def.ErrCanNotSetSliceAsMapKey,
			MethodAsWithCode: method,
		},
		{
			Name:             "Map.error.set.can.map",
			Code:             def.Map16,
			Data:             []byte{0, 1, def.Map16},
			ReadCount:        2,
			Error:            def.ErrCanNotSetMapAsMapKey,
			MethodAsWithCode: method,
		},
		{
			Name:             "Map.error.set.key",
			Code:             def.Map16,
			Data:             []byte{0, 1, def.FixStr + 1},
			ReadCount:        2,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Map.error.set.value",
			Code:             def.Map16,
			Data:             []byte{0, 1, def.FixStr + 1, 'a'},
			ReadCount:        3,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Ext.error",
			Code:             def.Fixext1,
			Data:             []byte{},
			ReadCount:        0,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "ExtCoder.error",
			Code:             def.Fixext1,
			Data:             []byte{3, 0},
			ReadCount:        2,
			Error:            ErrTestExtStreamDecoder,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Fixext1,
			Data:             []byte{4, 0},
			ReadCount:        2,
			Error:            def.ErrCanNotDecode,
			MethodAsWithCode: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

// TODO: to testutil
type testExt2StreamDecoder struct{}

var _ ext.StreamDecoder = (*testExt2StreamDecoder)(nil)

func (td *testExt2StreamDecoder) Code() int8 {
	return 3
}

func (td *testExt2StreamDecoder) IsType(_ byte, code int8, _ int) bool {
	return code == td.Code()
}

var ErrTestExtStreamDecoder = fmt.Errorf("testExtStreamDecoder")

func (td *testExt2StreamDecoder) ToValue(_ byte, _ []byte, k reflect.Kind) (any, error) {
	return nil, ErrTestExtStreamDecoder
}
