package decoding

import (
	"io"
	"reflect"
	"testing"
	"time"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_setStruct_ext(t *testing.T) {
	run := func(t *testing.T, rv reflect.Value) {

		method := func(d *decoder) func(byte, reflect.Kind) (any, error) {
			return func(code byte, k reflect.Kind) (any, error) {
				return nil, d.setStruct(code, rv, k)
			}
		}

		testcases := AsXXXTestCases[any]{
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
				Name:             "ExtCoder.ok",
				Code:             def.Fixext4,
				Data:             []byte{255, 0, 0, 0, 0},
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	ngDec := testExt2StreamDecoder{}
	AddExtDecoder(&ngDec)
	defer RemoveExtDecoder(&ngDec)

	v1 := new(time.Time)
	run(t, reflect.ValueOf(v1).Elem())
	tu.EqualEqualer(t, *v1, time.Unix(0, 0))
}

func Test_setStructFromMap(t *testing.T) {
	run := func(t *testing.T, rv reflect.Value) {
		method := func(d *decoder) func(byte, reflect.Kind) (any, error) {
			return func(code byte, k reflect.Kind) (any, error) {
				return nil, d.setStructFromMap(code, rv, k)
			}
		}

		testcases := AsXXXTestCases[any]{
			{
				Name:             "error.length",
				Code:             def.Map16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.key",
				Code:             def.Map16,
				Data:             []byte{0, 1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.decode",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'v', def.Array16},
				ReadCount:        4,
				Error:            def.ErrCanNotDecode,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.jump",
				Code:             def.Map16,
				Data:             []byte{0, 2, def.FixStr + 1, 'v', 0, def.FixStr + 1, 'b'},
				ReadCount:        6,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 7},
				ReadCount:        4,
				MethodAsWithCode: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	type st struct {
		V int `msgpack:"v"`
	}

	v1 := new(st)
	run(t, reflect.ValueOf(v1).Elem())
	tu.Equal(t, v1.V, 7)
}

func Test_setStructFromArray(t *testing.T) {
	run := func(t *testing.T, rv reflect.Value) {
		method := func(d *decoder) func(byte, reflect.Kind) (any, error) {
			return func(code byte, k reflect.Kind) (any, error) {
				return nil, d.setStructFromArray(code, rv, k)
			}
		}

		testcases := AsXXXTestCases[any]{
			{
				Name:             "error.length",
				Code:             def.Array16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.key",
				Code:             def.Array16,
				Data:             []byte{0, 1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.decode",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.Array16},
				ReadCount:        2,
				Error:            def.ErrCanNotDecode,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.jump",
				Code:             def.Array16,
				Data:             []byte{0, 2, 0},
				ReadCount:        2,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.PositiveFixIntMin + 8},
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		for _, tc := range testcases {
			tc.Run(t)
		}
	}

	type st struct {
		V int `msgpack:"v"`
	}

	v1 := new(st)
	run(t, reflect.ValueOf(v1).Elem())
	tu.Equal(t, v1.V, 8)
}
func Test_jumpOffset(t *testing.T) {
	method := func(d *decoder) (any, error) {
		return nil, d.jumpOffset()
	}

	testcases := AsXXXTestCases[any]{
		{
			Name:           "error.read.code",
			Data:           []byte{},
			ReadCount:      0,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "True.ok",
			Data:           []byte{def.True},
			ReadCount:      1,
			MethodAsCustom: method,
		},
		{
			Name:           "False.ok",
			Data:           []byte{def.False},
			ReadCount:      1,
			MethodAsCustom: method,
		},
		{
			Name:           "PositiveFixNum.ok",
			Data:           []byte{def.PositiveFixIntMin + 1},
			ReadCount:      1,
			MethodAsCustom: method,
		},
		{
			Name:           "NegativeFixNum.ok",
			Data:           []byte{0xf0},
			ReadCount:      1,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint8.error",
			Data:           []byte{def.Uint8},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Int8.error",
			Data:           []byte{def.Int8},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint8.ok",
			Data:           []byte{def.Uint8, 1},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Int8.ok",
			Data:           []byte{def.Int8, 1},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint16.error",
			Data:           []byte{def.Uint16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Int16.error",
			Data:           []byte{def.Int16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint16.ok",
			Data:           []byte{def.Uint16, 0, 1},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Int16.ok",
			Data:           []byte{def.Int16, 0, 1},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint32.error",
			Data:           []byte{def.Uint32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Int32.error",
			Data:           []byte{def.Int32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Float32.error",
			Data:           []byte{def.Float32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint32.ok",
			Data:           []byte{def.Uint32, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Int32.ok",
			Data:           []byte{def.Int32, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Float32.ok",
			Data:           []byte{def.Float32, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint64.error",
			Data:           []byte{def.Uint64},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Int64.error",
			Data:           []byte{def.Int64},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Float64.error",
			Data:           []byte{def.Float64},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Uint64.ok",
			Data:           []byte{def.Uint64, 0, 0, 0, 0, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Int64.ok",
			Data:           []byte{def.Int64, 0, 0, 0, 0, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Float64.ok",
			Data:           []byte{def.Float64, 0, 0, 0, 0, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "FixStr.ng",
			Data:           []byte{def.FixStr + 1},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "FixStr.ok",
			Data:           []byte{def.FixStr + 1, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Str8.ng.length",
			Data:           []byte{def.Str8},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Str8.ng.str",
			Data:           []byte{def.Str8, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Str8.ok",
			Data:           []byte{def.Str8, 1, 'a'},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin8.ng.length",
			Data:           []byte{def.Bin8},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin8.ng.str",
			Data:           []byte{def.Bin8, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin8.ok",
			Data:           []byte{def.Bin8, 1, 'a'},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Str16.ng.length",
			Data:           []byte{def.Str16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Str16.ng.str",
			Data:           []byte{def.Str16, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Str16.ok",
			Data:           []byte{def.Str16, 0, 1, 'a'},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin16.ng.length",
			Data:           []byte{def.Bin16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin16.ng.str",
			Data:           []byte{def.Bin16, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin16.ok",
			Data:           []byte{def.Bin16, 0, 1, 'a'},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Str32.ng.length",
			Data:           []byte{def.Str32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Str32.ng.str",
			Data:           []byte{def.Str32, 0, 0, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Str32.ok",
			Data:           []byte{def.Str32, 0, 0, 0, 1, 'a'},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin32.ng.length",
			Data:           []byte{def.Bin32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin32.ng.str",
			Data:           []byte{def.Bin32, 0, 0, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Bin32.ok",
			Data:           []byte{def.Bin32, 0, 0, 0, 1, 'a'},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "FixSlice.ng",
			Data:           []byte{def.FixArray + 1},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "FixSlice.ok",
			Data:           []byte{def.FixArray + 1, 0xc1},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Array16.ng.len",
			Data:           []byte{def.Array16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Array16.ng.jump",
			Data:           []byte{def.Array16, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Array16.ok",
			Data:           []byte{def.Array16, 0, 1, 0xc1},
			ReadCount:      3,
			MethodAsCustom: method,
		}, {
			Name:           "Array32.ng.len",
			Data:           []byte{def.Array32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Array32.ng.jump",
			Data:           []byte{def.Array32, 0, 0, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Array32.ok",
			Data:           []byte{def.Array32, 0, 0, 0, 1, 0xc1},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "FixMap.ng",
			Data:           []byte{def.FixMap + 1},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "FixMap.ok",
			Data:           []byte{def.FixMap + 1, 0xc1, 0xc1},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Map16.ng.len",
			Data:           []byte{def.Map16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Map16.ng.jump",
			Data:           []byte{def.Map16, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Map16.ok",
			Data:           []byte{def.Map16, 0, 1, 0xc1, 0xc1},
			ReadCount:      4,
			MethodAsCustom: method,
		}, {
			Name:           "Map32.ng.len",
			Data:           []byte{def.Map32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Map32.ng.jump",
			Data:           []byte{def.Map32, 0, 0, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Map32.ok",
			Data:           []byte{def.Map32, 0, 0, 0, 1, 0xc1, 0xc1},
			ReadCount:      4,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext1.ng",
			Data:           []byte{def.Fixext1},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext1.ok",
			Data:           []byte{def.Fixext1, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext2.ng",
			Data:           []byte{def.Fixext2},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext2.ok",
			Data:           []byte{def.Fixext2, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext4.ng",
			Data:           []byte{def.Fixext4},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext4.ok",
			Data:           []byte{def.Fixext4, 0, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext8.ng",
			Data:           []byte{def.Fixext8},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext8.ok",
			Data:           []byte{def.Fixext8, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext16.ng",
			Data:           []byte{def.Fixext16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Fixext16.ok",
			Data:           []byte{def.Fixext16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ReadCount:      2,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext8.ng.size",
			Data:           []byte{def.Ext8},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext8.ng.size.n",
			Data:           []byte{def.Ext8, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext8.ok",
			Data:           []byte{def.Ext8, 1, 0},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext16.ng.size",
			Data:           []byte{def.Ext16},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext16.ng.size.n",
			Data:           []byte{def.Ext16, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext16.ok",
			Data:           []byte{def.Ext16, 0, 1, 0},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext32.ng.size",
			Data:           []byte{def.Ext32},
			ReadCount:      1,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext32.ng.size.n",
			Data:           []byte{def.Ext32, 0, 0, 0, 1},
			ReadCount:      2,
			Error:          io.EOF,
			MethodAsCustom: method,
		},
		{
			Name:           "Ext32.ok",
			Data:           []byte{def.Ext32, 0, 0, 0, 1, 0},
			ReadCount:      3,
			MethodAsCustom: method,
		},
		{
			Name:           "Unexpected",
			Data:           []byte{0xc1},
			ReadCount:      1,
			MethodAsCustom: method,
		},
	}
	for _, tc := range testcases {
		tc.Run(t)
	}
}
