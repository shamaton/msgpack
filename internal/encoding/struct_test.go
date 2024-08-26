package encoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
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
			result: 0,
		},
		{
			name:   "u16",
			value:  math.MaxUint16,
			result: def.Byte2,
		},
		{
			name:   "u32",
			value:  math.MaxUint16 + 1,
			result: def.Byte4,
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
			result: 0,
		},
		{
			name:   "u16",
			value:  math.MaxUint16,
			result: def.Byte2,
		},
		{
			name:   "u32",
			value:  math.MaxUint16 + 1,
			result: def.Byte4,
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

			e.d = make([]byte, size+def.Byte1)
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

			e.d = make([]byte, size+def.Byte1)
			result := e.writeStructMap(rv, 0)
			tu.Equal(t, len(e.d), result)
			tu.Equal(t, e.d[0], tc.code)
		})
	}
}
