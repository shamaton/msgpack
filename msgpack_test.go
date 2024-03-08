package msgpack_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v2"
	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	extTime "github.com/shamaton/msgpack/v2/time"
)

var now time.Time

func init() {
	n := time.Now()
	now = time.Unix(n.Unix(), int64(n.Nanosecond()))
}

func TestInt(t *testing.T) {
	{
		args := []encdecArg[int]{{
			n: "FixInt",
			v: -8,
			c: func(d []byte) bool {
				return def.NegativeFixintMin <= int8(d[0]) && int8(d[0]) <= def.NegativeFixintMax
			}}, {
			n: "Int8",
			v: -108,
			c: func(d []byte) bool {
				return d[0] == def.Int8
			}}, {
			n: "Int16",
			v: -30108,
			c: func(d []byte) bool {
				return d[0] == def.Int16
			}}, {
			n: "Int32",
			v: -1030108,
			c: func(d []byte) bool {
				return d[0] == def.Int32
			}},
		}
		encdec(t, args...)
	}
	{
		arg := encdecArg[int64]{
			n: "Int64",
			v: int64(math.MinInt64 + 12345),
			c: func(d []byte) bool {
				return d[0] == def.Int64
			},
		}
		encdec(t, arg)
	}

	// error
	{
		arg := encdecArg[uint8]{
			n: "ErrorDecToUint8",
			v: -8,
			c: func(d []byte) bool {
				return def.NegativeFixintMin <= int8(d[0]) && int8(d[0]) <= def.NegativeFixintMax
			},
			e: "value different",
		}
		encdec(t, arg)
	}
	{
		arg := encdecArg[int32]{
			n: "ErrorDecToInt32",
			v: int64(math.MinInt64 + 12345),
			c: func(d []byte) bool {
				return d[0] == def.Int64
			},
			e: "value different",
		}
		encdec(t, arg)
	}
}

func TestUint(t *testing.T) {
	{
		args := []encdecArg[uint]{{
			n: "FixUint",
			v: uint(8),
			c: func(d []byte) bool {
				return def.PositiveFixIntMin <= uint8(d[0]) && uint8(d[0]) <= def.PositiveFixIntMax
			}}, {
			n: "Uint8",
			v: uint(130),
			c: func(d []byte) bool {
				return d[0] == def.Uint8
			}}, {
			n: "Uint16",
			v: uint(30130),
			c: func(d []byte) bool {
				return d[0] == def.Uint16
			}}, {
			n: "Uint32",
			v: uint(1030130),
			c: func(d []byte) bool {
				return d[0] == def.Uint32
			}},
		}
		encdec(t, args...)
	}
	{
		arg := encdecArg[uint64]{
			n: "Uint64",
			v: uint64(math.MaxUint64 - 12345),
			c: func(d []byte) bool {
				return d[0] == def.Uint64
			},
		}
		encdec(t, arg)
	}
}
func TestFloat(t *testing.T) {
	t.Run("Float32", func(t *testing.T) {
		c := func(d []byte) bool {
			return d[0] == def.Float32
		}
		args := []encdecArg[float32]{{
			n: "0",
			v: float32(0),
			c: c,
		}, {
			n: "-1",
			v: float32(-1),
			c: c,
		}, {
			n: "SmallestNonzeroFloat",
			v: float32(math.SmallestNonzeroFloat32),
			c: c,
		}, {
			n: "MaxFloat",
			v: float32(math.MaxFloat32),
			c: c,
		}}
		encdec(t, args...)
	})

	t.Run("Float64", func(t *testing.T) {
		c := func(d []byte) bool {
			return d[0] == def.Float64
		}
		args := []encdecArg[float64]{{
			n: "0",
			v: float64(0),
			c: c,
		}, {
			n: "-1",
			v: float64(-1),
			c: c,
		}, {
			n: "SmallestNonzeroFloat",
			v: math.SmallestNonzeroFloat64,
			c: c,
		}, {
			n: "MaxFloat",
			v: math.MaxFloat64,
			c: c,
		}}
		encdec(t, args...)
	})

	t.Run("FloatToInt", func(t *testing.T) {
		args := []encdecArg[int]{
			{
				n: "FromFloat32",
				v: float32(2.345),
				vc: func(v int) error {
					if v != 2 {
						return fmt.Errorf("different value: %d", v)
					}
					return nil
				},
				skipEq: true,
			},
			{
				n: "FromFloat64",
				v: float64(6.789),
				vc: func(v int) error {
					if v != 6 {
						return fmt.Errorf("different value: %d", v)
					}
					return nil
				},
				skipEq: true,
			},
		}
		encdec(t, args...)
	})

	// error
	t.Run("Float32ToFloat64", func(t *testing.T) {
		arg := encdecArg[float64]{
			n: "NotEqual",
			v: float32(math.MaxFloat32),
			c: func(d []byte) bool {
				return d[0] == def.Float32
			},
			e: "value different",
		}
		encdec(t, arg)
	})
	t.Run("Float64ToFloat32", func(t *testing.T) {
		arg := encdecArg[float32]{
			n: "ErrorDecToFloat32",
			v: math.MaxFloat64,
			c: func(d []byte) bool {
				return d[0] == def.Float64
			},
			e: "invalid code cb decoding",
		}
		encdec(t, arg)
	})
	t.Run("Float64ToString", func(t *testing.T) {
		arg := encdecArg[string]{
			n: "ErrorDecToString",
			v: math.MaxFloat64,
			c: func(d []byte) bool {
				return d[0] == def.Float64
			},
			e: "invalid code cb decoding",
		}
		encdec(t, arg)
	})
}

func TestBool(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		args := []encdecArg[bool]{
			{
				n: "True",
				v: true,
				c: func(d []byte) bool {
					return d[0] == def.True
				},
			},
			{
				n: "False",
				v: false,
				c: func(d []byte) bool {
					return d[0] == def.False
				},
			},
		}
		encdec(t, args...)
	})

	// error
	t.Run("BoolToUint8", func(t *testing.T) {
		arg := encdecArg[uint8]{
			n: "ErrorDecToUint8",
			v: true,
			c: func(d []byte) bool {
				return d[0] == def.True
			},
			e: "invalid code c3 decoding",
		}
		encdec(t, arg)
	})
}

func TestNil(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		args := []encdecArg[*map[any]any]{
			{
				n: "Nil",
				v: nil,
				c: func(d []byte) bool {
					return d[0] == def.Nil
				},
				vc: func(v *map[any]any) error {
					if v != nil {
						return fmt.Errorf("not nil: %v", v)
					}
					return nil
				},
			},
		}
		encdec(t, args...)
	})
}

