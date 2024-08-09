package decoding

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"

	"github.com/shamaton/msgpack/v2/internal/common"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

type AsXXXTestCase[T any] struct {
	Name             string
	Code             byte
	Data             []byte
	ReadCount        int
	Expected         T
	Error            error
	IsTemplateError  bool
	MethodAs         func(d *decoder) func(reflect.Kind) (T, error)
	MethodAsWithCode func(d *decoder) func(byte, reflect.Kind) (T, error)
	MethodAsCustom   func(d *decoder) (T, error)
}

type AsXXXTestCases[T any] []AsXXXTestCase[T]

func (tcs AsXXXTestCases[T]) Run(t *testing.T) {
	for _, tc := range tcs {
		tc.Run(t)
	}
}

func (tc *AsXXXTestCase[T]) Run(t *testing.T) {
	const kind = reflect.String
	t.Helper()

	if tc.MethodAs == nil && tc.MethodAsWithCode == nil && tc.MethodAsCustom == nil {
		t.Fatal("must set either method or methodAsWithCode or MethodAsCustom")
	}

	methodAs := func(d *decoder) (T, error) {
		if tc.MethodAs != nil {
			return tc.MethodAs(d)(kind)
		}
		if tc.MethodAsWithCode != nil {
			return tc.MethodAsWithCode(d)(tc.Code, kind)
		}
		if tc.MethodAsCustom != nil {
			return tc.MethodAsCustom(d)
		}
		panic("unreachable")
	}

	t.Run(tc.Name, func(t *testing.T) {
		r := tu.NewTestReader(tc.Data)
		d := decoder{
			r:   r,
			buf: common.GetBuffer(),
		}
		defer common.PutBuffer(d.buf)

		v, err := methodAs(&d)
		tu.Equal(t, r.Count(), tc.ReadCount)

		if tc.Error != nil {
			tu.IsError(t, err, tc.Error)
			return
		}
		if tc.IsTemplateError {
			tu.ErrorContains(t, err, fmt.Sprintf("msgpack : invalid code %x", tc.Code))
			return
		}
		tu.NoError(t, err)
		tu.Equal(t, v, tc.Expected)

		p := make([]byte, 1)
		n, err := d.r.Read(p)
		tu.IsError(t, err, io.EOF)
		tu.Equal(t, n, 0)
	})
}

func TestDecoding(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		v := new(int)
		err := Decode(nil, v, false)
		tu.Error(t, err)
		tu.Equal(t, err.Error(), "reader is nil")
	})
}

func Test_decodeWithCode(t *testing.T) {

	var target any
	method := func(d *decoder) func(code byte, _ reflect.Kind) (bool, error) {
		return func(code byte, _ reflect.Kind) (bool, error) {
			rv := reflect.ValueOf(target)
			return true, d.decodeWithCode(code, rv.Elem())
		}
	}

	t.Run("Int", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Int8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Int8,
				Data:             []byte{5},
				Expected:         true,
				ReadCount:        1,
				MethodAsWithCode: method,
			},
		}
		v := new(int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, 5)
	})
	t.Run("Uint", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Uint8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Uint8,
				Data:             []byte{5},
				Expected:         true,
				ReadCount:        1,
				MethodAsWithCode: method,
			},
		}
		v := new(uint)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, 5)
	})
	t.Run("Float32", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Float32,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Float32,
				Data:             []byte{63, 128, 0, 0},
				Expected:         true,
				ReadCount:        1,
				MethodAsWithCode: method,
			},
		}
		v := new(float32)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, 1)
	})
	t.Run("Float64", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Float64,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Float64,
				Data:             []byte{63, 240, 0, 0, 0, 0, 0, 0},
				Expected:         true,
				ReadCount:        1,
				MethodAsWithCode: method,
			},
		}
		v := new(float64)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, 1)
	})
	t.Run("BinString", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Bin8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Bin8,
				Data:             []byte{1, 'a'},
				Expected:         true,
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		v := new(string)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, "a")
	})
	t.Run("String", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Str8,
				Data:             []byte{},
				Expected:         false,
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Str8,
				Data:             []byte{1, 'b'},
				Expected:         true,
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		v := new(string)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, "b")
	})
	t.Run("Bool", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Int8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            ErrCanNotDecode,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.True,
				Data:             []byte{},
				Expected:         true,
				ReadCount:        0,
				MethodAsWithCode: method,
			},
		}
		v := new(bool)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, true)
	})
	t.Run("Slice.nil", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "ok",
				Code:             def.Nil,
				Data:             []byte{},
				Expected:         true,
				ReadCount:        0,
				MethodAsWithCode: method,
			},
		}
		v := new([]int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, nil)
	})
	t.Run("Slice.bin", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Bin8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Bin8,
				Data:             []byte{1, 2},
				Expected:         true,
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		v := new([]byte)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, []byte{2})
	})
	t.Run("Slice.string", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.strlen",
				Code:             def.Str8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.bytelen",
				Code:             def.Str8,
				Data:             []byte{1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Str8,
				Data:             []byte{1, 'c'},
				Expected:         true,
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		v := new([]byte)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, []byte{'c'})
	})
	t.Run("Slice.fixed", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.strlen",
				Code:             def.Array16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.slice",
				Code:             def.Array16,
				Data:             []byte{0, 1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.PositiveFixIntMin + 3},
				Expected:         true,
				ReadCount:        2,
				MethodAsWithCode: method,
			},
		}
		v := new([]int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, []int{3})
	})
	t.Run("Slice.struct", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.strlen",
				Code:             def.Array16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.slice",
				Code:             def.Array16,
				Data:             []byte{0, 1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected:         true,
				ReadCount:        5,
				MethodAsWithCode: method,
			},
		}
		type st struct {
			V int `msgpack:"v"`
		}
		v := new([]st)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, []st{{V: 3}})
	})
	t.Run("Slice.map", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.strlen",
				Code:             def.Array16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.slice",
				Code:             def.Array16,
				Data:             []byte{0, 1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected:         true,
				ReadCount:        5,
				MethodAsWithCode: method,
			},
		}
		v := new([]map[string]int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, []map[string]int{{"v": 3}})
	})
}
