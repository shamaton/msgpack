package decoding

import (
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
		tu.IsError(t, err, def.ErrNoData)
	})
	t.Run("not pointer", func(t *testing.T) {
		v := 0
		r := tu.NewTestReader([]byte{def.PositiveFixIntMax})
		err := Decode(r, v, false)
		tu.IsError(t, err, def.ErrReceiverNotPointer)
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
				Error:            def.ErrCanNotDecode,
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
				Name:             "error.struct",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.FixMap + 1, def.FixStr + 1, 'v'},
				ReadCount:        4,
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
				Name:             "error.map",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.FixMap + 1, def.FixStr + 1, 'v'},
				ReadCount:        4,
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
	t.Run("Complex64", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Fixext8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Fixext8,
				Data:             []byte{byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
				Expected:         true,
				ReadCount:        3,
				MethodAsWithCode: method,
			},
		}
		v := new(complex64)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, complex(1, 1))
	})
	t.Run("Complex128", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Fixext8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Fixext8,
				Data:             []byte{byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
				Expected:         true,
				ReadCount:        3,
				MethodAsWithCode: method,
			},
		}
		v := new(complex128)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, complex(1, 1))
	})

	t.Run("Array.nil", func(t *testing.T) {
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
		v := new([1]int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, [1]int{})
	})
	t.Run("Array.bin", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.bin",
				Code:             def.Bin8,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.len",
				Code:             def.Bin8,
				Data:             []byte{2, 1, 2},
				ReadCount:        2,
				Error:            def.ErrNotMatchArrayElement,
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
		v := new([1]byte)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, [1]byte{2})
	})
	t.Run("Array.string", func(t *testing.T) {
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
				Name:             "error.compare",
				Code:             def.Str8,
				Data:             []byte{2},
				ReadCount:        1,
				Error:            def.ErrNotMatchArrayElement,
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
		v := new([1]byte)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, [1]byte{'c'})
	})
	t.Run("Array.struct", func(t *testing.T) {
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
				Name:             "error.strlen",
				Code:             def.Array16,
				Data:             []byte{0, 2},
				ReadCount:        1,
				Error:            def.ErrNotMatchArrayElement,
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
				Name:             "error.struct",
				Code:             def.Array16,
				Data:             []byte{0, 1, def.FixMap + 1, def.FixStr + 1, 'v'},
				ReadCount:        4,
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
		v := new([1]st)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, [1]st{{V: 3}})
	})
	t.Run("Map.nil", func(t *testing.T) {
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
		v := new(map[string]int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, nil)
	})
	t.Run("Map.fixed", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.strlen",
				Code:             def.Map16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "error.map",
				Code:             def.Map16,
				Data:             []byte{0, 1},
				ReadCount:        1,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'a', def.PositiveFixIntMin + 3},
				Expected:         true,
				ReadCount:        4,
				MethodAsWithCode: method,
			},
		}
		v := new(map[string]int)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, map[string]int{"a": 3})
	})
	t.Run("Map.struct", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error.strlen",
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
				Name:             "error.value",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'a', def.FixMap + 1, def.FixStr + 1, 'v'},
				Expected:         true,
				ReadCount:        6,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'a', def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected:         true,
				ReadCount:        7,
				MethodAsWithCode: method,
			},
		}
		type st struct {
			V int `msgpack:"v"`
		}
		v := new(map[string]st)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, map[string]st{"a": {V: 3}})
	})
	t.Run("Struct", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Map16,
				Data:             []byte{},
				ReadCount:        0,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected:         true,
				ReadCount:        4,
				MethodAsWithCode: method,
			},
		}
		type st struct {
			V int `msgpack:"v"`
		}
		v := new(st)
		target = v
		testcases.Run(t)
		tu.Equal(t, *v, st{V: 3})
	})
	t.Run("Ptr.nil", func(t *testing.T) {
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
		v := new(int)
		target = &v
		testcases.Run(t)
		tu.Equal(t, *v, 0)
	})
	t.Run("Ptr", func(t *testing.T) {
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
				Data:             []byte{3},
				Expected:         true,
				ReadCount:        1,
				MethodAsWithCode: method,
			},
		}
		v := new(int)
		target = &v
		testcases.Run(t)
		tu.Equal(t, *v, 3)
	})
	t.Run("Interface.ptr", func(t *testing.T) {
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
				Data:             []byte{3},
				Expected:         true,
				ReadCount:        1,
				MethodAsWithCode: method,
			},
		}
		v := (any)(new(int))
		target = &v
		testcases.Run(t)
		vv := v.(*int)
		tu.Equal(t, *vv, 3)
	})
	t.Run("Interface", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:             "error",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'v', def.Int8},
				ReadCount:        4,
				Error:            io.EOF,
				MethodAsWithCode: method,
			},
			{
				Name:             "ok",
				Code:             def.Map16,
				Data:             []byte{0, 1, def.FixStr + 1, 'v', def.Int8, 3},
				Expected:         true,
				ReadCount:        5,
				MethodAsWithCode: method,
			},
		}
		type st struct {
			V any `msgpack:"v"`
		}
		v := new(st)
		target = v
		testcases.Run(t)
		var vv any = int8(3)
		tu.Equal(t, v.V, vv)
	})
}
