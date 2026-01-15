package encoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_writeMapLength(t *testing.T) {
	method := func(e *encoder) func(int) error {
		return e.writeMapLength
	}
	testcases := AsXXXTestCases[int]{
		{
			Name:         "FixMap.error",
			Value:        5,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:       "FixMap.ok",
			Value:      5,
			Expected:   []byte{def.FixMap + 5},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Map16.error.def",
			Value:        32,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Map16.error.value",
			Value:        32,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Map16},
			Method:       method,
		},
		{
			Name:       "Map16.ok",
			Value:      32,
			Expected:   []byte{def.Map16, 0x00, 0x20},
			BufferSize: 1,
			Method:     method,
		},
		{
			Name:         "Map32.error.def",
			Value:        65536,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:         "Map32.error.value",
			Value:        65536,
			BufferSize:   3,
			PreWriteSize: 1,
			Contains:     []byte{def.Map32},
			Method:       method,
		},
		{
			Name:       "Map32.ok",
			Value:      65536,
			Expected:   []byte{def.Map32, 0x00, 0x01, 0x00, 0x00},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}

func Test_writeFixedMap(t *testing.T) {
	method := func(e *encoder) func(reflect.Value) (bool, error) {
		return e.writeFixedMap
	}

	t.Run("map[string]int", func(t *testing.T) {
		value := map[string]int{"a": -1}
		testcases := AsXXXTestCases[map[string]int]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.FixStr + 1, 0x61, 0xff},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]uint", func(t *testing.T) {
		value := map[string]uint{"a": 1}
		testcases := AsXXXTestCases[map[string]uint]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.FixStr + 1, 0x61, 0x01},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]float32", func(t *testing.T) {
		value := map[string]float32{"a": 1}
		testcases := AsXXXTestCases[map[string]float32]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.FixStr + 1, 0x61, def.Float32, 0x3f, 0x80, 0x00, 0x00},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]float64", func(t *testing.T) {
		value := map[string]float64{"a": 1}
		testcases := AsXXXTestCases[map[string]float64]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Float64, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]bool", func(t *testing.T) {
		value := map[string]bool{"a": true}
		testcases := AsXXXTestCases[map[string]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.FixStr + 1, 0x61, def.True},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]string", func(t *testing.T) {
		value := map[string]string{"a": "b"}
		testcases := AsXXXTestCases[map[string]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.FixStr + 1, 0x61, def.FixStr + 1, 0x62},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]int8", func(t *testing.T) {
		value := map[string]int8{"a": math.MinInt8}
		testcases := AsXXXTestCases[map[string]int8]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Int8, 0x80,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]int16", func(t *testing.T) {
		value := map[string]int16{"a": math.MinInt16}
		testcases := AsXXXTestCases[map[string]int16]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Int16, 0x80, 0x00,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]int32", func(t *testing.T) {
		value := map[string]int32{"a": math.MinInt32}
		testcases := AsXXXTestCases[map[string]int32]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Int32, 0x80, 0x00, 0x00, 0x00,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]int64", func(t *testing.T) {
		value := map[string]int64{"a": math.MinInt64}
		testcases := AsXXXTestCases[map[string]int64]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]uint8", func(t *testing.T) {
		value := map[string]uint8{"a": math.MaxUint8}
		testcases := AsXXXTestCases[map[string]uint8]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Uint8, 0xff,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]uint16", func(t *testing.T) {
		value := map[string]uint16{"a": math.MaxUint16}
		testcases := AsXXXTestCases[map[string]uint16]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Uint16, 0xff, 0xff,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]uint32", func(t *testing.T) {
		value := map[string]uint32{"a": math.MaxUint32}
		testcases := AsXXXTestCases[map[string]uint32]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Uint32, 0xff, 0xff, 0xff, 0xff,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[string]uint64", func(t *testing.T) {
		value := map[string]uint64{"a": math.MaxUint64}
		testcases := AsXXXTestCases[map[string]uint64]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.FixStr + 1, 0x61},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.FixStr + 1, 0x61,
					def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int]string", func(t *testing.T) {
		value := map[int]string{math.MinInt8: "a"}
		testcases := AsXXXTestCases[map[int]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Int8, 0x80},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.Int8, 0x80, def.FixStr + 1, 0x61},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int]bool", func(t *testing.T) {
		value := map[int]bool{math.MinInt8: true}
		testcases := AsXXXTestCases[map[int]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Int8, 0x80},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.Int8, 0x80, def.True},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint]string", func(t *testing.T) {
		value := map[uint]string{math.MaxUint8: "a"}
		testcases := AsXXXTestCases[map[uint]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint8, 0xff},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.Uint8, 0xff, def.FixStr + 1, 0x61},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint]bool", func(t *testing.T) {
		value := map[uint]bool{math.MaxUint8: true}
		testcases := AsXXXTestCases[map[uint]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint8, 0xff},
				MethodForFixed: method,
			},
			{
				Name:           "ok",
				Value:          value,
				Expected:       []byte{def.Uint8, 0xff, def.True},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[float32]string", func(t *testing.T) {
		value := map[float32]string{1: "a"}
		testcases := AsXXXTestCases[map[float32]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     6,
				PreWriteSize:   1,
				Contains:       []byte{def.Float32, 0x3f, 0x80, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Float32, 0x3f, 0x80, 0x00, 0x00,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[float32]bool", func(t *testing.T) {
		value := map[float32]bool{1: true}
		testcases := AsXXXTestCases[map[float32]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     6,
				PreWriteSize:   1,
				Contains:       []byte{def.Float32, 0x3f, 0x80, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Float32, 0x3f, 0x80, 0x00, 0x00,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[float64]string", func(t *testing.T) {
		value := map[float64]string{1: "a"}
		testcases := AsXXXTestCases[map[float64]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     10,
				PreWriteSize:   1,
				Contains:       []byte{def.Float64, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Float64, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[float64]bool", func(t *testing.T) {
		value := map[float64]bool{1: true}
		testcases := AsXXXTestCases[map[float64]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     10,
				PreWriteSize:   1,
				Contains:       []byte{def.Float64, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Float64, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int8]string", func(t *testing.T) {
		value := map[int8]string{math.MinInt8: "a"}
		testcases := AsXXXTestCases[map[int8]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Int8, 0x80},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int8, 0x80,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int8]bool", func(t *testing.T) {
		value := map[int8]bool{math.MinInt8: true}
		testcases := AsXXXTestCases[map[int8]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Int8, 0x80},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int8, 0x80,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int16]string", func(t *testing.T) {
		value := map[int16]string{math.MinInt16: "a"}
		testcases := AsXXXTestCases[map[int16]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     4,
				PreWriteSize:   1,
				Contains:       []byte{def.Int16, 0x80, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int16, 0x80, 0x00,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int16]bool", func(t *testing.T) {
		value := map[int16]bool{math.MinInt16: true}
		testcases := AsXXXTestCases[map[int16]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     4,
				PreWriteSize:   1,
				Contains:       []byte{def.Int16, 0x80, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int16, 0x80, 0x00,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int32]string", func(t *testing.T) {
		value := map[int32]string{math.MinInt32: "a"}
		testcases := AsXXXTestCases[map[int32]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     6,
				PreWriteSize:   1,
				Contains:       []byte{def.Int32, 0x80, 0x00, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int32, 0x80, 0x00, 0x00, 0x00,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int32]bool", func(t *testing.T) {
		value := map[int32]bool{math.MinInt32: true}
		testcases := AsXXXTestCases[map[int32]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     6,
				PreWriteSize:   1,
				Contains:       []byte{def.Int32, 0x80, 0x00, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int32, 0x80, 0x00, 0x00, 0x00,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int64]string", func(t *testing.T) {
		value := map[int64]string{math.MinInt64: "a"}
		testcases := AsXXXTestCases[map[int64]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     10,
				PreWriteSize:   1,
				Contains:       []byte{def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[int64]bool", func(t *testing.T) {
		value := map[int64]bool{math.MinInt64: true}
		testcases := AsXXXTestCases[map[int64]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     10,
				PreWriteSize:   1,
				Contains:       []byte{def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Int64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})

	t.Run("map[uint8]string", func(t *testing.T) {
		value := map[uint8]string{math.MaxUint8: "a"}
		testcases := AsXXXTestCases[map[uint8]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint8, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint8, 0xff,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint8]bool", func(t *testing.T) {
		value := map[uint8]bool{math.MaxUint8: true}
		testcases := AsXXXTestCases[map[uint8]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     3,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint8, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint8, 0xff,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint16]string", func(t *testing.T) {
		value := map[uint16]string{math.MaxUint16: "a"}
		testcases := AsXXXTestCases[map[uint16]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     4,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint16, 0xff, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint16, 0xff, 0xff,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint16]bool", func(t *testing.T) {
		value := map[uint16]bool{math.MaxUint16: true}
		testcases := AsXXXTestCases[map[uint16]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     4,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint16, 0xff, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint16, 0xff, 0xff,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint32]string", func(t *testing.T) {
		value := map[uint32]string{math.MaxUint32: "a"}
		testcases := AsXXXTestCases[map[uint32]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     6,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint32, 0xff, 0xff, 0xff, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint32, 0xff, 0xff, 0xff, 0xff,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint32]bool", func(t *testing.T) {
		value := map[uint32]bool{math.MaxUint32: true}
		testcases := AsXXXTestCases[map[uint32]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     6,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint32, 0xff, 0xff, 0xff, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint32, 0xff, 0xff, 0xff, 0xff,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint64]string", func(t *testing.T) {
		value := map[uint64]string{math.MaxUint64: "a"}
		testcases := AsXXXTestCases[map[uint64]string]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     10,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					def.FixStr + 1, 0x61,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
	t.Run("map[uint64]bool", func(t *testing.T) {
		value := map[uint64]bool{math.MaxUint64: true}
		testcases := AsXXXTestCases[map[uint64]bool]{
			{
				Name:           "error.key",
				Value:          value,
				BufferSize:     1,
				PreWriteSize:   1,
				MethodForFixed: method,
			},
			{
				Name:           "error.value",
				Value:          value,
				BufferSize:     10,
				PreWriteSize:   1,
				Contains:       []byte{def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				MethodForFixed: method,
			},
			{
				Name:  "ok",
				Value: value,
				Expected: []byte{
					def.Uint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					def.True,
				},
				BufferSize:     1,
				MethodForFixed: method,
			},
		}
		testcases.Run(t)
	})
}