func TestString(t *testing.T) {
	// len 31
	const base = "abcdefghijklmnopqrstuvwxyz12345"

	t.Run("String", func(t *testing.T) {
		args := []encdecArg[string]{
			{
				n: "EmptyString",
				v: "",
				c: func(d []byte) bool {
					return def.FixStr <= d[0] && d[0] < def.FixStr+32
				},
			},
			{
				n: "FixStr",
				v: strings.Repeat(base, 1),
				c: func(d []byte) bool {
					return def.FixStr <= d[0] && d[0] < def.FixStr+32
				},
			},
			{
				n: "Str8",
				v: strings.Repeat(base, 8),
				c: func(d []byte) bool {
					return d[0] == def.Str8
				},
			},
			{
				n: "Str16",
				v: strings.Repeat(base, (math.MaxUint16/len(base))-1),
				c: func(d []byte) bool {
					return d[0] == def.Str16
				},
			},
			{
				n: "Str32",
				v: strings.Repeat(base, (math.MaxUint16/len(base))+1),
				c: func(d []byte) bool {
					return d[0] == def.Str32
				},
			},
		}
		encdec(t, args...)
	})

	// type different
	t.Run("StringToBytes", func(t *testing.T) {
		arg := encdecArg[[]byte]{
			n: "TypeDifferentButSameValue",
			v: strings.Repeat(base, 8),
			c: func(d []byte) bool {
				return d[0] == def.Str8
			},
			e: "value different",
		}
		encdec(t, arg)
	})
	t.Run("BytesToString", func(t *testing.T) {
		arg := encdecArg[string]{
			n: "TypeDifferentButSameValue",
			v: []byte(base),
			c: func(d []byte) bool {
				return d[0] == def.Bin8
			},
			e: "value different",
		}
		encdec(t, arg)
	})
}

func TestComplex(t *testing.T) {
	t.Run("Complex64", func(t *testing.T) {
		args := []encdecArg[complex64]{
			{
				n: "To64",
				v: complex64(complex(1, 2)),
				c: func(d []byte) bool {
					return d[0] == def.Fixext8 && int8(d[1]) == def.ComplexTypeCode()
				},
			},
			{
				n: "Nil",
				v: nil,
				e: "should not reach this line",
			},
		}
		encdec(t, args...)

		args2 := []encdecArg[complex128]{
			{
				n: "To128",
				v: complex64(complex(1, 2)),
				c: func(d []byte) bool {
					return d[0] == def.Fixext8 && int8(d[1]) == def.ComplexTypeCode()
				},
				vc: func(t complex128) error {
					if imag(t) == 0 || real(t) == 0 {
						return fmt.Errorf("somthing wrong %v", t)
					}
					return nil
				},
				skipEq: true,
			},
		}
		encdec(t, args2...)
	})

	t.Run("Complex128", func(t *testing.T) {
		args := []encdecArg[complex128]{
			{
				n: "To128",
				v: complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
				c: func(d []byte) bool {
					return d[0] == def.Fixext16 && int8(d[1]) == def.ComplexTypeCode()
				},
			},
			{
				n: "Nil",
				v: nil,
				e: "should not reach this line",
			},
		}
		encdec(t, args...)

		args2 := []encdecArg[complex64]{
			{
				n: "To64",
				v: complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
				c: func(d []byte) bool {
					return d[0] == def.Fixext16 && int8(d[1]) == def.ComplexTypeCode()
				},
				vc: func(t complex64) error {
					if imag(t) != 0 || real(t) == 0 {
						return fmt.Errorf("somthing wrong %v", t)
					}
					return nil
				},
				skipEq: true,
			},
		}
		encdec(t, args2...)
	})

	t.Run("ComplexTypeCode", func(t *testing.T) {
		b64, err := msgpack.Marshal(complex64(complex(3, 4)))
		NoError(t, err)

		b128, err := msgpack.Marshal(complex128(complex(5, 6)))
		NoError(t, err)

		// change complex type code
		msgpack.SetComplexTypeCode(int8(-99))

		data := []struct {
			name   string
			b      []byte
			errStr string
		}{
			{"From64", b64, "fixext8"},
			{"From128", b128, "fixext16"},
		}

		var r64 complex64
		var r128 complex128
		for _, d := range data {
			for _, u := range unmarshallers {
				t.Run(d.name+"To64"+u.name, func(t *testing.T) {
					ErrorContains(t, u.u(d.b, &r64), d.errStr)
				})
				t.Run(d.name+"To128"+u.name, func(t *testing.T) {
					ErrorContains(t, u.u(d.b, &r128), d.errStr)
				})
			}
		}
	})
}

