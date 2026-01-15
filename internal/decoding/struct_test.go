package decoding

import (
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3/def"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func Test_setStruct_ext(t *testing.T) {
	run := func(t *testing.T, rv reflect.Value) {

		method := func(d *decoder) func(int, reflect.Kind) (any, int, error) {
			return func(offset int, k reflect.Kind) (any, int, error) {
				o, err := d.setStruct(rv, offset, k)
				return nil, o, err
			}
		}

		testcases := AsXXXTestCases[any]{
			{
				Name:     "ExtCoder.error",
				Data:     []byte{def.Fixext1, 3, 0},
				Error:    ErrTestExtDecoder,
				MethodAs: method,
			},
			{
				Name:     "ExtCoder.ok",
				Data:     []byte{def.Fixext4, 255, 0, 0, 0, 0},
				MethodAs: method,
			},
		}
		testcases.Run(t)
	}

	ngDec := testExt2Decoder{}
	AddExtDecoder(&ngDec)
	defer RemoveExtDecoder(&ngDec)

	v1 := new(time.Time)
	run(t, reflect.ValueOf(v1).Elem())
	tu.EqualEqualer(t, *v1, time.Unix(0, 0))
}

func Test_setStructFromMap(t *testing.T) {
	run := func(t *testing.T, rv reflect.Value) {
		method := func(d *decoder) func(int, reflect.Kind) (any, int, error) {
			return func(offset int, k reflect.Kind) (any, int, error) {
				o, err := d.setStructFromMap(rv, offset, k)
				return nil, o, err
			}
		}

		testcases := AsXXXTestCases[any]{
			{
				Name:     "error.length",
				Data:     []byte{},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.required",
				Data:     []byte{def.Map16, 0, 1},
				Error:    def.ErrLackDataLengthToMap,
				MethodAs: method,
			},
			{
				Name:     "error.key",
				Data:     []byte{def.Map16, 0, 1, def.Str16, 0},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.decode",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'v', def.Array16},
				Error:    def.ErrCanNotDecode,
				MethodAs: method,
			},
			{
				Name:     "error.jump",
				Data:     []byte{def.Map16, 0, 2, def.FixStr + 1, 'v', 0, def.FixStr + 1, 'b'},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 7},
				MethodAs: method,
			},
		}
		testcases.Run(t)
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
		method := func(d *decoder) func(int, reflect.Kind) (any, int, error) {
			return func(offset int, k reflect.Kind) (any, int, error) {
				o, err := d.setStructFromArray(rv, offset, k)
				return nil, o, err
			}
		}

		testcases := AsXXXTestCases[any]{
			{
				Name:     "error.length",
				Data:     []byte{},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.required",
				Data:     []byte{def.Array16, 0, 1},
				Error:    def.ErrLackDataLengthToSlice,
				MethodAs: method,
			},
			{
				Name:     "error.decode",
				Data:     []byte{def.Array16, 0, 1, def.Array16},
				Error:    def.ErrCanNotDecode,
				MethodAs: method,
			},
			{
				Name:     "error.jump",
				Data:     []byte{def.Array16, 0, 2, 0, def.Array16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Array16, 0, 1, def.PositiveFixIntMin + 8},
				MethodAs: method,
			},
		}
		testcases.Run(t)
	}

	type st struct {
		V int `msgpack:"v"`
	}

	v1 := new(st)
	run(t, reflect.ValueOf(v1).Elem())
	tu.Equal(t, v1.V, 8)
}

