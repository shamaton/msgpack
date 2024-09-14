package decoding

import (
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

type AsXXXTestCase[T any] struct {
	Name           string
	Data           []byte
	Expected       T
	Error          error
	MethodAs       func(d *decoder) func(int, reflect.Kind) (T, int, error)
	MethodAsCustom func(d *decoder) (int, T, error)
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

	if tc.MethodAs == nil && tc.MethodAsCustom == nil {
		t.Fatal("must set either method or MethodAsCustom")
	}

	methodAs := func(d *decoder) (T, int, error) {
		if tc.MethodAs != nil {
			return tc.MethodAs(d)(0, kind)
		}
		if tc.MethodAsCustom != nil {
			v, o, err := tc.MethodAsCustom(d)
			return o, v, err
		}
		panic("unreachable")
	}

	t.Run(tc.Name, func(t *testing.T) {
		d := decoder{
			data: tc.Data,
		}

		v, offset, err := methodAs(&d)
		if tc.Error != nil {
			tu.IsError(t, err, tc.Error)
			return
		}
		tu.NoError(t, err)
		tu.Equal(t, v, tc.Expected)
		tu.Equal(t, offset, len(tc.Data))
	})
}

func TestDecoding(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		v := new(int)
		err := Decode(nil, v, false)
		tu.IsError(t, err, def.ErrNoData)
	})
	t.Run("not pointer", func(t *testing.T) {
		v := 0
		err := Decode([]byte{def.PositiveFixIntMax}, v, false)
		tu.IsError(t, err, def.ErrReceiverNotPointer)
	})
	t.Run("left data", func(t *testing.T) {
		v := new(int)
		err := Decode([]byte{def.PositiveFixIntMin, 0}, v, false)
		tu.IsError(t, err, def.ErrHasLeftOver)
	})
}