func TestAny(t *testing.T) {
	f := func(v interface{}) error {
		b, err := msgpack.Marshal(v)
		if err != nil {
			return err
		}
		var r interface{}
		err = msgpack.Unmarshal(b, &r)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", r) {
			return fmt.Errorf("different value %v, %v", v, r)
		}
		return err
	}

	a1 := make([]int, math.MaxUint16)
	a2 := make([]int, math.MaxUint16+1)
	m1 := map[string]int{}
	m2 := map[string]int{}

	for i := range a1 {
		a1[i] = i
		m1[fmt.Sprint(i)] = 1
	}
	for i := range a2 {
		a2[i] = i
		m2[fmt.Sprint(i)] = 1
	}

	vars := []interface{}{
		true, false,
		1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint32 + 1,
		-1, math.MinInt8, math.MinInt16, math.MinInt32, math.MinInt32 - 1,
		math.MaxFloat32, math.MaxFloat64,
		"a", strings.Repeat("a", math.MaxUint8), strings.Repeat("a", math.MaxUint16), strings.Repeat("a", math.MaxUint16+1),
		[]byte(strings.Repeat("a", math.MaxUint8)),
		[]byte(strings.Repeat("a", math.MaxUint16)),
		[]byte(strings.Repeat("a", math.MaxUint16+1)),
		[]interface{}{1, "a", 1.23}, a1, a2,
		map[interface{}]interface{}{"1": 1, 1.23: "a"}, m1, m2,
		time.Unix(now.Unix(), int64(now.Nanosecond())),
	}

	for i, v := range vars {
		if err := f(v); err != nil {
			t.Error(i, err)
		}
	}

	t.Run("Any", func(t *testing.T) {
		args := []encdecArg[any]{
			{n: "true", v: any(true)},
			{n: "false", v: any(false)},
			{n: "1", v: any(1)},
			{n: "MaxUint8", v: any(math.MaxUint8)},
			{n: "MaxUint16", v: any(math.MaxUint16)},
			{n: "MaxUint32", v: any(math.MaxUint32)},
			{n: "MaxUint32+1", v: any(math.MaxUint32 + 1)},
			{n: "-1", v: any(-1)},
			{n: "MinInt8", v: any(math.MinInt8)},
			{n: "MinInt16", v: any(math.MinInt16)},
			{n: "MinInt32", v: any(math.MinInt32)},
			{n: "MinInt32-1", v: any(math.MinInt32 - 1)},
			{n: "MaxFloat32", v: any(math.MaxFloat32)},
			{n: "MaxFloat64", v: any(math.MaxFloat64)},
			{n: "Str1", v: any("a")},
			{n: "Str255", v: any(strings.Repeat("a", math.MaxUint8))},
			{n: "Str65535", v: any(strings.Repeat("a", math.MaxUint16))},
			{n: "Str65536", v: any(strings.Repeat("a", math.MaxUint16+1))},
			{n: "Bin255", v: any([]byte(strings.Repeat("a", math.MaxUint8)))},
			{n: "Bin65535", v: any([]byte(strings.Repeat("a", math.MaxUint16)))},
			{n: "Bin65536", v: any([]byte(strings.Repeat("a", math.MaxUint16+1)))},
			{n: "Slice3", v: any([]any{1, "a", 1.23})},
			{n: "Slice65535", v: any(a1)},
			{n: "Slice65536", v: any(a2)},
			{n: "Map3", v: any(map[any]any{"1": 1, 1.23: "a"})},
			{n: "Map65535", v: any(m1)},
			{n: "Map65536", v: any(m2)},
			{n: "Time", v: any(time.Unix(now.Unix(), int64(now.Nanosecond())))},
		}
		for i := range args {
			i := i
			args[i].skipEq = true
			args[i].vc = func(t any) error {
				if fmt.Sprintf("%v", args[i].v) != fmt.Sprintf("%v", t) {
					return fmt.Errorf("different value %v, %v", args[i].v, t)
				}
				return nil
			}
		}
		encdec(t, args...)
	})

	// error
	t.Run("AnyError", func(t *testing.T) {
		var r any
		err := msgpack.Unmarshal([]byte{def.Ext32}, &r)
		ErrorContains(t, err, "invalid code")
	})
}

func TestBin(t *testing.T) {
	makeByteSlice := func(len int) []byte {
		v := make([]byte, len)
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		return v
	}

	t.Run("Slice", func(t *testing.T) {
		args := []encdecArg[[]byte]{
			{
				n: "Bin8",
				v: makeByteSlice(128),
				c: func(d []byte) bool {
					return d[0] == def.Bin8
				},
			},
			{
				n: "Bin16",
				v: makeByteSlice(31280),
				c: func(d []byte) bool {
					return d[0] == def.Bin16
				},
			},
			{
				n: "Bin32",
				v: makeByteSlice(1031280),
				c: func(d []byte) bool {
					return d[0] == def.Bin32
				},
			},
		}
		encdec(t, args...)
	})

	t.Run("Array", func(t *testing.T) {
		var (
			a128     [128]byte
			a31280   [31280]byte
			a1031280 [1031280]byte
		)
		for i := range a128 {
			a128[i] = byte(rand.Intn(0xff))
		}
		for i := range a31280 {
			a31280[i] = byte(rand.Intn(0xff))
		}
		for i := range a1031280 {
			a1031280[i] = byte(rand.Intn(0xff))
		}
		args128 := []encdecArg[[128]byte]{
			{
				n: "Bin8",
				v: a128,
				c: func(d []byte) bool {
					return d[0] == def.Bin8
				},
			},
		}
		encdec(t, args128...)

		args31280 := []encdecArg[[31280]byte]{
			{
				n: "Bin16",
				v: a31280,
				c: func(d []byte) bool {
					return d[0] == def.Bin16
				},
			},
		}
		encdec(t, args31280...)

		args1031280 := []encdecArg[[1031280]byte]{
			{
				n: "Bin32",
				v: a1031280,
				c: func(d []byte) bool {
					return d[0] == def.Bin32
				},
			},
		}
		encdec(t, args1031280...)

		args := []encdecArg[[1]byte]{
			{
				n: "Nil",
				v: nil,
				c: func(d []byte) bool {
					return d[0] == def.Nil
				},
				e: "value different",
			},
		}
		encdec(t, args...)
	})

	t.Run("Error", func(t *testing.T) {
		args1 := []encdecArg[[1]byte]{
			{
				n: "Nil",
				v: nil,
				c: func(d []byte) bool {
					return d[0] == def.Nil
				},
				e: "value different",
			},
		}
		encdec(t, args1...)

		var a128 [128]byte
		for i := range a128 {
			a128[i] = byte(rand.Intn(0xff))
		}
		args2 := []encdecArg[[127]byte]{
			{
				n: "Len",
				v: a128,
				c: func(d []byte) bool {
					return d[0] == def.Bin8
				},
				e: "[127]uint8 len is 127, but msgpack has 128 elements",
			},
		}
		encdec(t, args2...)
	})
}

