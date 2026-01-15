package encoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func Test_calcStructArray(t *testing.T) {
	type b struct {
		B []byte
	}

	t.Run("non-cache", func(t *testing.T) {
		value := b{B: make([]byte, math.MaxUint32+1)}
		e := encoder{}
		rv := reflect.ValueOf(value)
		_, err := e.calcStructArray(rv)
		tu.IsError(t, err, def.ErrUnsupportedType)
	})
	t.Run("cache", func(t *testing.T) {
		value := b{B: make([]byte, 1)}
		e := encoder{}
		rv := reflect.ValueOf(value)
		_, err := e.calcStructArray(rv)
		tu.NoError(t, err)

		value = b{B: make([]byte, math.MaxUint32+1)}
		rv = reflect.ValueOf(value)
		_, err = e.calcStructArray(rv)
		tu.IsError(t, err, def.ErrUnsupportedType)
	})

	testcases := []struct {
		name   string
		value  int
		result int
		error  error
	}{
		{
			name:   "0x0f",
			value:  0x0f,
			result: def.Byte1,
		},
		{
			name:   "u16",
			value:  math.MaxUint16,
			result: def.Byte1 + def.Byte2,
		},
		{
			name:   "u32",
			value:  math.MaxUint16 + 1,
			result: def.Byte1 + def.Byte4,
		},
		// can not run by out of memory
		//{
		//	name:   "u32over",
		//	value:  math.MaxUint32 + 1,
		//	result: 0,
		//},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := encoder{}
			v, _, bs := tu.CreateStruct(tc.value)
			rv := reflect.ValueOf(v).Elem()
			result, err := e.calcStructArray(rv)
			tu.IsError(t, err, tc.error)
			tu.Equal(t, result, tc.result+len(bs))
		})
	}
}

func Test_calcStructMap(t *testing.T) {
	type b struct {
		B []byte
	}

	t.Run("non-cache", func(t *testing.T) {
		value := b{B: make([]byte, math.MaxUint32+1)}
		e := encoder{}
		rv := reflect.ValueOf(value)
		_, err := e.calcStructMap(rv)
		tu.IsError(t, err, def.ErrUnsupportedType)
	})
	t.Run("cache", func(t *testing.T) {
		value := b{B: make([]byte, 1)}
		e := encoder{}
		rv := reflect.ValueOf(value)
		_, err := e.calcStructMap(rv)
		tu.NoError(t, err)

		value = b{B: make([]byte, math.MaxUint32+1)}
		rv = reflect.ValueOf(value)
		_, err = e.calcStructMap(rv)
		tu.IsError(t, err, def.ErrUnsupportedType)
	})

	testcases := []struct {
		name   string
		value  int
		result int
		error  error
	}{
		{
			name:   "0x0f",
			value:  0x0f,
			result: def.Byte1,
		},
		{
			name:   "u16",
			value:  math.MaxUint16,
			result: def.Byte1 + def.Byte2,
		},
		{
			name:   "u32",
			value:  math.MaxUint16 + 1,
			result: def.Byte1 + def.Byte4,
		},
		// can not run by out of memory
		//{
		//	name:   "u32over",
		//	value:  math.MaxUint32 + 1,
		//	result: 0,
		//},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := encoder{}
			v, bs, _ := tu.CreateStruct(tc.value)
			rv := reflect.ValueOf(v).Elem()
			result, err := e.calcStructMap(rv)
			tu.IsError(t, err, tc.error)
			tu.Equal(t, result, tc.result+len(bs))
		})
	}
}

func Test_writeStructArray(t *testing.T) {
	testcases := []struct {
		name  string
		value int
		code  byte
	}{
		{
			name:  "0x0f",
			value: 0x0f,
			code:  def.FixArray + 0x0f,
		},
		{
			name:  "u16",
			value: math.MaxUint16,
			code:  def.Array16,
		},
		{
			name:  "u32",
			value: math.MaxUint16 + 1,
			code:  def.Array32,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := encoder{}
			v, _, _ := tu.CreateStruct(tc.value)
			rv := reflect.ValueOf(v).Elem()
			size, err := e.calcStructArray(rv)
			tu.NoError(t, err)

			e.d = make([]byte, size)
			result := e.writeStructArray(rv, 0)
			tu.Equal(t, len(e.d), result)
			tu.Equal(t, e.d[0], tc.code)
		})
	}
}

func Test_writeStructMap(t *testing.T) {
	testcases := []struct {
		name  string
		value int
		code  byte
	}{
		{
			name:  "0x0f",
			value: 0x0f,
			code:  def.FixMap + 0x0f,
		},
		{
			name:  "u16",
			value: math.MaxUint16,
			code:  def.Map16,
		},
		{
			name:  "u32",
			value: math.MaxUint16 + 1,
			code:  def.Map32,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := encoder{}
			v, _, _ := tu.CreateStruct(tc.value)
			rv := reflect.ValueOf(v).Elem()
			size, err := e.calcStructMap(rv)
			tu.NoError(t, err)

			e.d = make([]byte, size)
			result := e.writeStructMap(rv, 0)
			tu.Equal(t, len(e.d), result)
			tu.Equal(t, e.d[0], tc.code)
		})
	}
}

func Test_calcSizeWithOmitEmpty(t *testing.T) {
	e := encoder{}
	var v any
	v = func() {}
	_, err := e.calcSizeWithOmitEmpty(reflect.ValueOf(v), "a", false)
	tu.Error(t, err)

	v = 1
	_, err = e.calcSizeWithOmitEmpty(reflect.ValueOf(v), "a", false)
	tu.NoError(t, err)
}

func Test_structCache(t *testing.T) {
	type forMap struct {
		Int  int
		Uint uint `msgpack:",omitempty"`
		Str  string
	}
	type forMapNoOmit struct {
		Int  int
		Uint uint
		Str  string
	}
	type forArray struct {
		Int  int `msgpack:",omitempty"`
		Uint uint
		Str  string `msgpack:",omitempty"`
	}
	type forArrayNoOmit struct {
		Int  int
		Uint uint
		Str  string
	}

	e := encoder{}
	testcases := []struct {
		v         any
		omitCount int
		f         func(rv reflect.Value) (int, error)
	}{
		{
			v:         forMap{},
			omitCount: 1,
			f:         e.calcStructMap,
		},

		{
			v:         forMapNoOmit{},
			omitCount: 0,
			f:         e.calcStructMap,
		},
		{
			v:         forArray{},
			omitCount: 2,
			f:         e.calcStructArray,
		},

		{
			v:         forArrayNoOmit{},
			omitCount: 0,
			f:         e.calcStructArray,
		},
	}

	for _, c := range testcases {
		rv := reflect.ValueOf(c.v)
		t.Run(rv.String(), func(t *testing.T) {
			_, found := cachemap.Load(rv.Type())
			tu.Equal(t, found, false)

			_, err := c.f(rv)
			tu.NoError(t, err)

			cache, _ := cachemap.Load(rv.Type())
			ca, ok := cache.(*structCache)
			tu.Equal(t, ok, true)
			tu.Equal(t, len(ca.omits), rv.NumField())
			count := 0
			for _, b := range ca.omits {
				if b {
					count++
				}
			}
			tu.Equal(t, count, c.omitCount)
		})
	}
}
