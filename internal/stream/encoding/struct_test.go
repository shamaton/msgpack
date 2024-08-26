package encoding

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_writeStruct(t *testing.T) {
	method := func(e *encoder) func(reflect.Value) error {
		return e.writeStruct
	}

	t.Run("Ext", func(t *testing.T) {
		value := time.Time{}
		testcases := AsXXXTestCases[time.Time]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Ext8,
					12, 0xff,
					0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xf1, 0x88, 0x6e, 0x09, 0x00,
				},
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("Array", func(t *testing.T) {
		st, _, b := tu.CreateStruct(1)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				AsArray:         true,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.FixArray + 1}, b...),
				AsArray:         true,
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("Map", func(t *testing.T) {
		st, b, _ := tu.CreateStruct(1)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				AsArray:         false,
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.FixMap + 1}, b...),
				AsArray:         false,
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
}

func Test_writeStructArray(t *testing.T) {
	method := func(e *encoder) func(reflect.Value) error {
		return e.writeStructArray
	}

	t.Run("FixArray", func(t *testing.T) {
		st, _, b := tu.CreateStruct(0x0f)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error.def",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "error.value",
				Value:           value,
				BufferSize:      2,
				PreWriteSize:    1,
				Contains:        []byte{def.FixArray + byte(0x0f)},
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.FixArray + byte(0x0f)}, b...),
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("Array16", func(t *testing.T) {
		st, _, b := tu.CreateStruct(math.MaxUint16)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error.def",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "error.num",
				Value:           value,
				BufferSize:      2,
				PreWriteSize:    1,
				Contains:        []byte{def.Array16},
				MethodForStruct: method,
			},
			{
				Name:            "error.value",
				Value:           value,
				BufferSize:      4,
				PreWriteSize:    1,
				Contains:        []byte{def.Array16, 0xff, 0xff},
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.Array16, 0xff, 0xff}, b...),
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("Array32", func(t *testing.T) {
		st, _, b := tu.CreateStruct(math.MaxUint16 + 1)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error.def",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "error.num",
				Value:           value,
				BufferSize:      2,
				PreWriteSize:    1,
				Contains:        []byte{def.Array32},
				MethodForStruct: method,
			},
			{
				Name:            "error.value",
				Value:           value,
				BufferSize:      6,
				PreWriteSize:    1,
				Contains:        []byte{def.Array32, 0x00, 0x01, 0x00, 0x00},
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.Array32, 0x00, 0x01, 0x00, 0x00}, b...),
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
}

func Test_writeStructMap(t *testing.T) {
	method := func(e *encoder) func(reflect.Value) error {
		return e.writeStructMap
	}

	t.Run("FixMap", func(t *testing.T) {
		st, b, _ := tu.CreateStruct(0x0f)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error.def",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "error.key",
				Value:           value,
				BufferSize:      2,
				PreWriteSize:    1,
				Contains:        []byte{def.FixMap + byte(0x0f)},
				MethodForStruct: method,
			},
			{
				Name:            "error.value",
				Value:           value,
				BufferSize:      5,
				PreWriteSize:    1,
				Contains:        append([]byte{def.FixMap + byte(0x0f)}, b[:3]...),
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.FixMap + byte(0x0f)}, b...),
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("Map16", func(t *testing.T) {
		st, b, _ := tu.CreateStruct(math.MaxUint16)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error.def",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "error.num",
				Value:           value,
				BufferSize:      2,
				PreWriteSize:    1,
				Contains:        []byte{def.Map16},
				MethodForStruct: method,
			},
			{
				Name:            "error.key",
				Value:           value,
				BufferSize:      4,
				PreWriteSize:    1,
				Contains:        []byte{def.Map16, 0xff, 0xff},
				MethodForStruct: method,
			},
			{
				Name:            "error.value",
				Value:           value,
				BufferSize:      7,
				PreWriteSize:    1,
				Contains:        append([]byte{def.Map16, 0xff, 0xff}, b[:3]...),
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.Map16, 0xff, 0xff}, b...),
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("Map32", func(t *testing.T) {
		st, b, _ := tu.CreateStruct(math.MaxUint16 + 1)
		value := reflect.ValueOf(st).Elem().Interface()
		testcases := AsXXXTestCases[any]{
			{
				Name:            "error.def",
				Value:           value,
				BufferSize:      1,
				PreWriteSize:    1,
				MethodForStruct: method,
			},
			{
				Name:            "error.num",
				Value:           value,
				BufferSize:      2,
				PreWriteSize:    1,
				Contains:        []byte{def.Map32},
				MethodForStruct: method,
			},
			{
				Name:            "error.key",
				Value:           value,
				BufferSize:      6,
				PreWriteSize:    1,
				Contains:        []byte{def.Map32, 0x00, 0x01, 0x00, 0x00},
				MethodForStruct: method,
			},
			{
				Name:            "error.value",
				Value:           value,
				BufferSize:      9,
				PreWriteSize:    1,
				Contains:        append([]byte{def.Map32, 0x00, 0x01, 0x00, 0x00}, b[:3]...),
				MethodForStruct: method,
			},
			{
				Name:            "ok",
				Value:           value,
				Expected:        append([]byte{def.Map32, 0x00, 0x01, 0x00, 0x00}, b...),
				BufferSize:      1,
				MethodForStruct: method,
			},
		}
		testcases.Run(t)
	})
}