func TestArray(t *testing.T) {
	// slice
	{
		var v, r []int
		v = nil
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.Nil == code
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int
		v = make([]int, 15)
		for i := range v {
			v[i] = rand.Intn(math.MaxInt32)
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int
		v = make([]int, 30015)
		for i := range v {
			v[i] = rand.Intn(math.MaxInt32)
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Array16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int
		v = make([]int, 1030015)
		for i := range v {
			v[i] = rand.Intn(math.MaxInt32)
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Array32
		}); err != nil {
			t.Error(err)
		}
	}
	// array
	{
		var v, r [8]float32
		for i := range v {
			v[i] = float32(rand.Intn(0xff))
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil && strings.Contains(err.Error(), "value different") {
			for i := range v {
				if v[i] != r[i] {
					t.Error("value different")
				}
			}
		} else if err != nil {
			t.Error(err)
		}
	}
	{
		var v, r [31280]string
		for i := range v {
			v[i] = "a"
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Array16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r [1031280]bool
		for i := range v {
			v[i] = rand.Intn(0xff) > 0x7f
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Array32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v []int
		var r [1]int
		v = nil
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Nil
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
	}
	{
		v := "abcde"
		var r [5]byte
		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		err = msgpack.Unmarshal(b, &r)
		if err != nil {
			t.Error(err)
		}
		if v != string(r[:]) {
			t.Errorf("value different %v, %v", v, string(r[:]))
		}
	}
}

func TestFixedSlice(t *testing.T) {
	{
		var v, r []int
		v = []int{-1, 1}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint
		v = []uint{0, 100}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int8
		v = []int8{math.MinInt8, math.MaxInt8}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int16
		v = []int16{math.MinInt16, math.MaxInt16}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int32
		v = []int32{math.MinInt32, math.MaxInt32}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int64
		v = []int64{math.MinInt64, math.MaxInt64}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		// byte array
		var v, r []uint8
		v = []uint8{0, math.MaxUint8}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.Bin8 == code
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint16
		v = []uint16{0, math.MaxUint16}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint32
		v = []uint32{0, math.MaxUint32}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint64
		v = []uint64{0, math.MaxUint64}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []float32
		v = []float32{math.SmallestNonzeroFloat32, math.MaxFloat32}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []float64
		v = []float64{math.SmallestNonzeroFloat64, math.MaxFloat64}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []string
		v = []string{"aaa", "bbb"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []bool
		v = []bool{true, false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestFixedMap(t *testing.T) {
	{
		var v, r map[string]int
		v = map[string]int{"a": 1, "b": 2}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r map[string]uint
		v = map[string]uint{"a": math.MaxUint32, "b": 0}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]string
		v = map[string]string{"a": "12345", "abcdefghijklmnopqrstuvwxyz": "abcdefghijklmnopqrstuvwxyz"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]float32
		v = map[string]float32{"a": math.MaxFloat32, "b": math.SmallestNonzeroFloat32}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]float64
		v = map[string]float64{"a": math.MaxFloat64, "b": math.SmallestNonzeroFloat64}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]bool
		v = map[string]bool{"a": true, "b": false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int8
		v = map[string]int8{"a": math.MinInt8, "b": math.MaxInt8}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int16
		v = map[string]int16{"a": math.MaxInt16, "b": math.MinInt16}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int32
		v = map[string]int32{"a": math.MaxInt32, "b": math.MinInt32}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int64
		v = map[string]int64{"a": math.MinInt64, "b": math.MaxInt64}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint8
		v = map[string]uint8{"a": 0, "b": math.MaxUint8}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint16
		v = map[string]uint16{"a": 0, "b": math.MaxUint16}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint32
		v = map[string]uint32{"a": 0, "b": math.MaxUint32}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint64
		v = map[string]uint64{"a": 0, "b": math.MaxUint64}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int]string
		v = map[int]string{0: "a", 1: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int]bool
		v = map[int]bool{1: true, 2: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint]string
		v = map[uint]string{0: "a", 1: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint]bool
		v = map[uint]bool{0: true, 255: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float32]string
		v = map[float32]string{math.MaxFloat32: "a", math.SmallestNonzeroFloat32: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float32]bool
		v = map[float32]bool{math.SmallestNonzeroFloat32: true, math.MaxFloat32: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float64]string
		v = map[float64]string{math.MaxFloat64: "a", math.SmallestNonzeroFloat64: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float64]bool
		v = map[float64]bool{math.SmallestNonzeroFloat64: true, math.MaxFloat64: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int8]string
		v = map[int8]string{math.MinInt8: "a", math.MaxInt8: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int8]bool
		v = map[int8]bool{math.MinInt8: true, math.MaxInt8: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int16]string
		v = map[int16]string{math.MaxInt16: "a", math.MinInt16: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int16]bool
		v = map[int16]bool{math.MaxInt16: true, math.MinInt16: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int32]string
		v = map[int32]string{math.MinInt32: "a", math.MaxInt32: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int32]bool
		v = map[int32]bool{math.MinInt32: true, math.MaxInt32: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int64]string
		v = map[int64]string{math.MaxInt64: "a", math.MinInt64: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int64]bool
		v = map[int64]bool{math.MaxInt64: true, math.MinInt64: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint8]string
		v = map[uint8]string{0: "a", math.MaxUint8: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint8]bool
		v = map[uint8]bool{0: true, math.MaxUint8: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint16]string
		v = map[uint16]string{0: "a", math.MaxUint16: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint16]bool
		v = map[uint16]bool{0: true, math.MaxUint16: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint32]string
		v = map[uint32]string{0: "a", math.MaxUint32: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint32]bool
		v = map[uint32]bool{0: true, math.MaxUint32: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint64]string
		v = map[uint64]string{0: "a", math.MaxUint64: "b"}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint64]bool
		v = map[uint64]bool{0: true, math.MaxUint64: false}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

}

func TestTime(t *testing.T) {
	{
		var v, r time.Time
		v = time.Unix(now.Unix(), 0)
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Fixext4
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r time.Time
		v = time.Unix(now.Unix(), int64(now.Nanosecond()))
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Fixext8
		}); err != nil {
			t.Error(err)
		}
	}

	{
		now := time.Now().Unix()
		nowByte := make([]byte, 4)
		binary.BigEndian.PutUint32(nowByte, uint32(now))

		var r time.Time
		c := def.TimeStamp

		nanoByte := make([]byte, 64)
		for i := range nanoByte[:30] {
			nanoByte[i] = 0xff
		}

		b := append([]byte{def.Fixext8, byte(c)}, nanoByte...)
		err := msgpack.UnmarshalAsArray(b, &r)
		if err == nil || !strings.Contains(err.Error(), "In timestamp 64 formats") {
			t.Error(err)
		}

		nanoByte = make([]byte, 96)
		for i := range nanoByte[:32] {
			nanoByte[i] = 0xff
		}
		b = append([]byte{def.Ext8, byte(12), byte(c)}, nanoByte...)
		err = msgpack.UnmarshalAsArray(b, &r)
		if err == nil || !strings.Contains(err.Error(), "In timestamp 96 formats") {
			t.Error(err)
		}

		notReach := []byte{def.Fixext1}
		_, _, err = extTime.Decoder.AsValue(0, reflect.Bool, &notReach)
		if err == nil || !strings.Contains(err.Error(), "should not reach this line") {
			t.Error("something wrong", err)
		}
	}
}

func TestMap(t *testing.T) {
	{
		var v, r map[int]int
		v = map[int]int{1: 2, 3: 4, 5: 6, 7: 8, 9: 10}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r map[int]int
		v = make(map[int]int, 1000)
		for i := 0; i < 1000; i++ {
			v[i] = i + 1
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Map16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r map[int]int
		v = make(map[int]int, math.MaxUint16+1)
		for i := 0; i < math.MaxUint16+1; i++ {
			v[i] = i + 1
		}
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Map32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v map[int]int
		var r map[uint]uint
		v = make(map[int]int, 100)
		for i := 0; i < 100; i++ {
			v[i] = i + 1
		}
		d, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Map16 {
			t.Error("code diffenrent")
		}
		err = msgpack.Unmarshal(d, &r)
		if err != nil {
			t.Error(err)
		}
		for k, vv := range v {
			if rv, ok := r[uint(k)]; !ok || rv != uint(vv) {
				t.Error("value diffrent")
			}
		}
		if len(v) != len(r) {
			t.Error("value different")
		}
	}
	// error
	{
		var v map[string]int
		var r map[int]int
		v = make(map[string]int, 100)
		for i := 0; i < 100; i++ {
			v[fmt.Sprintf("%03d", i)] = i
		}
		d, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Map16 {
			t.Error("code diffenrent")
		}
		err = msgpack.Unmarshal(d, &r)
		if err == nil || !strings.Contains(err.Error(), "invalid code a3 decoding") {
			t.Error("error")
		}
	}
	{
		var v map[int]string
		var r map[int]int
		v = make(map[int]string, 100)
		for i := 0; i < 100; i++ {
			v[i] = fmt.Sprint(i % 10)
		}
		d, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Map16 {
			t.Error("code diffenrent")
		}
		err = msgpack.Unmarshal(d, &r)
		if err == nil || !strings.Contains(err.Error(), "invalid code a1 decoding") {
			t.Error("error", err)
		}
	}
}

func TestPointer(t *testing.T) {
	{
		var v, r *int
		vv := 250
		v = &vv
		if err := encdec3(t, v, &r, func(code byte) bool {
			return code == def.Uint8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r *int
		d, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Nil {
			t.Error("code diffenrent")
		}
		err = msgpack.Unmarshal(d, &r)
		if err != nil {
			t.Error(err)
		}
		if v != r {
			t.Error("value different")
		}
	}
	// error
	{
		var v *int
		var r int
		if err := encdec3(t, v, r, func(code byte) bool {
			return code == def.Nil
		}); err == nil || !strings.Contains(err.Error(), "holder must set pointer value. but got:") {
			t.Error(err)
		}
	}
}

func TestUnsupported(t *testing.T) {
	b := []byte{0xc0}
	{
		var v, r uintptr
		_, err := msgpack.Marshal(v)
		if !strings.Contains(err.Error(), "type(uintptr) is unsupported") {
			t.Error("test error", err)
		}
		err = msgpack.Unmarshal(b, &r)
		if !strings.Contains(err.Error(), "type(uintptr) is unsupported") {
			t.Error("test error", err)
		}
	}
	{
		var v, r chan string
		_, err := msgpack.Marshal(v)
		if !strings.Contains(err.Error(), "type(chan) is unsupported") {
			t.Error("test error", err)
		}
		err = msgpack.Unmarshal(b, &r)
		if !strings.Contains(err.Error(), "type(chan) is unsupported") {
			t.Error("test error", err)
		}
	}
	{
		var v, r func()
		_, err := msgpack.Marshal(v)
		if !strings.Contains(err.Error(), "type(func) is unsupported") {
			t.Error("test error", err)
		}
		err = msgpack.Unmarshal(b, &r)
		if !strings.Contains(err.Error(), "type(func) is unsupported") {
			t.Error("test error", err)
		}
	}
	{
		// error reflect kind is invalid. current version set nil (0xc0)
		var v, r error
		bb, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if bb[0] != def.Nil {
			t.Errorf("code is different %d, %d", bb[0], def.Nil)
		}
		err = msgpack.Unmarshal(b, &r)
		if err != nil {
			t.Error(err)
		}
		if r != nil {
			t.Error("error should be nil")
		}
	}
}

/////////////////////////////////////////////////////////////////

func TestStruct(t *testing.T) {
	testSturctCode(t)
	testStructTag(t)
	testStructArray(t)
	testEmbedded(t)
	testStructJump(t)

	testStructUseCase(t)
	msgpack.StructAsArray = true
	testStructUseCase(t)
}

func testEmbedded(t *testing.T) {
	type Emb struct {
		Int int
	}
	type A struct {
		Emb
	}
	v := A{Emb: Emb{Int: 2}}
	b, err := msgpack.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	var vv A
	err = msgpack.Unmarshal(b, &vv)
	if err != nil {
		t.Error(err)
	}
	if v.Int != vv.Int {
		t.Errorf("value is different %v, %v", v, vv)
	}
}

func testStructTag(t *testing.T) {
	type vSt struct {
		One int    `msgpack:"Three"`
		Two string `msgpack:"four"`
		Hfn bool   `msgpack:"-"`
	}
	type rSt struct {
		Three int
		Four  string `msgpack:"four"`
		Hfn   bool
	}

	msgpack.StructAsArray = false

	v := vSt{One: 1, Two: "2", Hfn: true}
	r := rSt{}

	d, err := msgpack.MarshalAsMap(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.FixMap+0x02 {
		t.Error("code different")
	}
	err = msgpack.UnmarshalAsMap(d, &r)
	if err != nil {
		t.Error(err)
	}
	if v.One != r.Three || v.Two != r.Four || r.Hfn != false {
		t.Error("error:", v, r)
	}
}

func testStructArray(t *testing.T) {
	type vSt struct {
		One  int
		Two  string
		Ten  float32
		Skip float32
	}
	type rSt struct {
		Three int
		Four  string
		Tem   float32
	}

	msgpack.StructAsArray = true

	v := vSt{One: 1, Two: "2", Ten: 1.234}
	r := rSt{}

	d, err := msgpack.MarshalAsArray(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.FixArray+0x04 {
		t.Error("code different")
	}
	err = msgpack.UnmarshalAsArray(d, &r)
	if err != nil {
		t.Error(err)
	}
	if v.One != r.Three || v.Two != r.Four || v.Ten != r.Tem {
		t.Error("error:", v, r)
	}
}

func testSturctCode(t *testing.T) {
	type st1 struct {
		Int int
	}
	type st16 struct {
		I1  int
		I2  int
		I3  int
		I4  int
		I5  int
		I6  int
		I7  int
		I8  int
		I9  int
		I10 int
		I11 int
		I12 int
		I13 int
		I14 int
		I15 int
		I16 int
	}
	v1 := st1{Int: math.MinInt32}
	v16 := st16{I1: 1, I2: 2, I3: 3, I4: 4, I5: 5, I6: 6, I7: 7, I8: 8, I9: 9, I10: 10,
		I11: 11, I12: 12, I13: 13, I14: 14, I15: 15, I16: 16}

	{
		var r1 st1
		var r16 st16
		if err := encdec3(t, v1, &r1, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
		if err := encdec3(t, v16, &r16, func(code byte) bool {
			return code == def.Map16
		}); err != nil {
			t.Error(err)
		}
	}
	msgpack.StructAsArray = true
	{
		var r1 st1
		var r16 st16
		if err := encdec3(t, v1, &r1, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
		if err := encdec3(t, v16, &r16, func(code byte) bool {
			return code == def.Array16
		}); err != nil {
			t.Error(err)
		}
	}
	msgpack.StructAsArray = false
}

func testStructUseCase(t *testing.T) {
	type child3 struct {
		Int int
	}
	type child2 struct {
		Int2Uint        map[int]uint
		Float2Bool      map[float32]bool
		Duration2Struct map[time.Duration]child3
	}
	type child struct {
		IntArray      []int
		UintArray     []uint
		FloatArray    []float32
		BoolArray     []bool
		StringArray   []string
		TimeArray     []time.Time
		DurationArray []time.Duration
		Child         child2
	}
	type st struct {
		Int8        int8
		Int16       int16
		Int32       int32
		Int64       int64
		Uint8       byte
		Uint16      uint16
		Uint32      uint32
		Uint64      uint64
		Float       float32
		Double      float64
		Bool        bool
		String      string
		Pointer     *int
		Nil         *int
		Time        time.Time
		Duration    time.Duration
		Child       child
		Child3Array []child3
	}
	p := rand.Int()
	v := &st{
		Int32:    -32,
		Int8:     -8,
		Int16:    -16,
		Int64:    -64,
		Uint32:   32,
		Uint8:    8,
		Uint16:   16,
		Uint64:   64,
		Float:    1.23,
		Double:   2.3456,
		Bool:     true,
		String:   "Parent",
		Pointer:  &p,
		Nil:      nil,
		Time:     now,
		Duration: time.Duration(123 * time.Second),

		// child
		Child: child{
			IntArray:      []int{-1, -2, -3, -4, -5},
			UintArray:     []uint{1, 2, 3, 4, 5},
			FloatArray:    []float32{-1.2, -3.4, -5.6, -7.8},
			BoolArray:     []bool{true, true, false, false, true},
			StringArray:   []string{"str", "ing", "arr", "ay"},
			TimeArray:     []time.Time{now, now, now},
			DurationArray: []time.Duration{time.Duration(1 * time.Nanosecond), time.Duration(2 * time.Nanosecond)},

			// childchild
			Child: child2{
				Int2Uint:        map[int]uint{-1: 2, -3: 4},
				Float2Bool:      map[float32]bool{-1.1: true, -2.2: false},
				Duration2Struct: map[time.Duration]child3{time.Duration(1 * time.Hour): child3{Int: 1}, time.Duration(2 * time.Hour): child3{Int: 2}},
			},
		},

		Child3Array: []child3{child3{Int: 100}, child3{Int: 1000000}, child3{Int: 100000000}},
	}

	r1, r2 := st{}, st{}
	d1, d2, err := encSt(t, v, false)
	if err != nil {
		t.Error(err)
	}
	err = decSt(t, d1, d2, &r1, &r2, false)
	if err != nil {
		t.Error(err)
	}
	if err := equalCheck(v, r1); err != nil {
		t.Error(err)
	}
	if err := equalCheck(v, r2); err != nil {
		t.Error(err)
	}
}

func testStructJump(t *testing.T) {
	type v1 struct{ A interface{} }
	type r1 struct{ B interface{} }

	f := func(v v1) error {
		b, err := msgpack.MarshalAsMap(v)
		if err != nil {
			return err
		}
		var r r1
		err = msgpack.UnmarshalAsMap(b, &r)
		if err != nil {
			return err
		}
		if fmt.Sprint(v.A) == fmt.Sprint(r.B) {
			return fmt.Errorf("value equal %v, %v", v, r)
		}
		return nil
	}

	a1 := make([]int, math.MaxUint16)
	a2 := make([]int, math.MaxUint16+1)
	m1 := map[string]int{}
	m2 := map[string]int{}

	for i := range a1 {
		a1[i] = i
		m1[fmt.Sprint(i)] = 1
	}
	for i := range a2 {
		a2[i] = i
		m2[fmt.Sprint(i)] = 1
	}

	vs := []v1{
		{A: true},
		{A: 1}, {A: -1},
		{A: math.MaxUint8}, {A: math.MinInt8},
		{A: math.MaxUint16}, {A: math.MinInt16},
		{A: math.MaxUint32 + 1}, {A: math.MinInt32 - 1}, {A: math.MaxFloat64},
		{A: "a"},
		{A: strings.Repeat("b", math.MaxUint8)}, {A: []byte(strings.Repeat("c", math.MaxUint8))},
		{A: strings.Repeat("e", math.MaxUint16)}, {A: []byte(strings.Repeat("d", math.MaxUint16))},
		{A: strings.Repeat("f", math.MaxUint16+1)}, {A: []byte(strings.Repeat("g", math.MaxUint16+1))},
		{A: []int{1}}, {A: a1}, {A: a2},
		{A: map[string]int{"a": 1}}, {A: m1}, {A: m2},
		{A: time.Unix(now.Unix(), int64(now.Nanosecond()))},
	}

	for i, v := range vs {
		if err := f(v); err != nil {
			t.Error(i, err)
		}
	}

}

func encSt(t *testing.T, in interface{}, isDebug bool) ([]byte, []byte, error) {
	var d1, d2 []byte
	var err error
	d1, err = msgpack.Marshal(in)
	if err != nil {
		return nil, nil, err
	}

	if msgpack.StructAsArray {
		d2, err = msgpack.MarshalAsMap(in)
	} else {
		d2, err = msgpack.MarshalAsArray(in)
	}
	if err != nil {
		return nil, nil, err
	}

	if isDebug {
		t.Log(in, " -- to byte --> ", d1)
		t.Log(in, " -- to byte --> ", d2)
	}
	return d1, d2, nil
}

func decSt(t *testing.T, d1, d2 []byte, out1, out2 interface{}, isDebug bool) error {
	if err := msgpack.Unmarshal(d1, out1); err != nil {
		return err
	}

	if msgpack.StructAsArray {
		err := msgpack.UnmarshalAsMap(d2, out2)
		if err != nil {
			return err
		}
	} else {
		err := msgpack.UnmarshalAsArray(d2, out2)
		if err != nil {
			return err
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////

func TestExt(t *testing.T) {
	err := msgpack.AddExtCoder(encoder, decoder)
	if err != nil {
		t.Error(err)
	}

	{
		v := ExtInt{
			Int8:        math.MinInt8,
			Int16:       math.MinInt16,
			Int32:       math.MinInt32,
			Int64:       math.MinInt32 - 1,
			Uint8:       math.MaxUint8,
			Uint16:      math.MaxUint16,
			Uint32:      math.MaxUint32,
			Uint64:      math.MaxUint32 + 1,
			Byte2Int:    rand.Intn(math.MaxUint16),
			Byte4Int:    rand.Intn(math.MaxInt32) - rand.Intn(math.MaxInt32-1),
			Byte4Uint32: math.MaxUint32 - 1,
			Bytes:       []byte{1, 2, 3}, // 3
		}
		var r1, r2 ExtInt
		d1, d2, err := encSt(t, v, false)
		if err != nil {
			t.Error(err)
		}
		err = decSt(t, d1, d2, &r1, &r2, false)
		if err != nil {
			t.Error(err)
		}
		if err := equalCheck(v, r1); err != nil {
			t.Error(err)
		}
		if err := equalCheck(r1, r2); err != nil {
			t.Error(err)
		}
	}
	err = msgpack.RemoveExtCoder(encoder, decoder)
	if err != nil {
		t.Error(err)
	}
	{
		v := ExtInt{
			Int8:        math.MinInt8,
			Int16:       math.MinInt16,
			Int32:       math.MinInt32,
			Int64:       math.MinInt32 - 1,
			Uint8:       math.MaxUint8,
			Uint16:      math.MaxUint16,
			Uint32:      math.MaxUint32,
			Uint64:      math.MaxUint32 + 1,
			Byte2Int:    rand.Intn(math.MaxUint16),
			Byte4Int:    rand.Intn(math.MaxInt32) - rand.Intn(math.MaxInt32-1),
			Byte4Uint32: math.MaxUint32 - 1,
			Bytes:       []byte{4, 5, 6}, // 3
		}
		var r1, r2 ExtInt
		d1, d2, err := encSt(t, v, false)
		if err != nil {
			t.Error(err)
		}
		err = decSt(t, d1, d2, &r1, &r2, false)
		if err != nil {
			t.Error(err)
		}
		if err := equalCheck(v, r1); err != nil {
			t.Error(err)
		}
		if err := equalCheck(r1, r2); err != nil {
			t.Error(err)
		}
	}

	// error
	enc2, dec2 := new(testExt2Encoder), new(testExt2Decoder)
	err = msgpack.AddExtCoder(enc2, dec2)
	if err != nil && strings.Contains(err.Error(), "code different") {
		// ok
	} else {
		t.Error("unreachable", err)
	}
	err = msgpack.RemoveExtCoder(enc2, dec2)
	if err != nil && strings.Contains(err.Error(), "code different") {
		// ok
	} else {
		t.Error("unreachable", err)
	}
}

type ExtStruct struct {
	Int8        int8
	Int16       int16
	Int32       int32
	Int64       int64
	Uint8       uint8
	Uint16      uint16
	Uint32      uint32
	Uint64      uint64
	Byte2Int    int
	Byte4Int    int
	Byte4Uint32 uint32
	Bytes       []byte // 3
}
type ExtInt ExtStruct

var decoder = new(testDecoder)

type testDecoder struct {
	ext.DecoderCommon
}

var extIntCode = int8(-2)

func (td *testDecoder) Code() int8 {
	return extIntCode
}

func (td *testDecoder) IsType(offset int, d *[]byte) bool {
	code, offset := td.ReadSize1(offset, d)
	if code == def.Ext8 {
		c, offset := td.ReadSize1(offset, d)
		t, _ := td.ReadSize1(offset, d)
		return c == 15+15+10+3 && int8(t) == td.Code()
	}
	return false
}

func (td *testDecoder) AsValue(offset int, k reflect.Kind, d *[]byte) (interface{}, int, error) {
	code, offset := td.ReadSize1(offset, d)

	switch code {
	case def.Ext8:
		// size
		_, offset = td.ReadSize1(offset, d)
		// code
		_, offset = td.ReadSize1(offset, d)
		i8, offset := td.ReadSize1(offset, d)
		i16, offset := td.ReadSize2(offset, d)
		i32, offset := td.ReadSize4(offset, d)
		i64, offset := td.ReadSize8(offset, d)
		u8, offset := td.ReadSize1(offset, d)
		u16, offset := td.ReadSize2(offset, d)
		u32, offset := td.ReadSize4(offset, d)
		u64, offset := td.ReadSize8(offset, d)
		b16, offset := td.ReadSize2(offset, d)
		b32, offset := td.ReadSize4(offset, d)
		bu32, offset := td.ReadSize4(offset, d)
		bs, offset := td.ReadSizeN(offset, 3, d)
		return ExtInt{
			Int8:        int8(i8),
			Int16:       int16(binary.BigEndian.Uint16(i16)),
			Int32:       int32(binary.BigEndian.Uint32(i32)),
			Int64:       int64(binary.BigEndian.Uint64(i64)),
			Uint8:       u8,
			Uint16:      binary.BigEndian.Uint16(u16),
			Uint32:      binary.BigEndian.Uint32(u32),
			Uint64:      binary.BigEndian.Uint64(u64),
			Byte2Int:    int(binary.BigEndian.Uint16(b16)),
			Byte4Int:    int(int32(binary.BigEndian.Uint32(b32))),
			Byte4Uint32: binary.BigEndian.Uint32(bu32),
			Bytes:       bs,
		}, offset, nil
	}

	return ExtInt{}, 0, fmt.Errorf("should not reach this line!! code %x decoding %v", code, k)
}

var encoder = new(testEncoder)

type testEncoder struct {
	ext.EncoderCommon
}

func (s *testEncoder) Code() int8 {
	return extIntCode
}

func (s *testEncoder) Type() reflect.Type {
	return reflect.TypeOf(ExtInt{})
}

func (s *testEncoder) CalcByteSize(value reflect.Value) (int, error) {
	t := value.Interface().(ExtInt)
	return def.Byte1 + def.Byte1 + 15 + 15 + 10 + len(t.Bytes), nil
}

func (s *testEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	t := value.Interface().(ExtInt)
	offset = s.SetByte1Int(def.Ext8, offset, bytes)
	offset = s.SetByte1Int(15+15+10+len(t.Bytes), offset, bytes)
	offset = s.SetByte1Int(int(s.Code()), offset, bytes)

	offset = s.SetByte1Int64(int64(t.Int8), offset, bytes)
	offset = s.SetByte2Int64(int64(t.Int16), offset, bytes)
	offset = s.SetByte4Int64(int64(t.Int32), offset, bytes)
	offset = s.SetByte8Int64(t.Int64, offset, bytes)

	offset = s.SetByte1Uint64(uint64(t.Uint8), offset, bytes)
	offset = s.SetByte2Uint64(uint64(t.Uint16), offset, bytes)
	offset = s.SetByte4Uint64(uint64(t.Uint32), offset, bytes)
	offset = s.SetByte8Uint64(t.Uint64, offset, bytes)

	offset = s.SetByte2Int(t.Byte2Int, offset, bytes)
	offset = s.SetByte4Int(t.Byte4Int, offset, bytes)

	offset = s.SetByte4Uint32(t.Byte4Uint32, offset, bytes)
	offset = s.SetBytes(t.Bytes, offset, bytes)
	return offset
}

/////////////////////////////////////////////////////////

type Ext2Struct struct {
	V int
}
type Ext2Int Ext2Struct

type testExt2Decoder struct {
	ext.DecoderCommon
}

func (td *testExt2Decoder) Code() int8 {
	return 3
}

func (td *testExt2Decoder) IsType(offset int, d *[]byte) bool {
	return false
}

func (td *testExt2Decoder) AsValue(offset int, k reflect.Kind, d *[]byte) (interface{}, int, error) {
	return Ext2Int{}, 0, fmt.Errorf("should not reach this line!! code %x decoding %v", 3, k)
}

type testExt2Encoder struct {
	ext.EncoderCommon
}

func (s *testExt2Encoder) Code() int8 {
	return -3
}

func (s *testExt2Encoder) Type() reflect.Type {
	return reflect.TypeOf(ExtInt{})
}

func (s *testExt2Encoder) CalcByteSize(value reflect.Value) (int, error) {
	return 0, nil
}

func (s *testExt2Encoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	return offset
}

/////////////////////////////////////////////////////////

type (
	checker      func(data []byte) bool
	marshaller   func(v any) ([]byte, error)
	unmarshaller func(data []byte, v any) error
)

var marshallers = []struct {
	name string
	m    marshaller
}{
	{"Marshal", msgpack.Marshal},
	{"MarshalWrite", func(v any) ([]byte, error) {
		buf := bytes.Buffer{}
		err := msgpack.MarshalWrite(&buf, v)
		return buf.Bytes(), err
	}},
}
var unmarshallers = []struct {
	name string
	u    unmarshaller
}{
	{"Unmarshal", msgpack.Unmarshal},
	{"UnmarshalRead", func(data []byte, v any) error {
		return msgpack.UnmarshalRead(bytes.NewReader(data), v)
	}},
}

type encdecArg[T any] struct {
	n      string
	v      any
	r      T
	c      checker
	skipEq bool
	vc     func(t T) error
	e      string
}

func encdec[T any](t *testing.T, args ...encdecArg[T]) {
	t.Helper()

	for _, arg := range args {
		t.Run(arg.n, func(t *testing.T) {
			for _, m := range marshallers {
				for _, u := range unmarshallers {
					t.Run(m.name+"-"+u.name, func(t *testing.T) {
						var e error
						defer func() {
							if e != nil {
								if len(arg.e) < 1 {
									t.Fatalf("unexpected error: %v", e)
								}
								if !strings.Contains(e.Error(), arg.e) {
									t.Fatalf("error does not contain '%s'. err: %v", arg.e, e)
								}
							} else {
								if len(arg.e) > 0 {
									t.Fatalf("error should occur, but nil. expected: %s", arg.e)
								}
							}
						}()

						d, err := m.m(arg.v)
						if err != nil {
							e = err
							return
						}
						if arg.c != nil && !arg.c(d) {
							e = fmt.Errorf("different %s", hex.Dump(d))
							return
						}
						if err = u.u(d, &arg.r); err != nil {
							e = err
							return
						}
						if !arg.skipEq {
							if err = equalCheck(arg.v, arg.r); err != nil {
								e = err
								return
							}
						}
						if arg.vc != nil {
							if err = arg.vc(arg.r); err != nil {
								e = err
								return
							}
						}
					})
				}
			}
		})
	}
}

func encdec3(t *testing.T, v, r interface{}, j func(byte) bool, errStr ...string) error {
	t.Helper()

	marshallers := []struct {
		name string
		m    marshaller
	}{
		{"Marshal", msgpack.Marshal},
		{"MarshalWrite", func(v any) ([]byte, error) {
			buf := bytes.Buffer{}
			err := msgpack.MarshalWrite(&buf, v)
			return buf.Bytes(), err
		}},
	}
	unmarshallers := []struct {
		name string
		u    unmarshaller
	}{
		{"Unmarshal", msgpack.Unmarshal},
		{"UnmarshalRead", func(data []byte, v any) error {
			return msgpack.UnmarshalRead(bytes.NewReader(data), r)
		}},
	}

	for _, m := range marshallers {
		for _, u := range unmarshallers {
			t.Run(m.name+"-"+u.name, func(t *testing.T) {
				var e error
				defer func() {
					if e != nil {
						if len(errStr) < 1 {
							t.Fatalf("unexpected error: %v", e)
						}
						if !strings.Contains(e.Error(), errStr[0]) {
							t.Fatalf("error does not contain '%s'. err: %v", errStr[0], e)
						}
					} else {
						if len(errStr) > 0 {
							t.Fatalf("error should occur, but nil")
						}
					}
				}()

				d, err := m.m(v)
				if err != nil {
					e = err
					return
				}
				if j != nil && !j(d[0]) {
					e = fmt.Errorf("different %s", hex.Dump(d))
					return
				}
				if err = u.u(d, r); err != nil {
					e = err
					return
				}
				if err = equalCheck(v, r); err != nil {
					e = err
					return
				}
			})
		}
	}

	return nil
}

// for check value
func getValue(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Invalid {
		return nil
	}
	return rv.Interface()
}

func equalCheck(in, out interface{}) error {
	i := getValue(in)
	o := getValue(out)
	if !reflect.DeepEqual(i, o) {
		return errors.New(fmt.Sprint("value different \n[in]:", i, " \n[out]:", o))
	}
	return nil
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func ErrorContains(t *testing.T, err error, errStr string) {
	if err == nil {
		t.Fatal("error should occur")
	}
	if !strings.Contains(err.Error(), errStr) {
		t.Fatalf("error does not contain '%s'. err: %v", errStr, err)
	}
}