func Test_decodeWithCode(t *testing.T) {
	var target any
	method := func(d *decoder) func(offset int, _ reflect.Kind) (bool, int, error) {
		return func(offset int, _ reflect.Kind) (bool, int, error) {
			rv := reflect.ValueOf(target)
			o, err := d.decode(rv.Elem(), offset)
			return true, o, err
		}
	}

	t.Run("Int", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:     "error",
				Data:     []byte{def.Int8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Int8, 5},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Uint8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Uint8, 5},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Float32},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Float32, 63, 128, 0, 0},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Float64},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Float64, 63, 240, 0, 0, 0, 0, 0, 0},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Bin8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Bin8, 1, 'a'},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Str8},
				Expected: false,
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Str8, 1, 'b'},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Int8},
				Error:    def.ErrCanNotDecode,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.True},
				Expected: true,
				MethodAs: method,
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
				Name:     "ok",
				Data:     []byte{def.Nil},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Bin8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Bin8, 1, 2},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Str8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.bytelen",
				Data:     []byte{def.Str8, 1},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Str8, 1, 'c'},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Array16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.require",
				Data:     []byte{def.Array16, 0, 1},
				Error:    def.ErrLackDataLengthToSlice,
				MethodAs: method,
			},
			{
				Name:     "error.slice",
				Data:     []byte{def.Array16, 0, 1, def.Int8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Array16, 0, 1, def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Array16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.require",
				Data:     []byte{def.Array16, 0, 1},
				Error:    def.ErrLackDataLengthToSlice,
				MethodAs: method,
			},
			{
				Name:     "error.slice",
				Data:     []byte{def.Array16, 0, 1, def.Map16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.struct",
				Data:     []byte{def.Array16, 0, 1, def.FixMap + 1, def.FixStr + 1, 'v'},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Array16, 0, 1, def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Array16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.require",
				Data:     []byte{def.Array16, 0, 1},
				Error:    def.ErrLackDataLengthToSlice,
				MethodAs: method,
			},
			{
				Name:     "error.slice",
				Data:     []byte{def.Array16, 0, 1, def.Map16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.map",
				Data:     []byte{def.Array16, 0, 1, def.FixMap + 1, def.FixStr + 1, 'v'},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Array16, 0, 1, def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Fixext8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Fixext8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Fixext8, byte(def.ComplexTypeCode()), 63, 128, 0, 0, 63, 128, 0, 0},
				Expected: true,
				MethodAs: method,
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
				Name:     "ok",
				Data:     []byte{def.Nil},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.bin",
				Data:     []byte{def.Bin8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.len",
				Data:     []byte{def.Bin8, 2, 1, 2},
				Error:    def.ErrNotMatchArrayElement,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Bin8, 1, 2},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Str8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.compare",
				Data:     []byte{def.Str8, 2},
				Error:    def.ErrNotMatchArrayElement,
				MethodAs: method,
			},
			{
				Name:     "error.bytelen",
				Data:     []byte{def.Str8, 1},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Str8, 1, 'c'},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Array16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.len.match",
				Data:     []byte{def.Array16, 0, 2},
				Error:    def.ErrNotMatchArrayElement,
				MethodAs: method,
			},
			{
				Name:     "error.slice",
				Data:     []byte{def.Array16, 0, 1},
				Error:    def.ErrLackDataLengthToSlice,
				MethodAs: method,
			},
			{
				Name:     "error.struct",
				Data:     []byte{def.Array16, 0, 1, def.FixMap + 1, def.FixStr + 1, 'v'},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Array16, 0, 1, def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "ok",
				Data:     []byte{def.Nil},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Map16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.require",
				Data:     []byte{def.Map16, 0, 1},
				Error:    def.ErrLackDataLengthToMap,
				MethodAs: method,
			},
			{
				Name:     "error.map",
				Data:     []byte{def.Map16, 0, 1, def.Str16, 0},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'a', def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error.strlen",
				Data:     []byte{def.Map16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "error.require",
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
				Name:     "error.value",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'a', def.FixMap + 1, def.FixStr + 1, 'v'},
				Expected: true,
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'a', def.FixMap + 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Map16},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'v', def.PositiveFixIntMin + 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "ok",
				Data:     []byte{def.Nil},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Int8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Int8, 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Int8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Int8, 3},
				Expected: true,
				MethodAs: method,
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
				Name:     "error",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'v', def.Int8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'v', def.Int8, 3},
				Expected: true,
				MethodAs: method,
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

	t.Run("StructWithCustomDecodeMethod", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:     "error",
				Data:     []byte{def.Int8},
				Error:    def.ErrTooShortBytes,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Int8, 3},
				Expected: true,
				MethodAs: method,
			},
		}
		v := new(StructWithCustomDecodeMethod)
		target = &v
		testcases.Run(t)
		tu.Equal(t, v.A, 9)
		tu.Equal(t, v.B, "Happy birthday!")
	})

	t.Run("CustomDecodeNetAddrWhitelist", func(t *testing.T) {
		testcases := AsXXXTestCases[bool]{
			{
				Name:     "error",
				Data:     []byte{def.Int8, 0},
				Error:    def.ErrCanNotDecode,
				MethodAs: method,
			},
			{
				Name:     "ok",
				Data:     []byte{def.Str8, 24, '1', '9', '2', '.', '1', '6', '8', '.', '1', '.', '0', '/', '2', '4', ',', '0', '.', '0', '.', '0', '.', '0', '/', '0'},
				Expected: true,
				MethodAs: method,
			},
		}
		v := new(NetAddrWhitelist)
		target = &v
		testcases.Run(t)
		tu.Equal(t, len(*v), 2)
		tu.Equal(t, (*v)[0].String(), "192.168.1.0/24")
		tu.Equal(t, (*v)[1].String(), "0.0.0.0/0")
	})

}

type StructWithCustomDecodeMethod struct {
	A int
	B string
}

func (s *StructWithCustomDecodeMethod) UnmarshalMsgpack(value any) error {
	if v, ok := value.(int8); ok {
		s.A = int(v) * 3
		s.B = "Happy birthday!"
		return nil
	}
	return def.ErrCanNotDecode
}

type NetAddrWhitelist []net.IPNet

func (a *NetAddrWhitelist) UnmarshalMsgpack(value any) error {
	if v, ok := value.(string); ok {
		nets := strings.Split(v, ",")
		for _, netw := range nets {
			_, ipnet, err := net.ParseCIDR(netw)
			if err != nil {
				return err
			}
			*a = append(*a, *ipnet)
		}
		return nil
	}
	return def.ErrCanNotDecode
}