func Test_jumpOffset(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (any, int, error) {
		return func(offset int, _ reflect.Kind) (any, int, error) {
			o, err := d.jumpOffset(offset)
			return nil, o, err
		}
	}

	testcases := AsXXXTestCases[any]{
		{
			Name:     "error.read.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "True.ok",
			Data:     []byte{def.True},
			MethodAs: method,
		},
		{
			Name:     "False.ok",
			Data:     []byte{def.False},
			MethodAs: method,
		},
		{
			Name:     "PositiveFixNum.ok",
			Data:     []byte{def.PositiveFixIntMin + 1},
			MethodAs: method,
		},
		{
			Name:     "NegativeFixNum.ok",
			Data:     []byte{0xf0},
			MethodAs: method,
		},
		{
			Name:     "Uint8.ok",
			Data:     []byte{def.Uint8, 1},
			MethodAs: method,
		},
		{
			Name:     "Int8.ok",
			Data:     []byte{def.Int8, 1},
			MethodAs: method,
		},
		{
			Name:     "Uint16.ok",
			Data:     []byte{def.Uint16, 0, 1},
			MethodAs: method,
		},
		{
			Name:     "Int16.ok",
			Data:     []byte{def.Int16, 0, 1},
			MethodAs: method,
		},
		{
			Name:     "Uint32.ok",
			Data:     []byte{def.Uint32, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Int32.ok",
			Data:     []byte{def.Int32, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Float32.ok",
			Data:     []byte{def.Float32, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Uint64.ok",
			Data:     []byte{def.Uint64, 0, 0, 0, 0, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Int64.ok",
			Data:     []byte{def.Int64, 0, 0, 0, 0, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Float64.ok",
			Data:     []byte{def.Float64, 0, 0, 0, 0, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "FixStr.ok",
			Data:     []byte{def.FixStr + 1, 0},
			MethodAs: method,
		},
		{
			Name:     "Str8.ng.length",
			Data:     []byte{def.Str8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str8.ok",
			Data:     []byte{def.Str8, 1, 'a'},
			MethodAs: method,
		},
		{
			Name:     "Bin8.ng.length",
			Data:     []byte{def.Bin8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin8.ok",
			Data:     []byte{def.Bin8, 1, 'a'},
			MethodAs: method,
		},
		{
			Name:     "Str16.ng.length",
			Data:     []byte{def.Str16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str16.ok",
			Data:     []byte{def.Str16, 0, 1, 'a'},
			MethodAs: method,
		},
		{
			Name:     "Bin16.ng.length",
			Data:     []byte{def.Bin16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin16.ok",
			Data:     []byte{def.Bin16, 0, 1, 'a'},
			MethodAs: method,
		},
		{
			Name:     "Str32.ng.length",
			Data:     []byte{def.Str32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str32.ok",
			Data:     []byte{def.Str32, 0, 0, 0, 1, 'a'},
			MethodAs: method,
		},
		{
			Name:     "Bin32.ng.length",
			Data:     []byte{def.Bin32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin32.ok",
			Data:     []byte{def.Bin32, 0, 0, 0, 1, 'a'},
			MethodAs: method,
		},
		{
			Name:     "FixSlice.ng",
			Data:     []byte{def.FixArray + 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "FixSlice.ok",
			Data:     []byte{def.FixArray + 1, 0xc1},
			MethodAs: method,
		},
		{
			Name:     "Array16.ng.len",
			Data:     []byte{def.Array16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array16.ng.jump",
			Data:     []byte{def.Array16, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array16.ok",
			Data:     []byte{def.Array16, 0, 1, 0xc1},
			MethodAs: method,
		}, {
			Name:     "Array32.ng.len",
			Data:     []byte{def.Array32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array32.ng.jump",
			Data:     []byte{def.Array32, 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array32.ok",
			Data:     []byte{def.Array32, 0, 0, 0, 1, 0xc1},
			MethodAs: method,
		},
		{
			Name:     "FixMap.ng",
			Data:     []byte{def.FixMap + 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "FixMap.ok",
			Data:     []byte{def.FixMap + 1, 0xc1, 0xc1},
			MethodAs: method,
		},
		{
			Name:     "Map16.ng.len",
			Data:     []byte{def.Map16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map16.ng.jump",
			Data:     []byte{def.Map16, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map16.ok",
			Data:     []byte{def.Map16, 0, 1, 0xc1, 0xc1},
			MethodAs: method,
		}, {
			Name:     "Map32.ng.len",
			Data:     []byte{def.Map32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map32.ng.jump",
			Data:     []byte{def.Map32, 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map32.ok",
			Data:     []byte{def.Map32, 0, 0, 0, 1, 0xc1, 0xc1},
			MethodAs: method,
		},
		{
			Name:     "Fixext1.ok",
			Data:     []byte{def.Fixext1, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Fixext2.ok",
			Data:     []byte{def.Fixext2, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Fixext4.ok",
			Data:     []byte{def.Fixext4, 0, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Fixext8.ok",
			Data:     []byte{def.Fixext8, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Fixext16.ok",
			Data:     []byte{def.Fixext16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Ext8.ng.size",
			Data:     []byte{def.Ext8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Ext8.ok",
			Data:     []byte{def.Ext8, 1, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Ext16.ng.size",
			Data:     []byte{def.Ext16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Ext16.ok",
			Data:     []byte{def.Ext16, 0, 1, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Ext32.ng.size",
			Data:     []byte{def.Ext32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Ext32.ok",
			Data:     []byte{def.Ext32, 0, 0, 0, 1, 0, 0},
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{0xc1},
			MethodAs: method,
		},
	}
	testcases.Run(t)
}
