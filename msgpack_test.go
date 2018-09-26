package msgpack_test

import (
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

	"github.com/shamaton/msgpack"
	"github.com/shamaton/msgpack/def"
	"github.com/shamaton/msgpack/ext"
)

var now time.Time

func init() {
	n := time.Now()
	now = time.Unix(n.Unix(), int64(n.Nanosecond()))
}

func TestInt(t *testing.T) {
	{
		var r int
		if err := encdec(-8, &r, func(code byte) bool {
			return def.NegativeFixintMin <= int8(code) && int8(code) <= def.NegativeFixintMax
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var r int
		if err := encdec(-108, &r, func(code byte) bool {
			return code == def.Int8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var r int
		if err := encdec(-30108, &r, func(code byte) bool {
			return code == def.Int16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var r int
		if err := encdec(-1030108, &r, func(code byte) bool {
			return code == def.Int32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var r int64
		if err := encdec(int64(math.MinInt64+12345), &r, func(code byte) bool {
			return code == def.Int64
		}); err != nil {
			t.Error(err)
		}
	}

	// error
	{
		var r uint8
		if err := encdec(-8, &r, func(code byte) bool {
			return def.NegativeFixintMin <= int8(code) && int8(code) <= def.NegativeFixintMax
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
	}
	{
		var r int32
		if err := encdec(int64(math.MinInt64+12345), &r, func(code byte) bool {
			return code == def.Int64
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
	}
}

func TestUint(t *testing.T) {
	{
		var v, r uint
		v = 8
		if err := encdec(v, &r, func(code byte) bool {
			return def.PositiveFixIntMin <= uint8(code) && uint8(code) <= def.PositiveFixIntMax
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r uint
		v = 130
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Uint8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r uint
		v = 30130
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Uint16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r uint
		v = 1030130
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Uint32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var r uint64
		if err := encdec(uint64(math.MaxUint64-12345), &r, func(code byte) bool {
			return code == def.Uint64
		}); err != nil {
			t.Error(err)
		}
	}
}
func TestFloat(t *testing.T) {
	{
		var v, r float32
		v = 0
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r float32
		v = -1
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r float32
		v = math.SmallestNonzeroFloat32
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r float32
		v = math.MaxFloat32
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float32
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r float64
		v = 0
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r float64
		v = -1
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r float64
		v = math.SmallestNonzeroFloat64
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r float64
		v = math.MaxFloat64
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err != nil {
			t.Error(err)
		}
	}
	// error
	{
		var v float32
		var r float64
		v = math.MaxFloat32
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float32
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error(err)
		}
	}
	{
		var v float64
		var r float32
		v = math.MaxFloat64
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err == nil || !strings.Contains(err.Error(), "invalid code cb decoding") {
			t.Error("error")
		}
	}
	{
		var v float64
		var r string
		v = math.MaxFloat64
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err == nil || !strings.Contains(err.Error(), "invalid code cb decoding") {
			t.Error("error")
		}
	}
}
func TestBool(t *testing.T) {
	{
		var v, r bool
		v = true
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.True
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r bool
		v = false
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.False
		}); err != nil {
			t.Error(err)
		}
	}
	// error
	{
		var v bool
		var r uint8
		v = true
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.True
		}); err == nil || !strings.Contains(err.Error(), "invalid code c3 decoding") {
			t.Error("error")
		}
	}
}

func TestNil(t *testing.T) {
	{
		var r *map[interface{}]interface{}
		d, err := msgpack.Encode(nil)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Nil {
			t.Error("not nil type")
		}
		err = msgpack.Decode(d, &r)
		if err != nil {
			t.Error(err)
		}
		if r != nil {
			t.Error("not nil")
		}
	}
}

func TestString(t *testing.T) {
	// len 31
	base := "abcdefghijklmnopqrstuvwxyz12345"

	{
		var v, r string
		v = ""
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixStr <= code && code < def.FixStr+32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r string
		v = strings.Repeat(base, 1)
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixStr <= code && code < def.FixStr+32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r string
		v = strings.Repeat(base, 8)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Str8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r string
		v = strings.Repeat(base, (math.MaxUint16/len(base))-1)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Str16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r string
		v = strings.Repeat(base, (math.MaxUint16/len(base))+1)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Str32
		}); err != nil {
			t.Error(err)
		}
	}

	// type different
	{
		var v string
		var r []byte
		v = strings.Repeat(base, 8)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Str8
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
		if v != string(r) {
			t.Error("string error")
		}
	}
	{
		v := []byte(base)
		var r string
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin8
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
		if base != r {
			t.Error("string error")
		}
	}
}
func TestBin(t *testing.T) {
	// slice
	{
		var v, r []byte
		v = make([]byte, 128)
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []byte
		v = make([]byte, 31280)
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []byte
		v = make([]byte, 1031280)
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin32
		}); err != nil {
			t.Error(err)
		}
	}
	// array
	{
		var v, r [128]byte
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r [31280]byte
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin16
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r [1031280]byte
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		}
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v []byte
		var r [1]byte
		v = nil
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Nil
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
	}
	// error
	{
		var v [128]byte
		var r [127]byte
		for i := range v {
			v[i] = byte(rand.Intn(0xff))
		} //%v len is %d, but msgpack has %d elements
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Bin8
		}); err == nil || !strings.Contains(err.Error(), "[127]uint8 len is 127, but msgpack has 128 elements") {
			t.Error("error", err)
		}
	}
}
func TestArray(t *testing.T) {
	// slice
	{
		var v, r []int
		v = make([]int, 15)
		for i := range v {
			v[i] = rand.Intn(math.MaxInt32)
		}
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Array32
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v []int
		var r [1]int
		v = nil
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Nil
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error")
		}
	}
}

func TestFixedSlice(t *testing.T) {
	{
		var v, r []int
		v = []int{-1, 1}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint
		v = []uint{0, 100}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int8
		v = []int8{math.MinInt8, math.MaxInt8}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int16
		v = []int16{math.MinInt16, math.MaxInt16}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int32
		v = []int32{math.MinInt32, math.MaxInt32}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []int64
		v = []int64{math.MinInt64, math.MaxInt64}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		// byte array
		var v, r []uint8
		v = []uint8{0, math.MaxUint8}
		if err := encdec(v, &r, func(code byte) bool {
			return def.Bin8 == code
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint16
		v = []uint16{0, math.MaxUint16}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint32
		v = []uint32{0, math.MaxUint32}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []uint64
		v = []uint64{0, math.MaxUint64}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []float32
		v = []float32{math.SmallestNonzeroFloat32, math.MaxFloat32}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []float64
		v = []float64{math.SmallestNonzeroFloat64, math.MaxFloat64}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []string
		v = []string{"aaa", "bbb"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r []bool
		v = []bool{true, false}
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r map[string]uint
		v = map[string]uint{"a": math.MaxUint32, "b": 0}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]string
		v = map[string]string{"a": "12345", "abcdefghijklmnopqrstuvwxyz": "abcdefghijklmnopqrstuvwxyz"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]float32
		v = map[string]float32{"a": math.MaxFloat32, "b": math.SmallestNonzeroFloat32}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]float64
		v = map[string]float64{"a": math.MaxFloat64, "b": math.SmallestNonzeroFloat64}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]bool
		v = map[string]bool{"a": true, "b": false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int8
		v = map[string]int8{"a": math.MinInt8, "b": math.MaxInt8}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int16
		v = map[string]int16{"a": math.MaxInt16, "b": math.MinInt16}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int32
		v = map[string]int32{"a": math.MaxInt32, "b": math.MinInt32}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]int64
		v = map[string]int64{"a": math.MinInt64, "b": math.MaxInt64}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint8
		v = map[string]uint8{"a": 0, "b": math.MaxUint8}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint16
		v = map[string]uint16{"a": 0, "b": math.MaxUint16}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint32
		v = map[string]uint32{"a": 0, "b": math.MaxUint32}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[string]uint64
		v = map[string]uint64{"a": 0, "b": math.MaxUint64}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int]string
		v = map[int]string{0: "a", 1: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int]bool
		v = map[int]bool{1: true, 2: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint]string
		v = map[uint]string{0: "a", 1: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint]bool
		v = map[uint]bool{0: true, 255: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float32]string
		v = map[float32]string{math.MaxFloat32: "a", math.SmallestNonzeroFloat32: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float32]bool
		v = map[float32]bool{math.SmallestNonzeroFloat32: true, math.MaxFloat32: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float64]string
		v = map[float64]string{math.MaxFloat64: "a", math.SmallestNonzeroFloat64: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[float64]bool
		v = map[float64]bool{math.SmallestNonzeroFloat64: true, math.MaxFloat64: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int8]string
		v = map[int8]string{math.MinInt8: "a", math.MaxInt8: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int8]bool
		v = map[int8]bool{math.MinInt8: true, math.MaxInt8: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int16]string
		v = map[int16]string{math.MaxInt16: "a", math.MinInt16: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int16]bool
		v = map[int16]bool{math.MaxInt16: true, math.MinInt16: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int32]string
		v = map[int32]string{math.MinInt32: "a", math.MaxInt32: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int32]bool
		v = map[int32]bool{math.MinInt32: true, math.MaxInt32: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int64]string
		v = map[int64]string{math.MaxInt64: "a", math.MinInt64: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[int64]bool
		v = map[int64]bool{math.MaxInt64: true, math.MinInt64: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint8]string
		v = map[uint8]string{0: "a", math.MaxUint8: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint8]bool
		v = map[uint8]bool{0: true, math.MaxUint8: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint16]string
		v = map[uint16]string{0: "a", math.MaxUint16: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint16]bool
		v = map[uint16]bool{0: true, math.MaxUint16: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint32]string
		v = map[uint32]string{0: "a", math.MaxUint32: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint32]bool
		v = map[uint32]bool{0: true, math.MaxUint32: false}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint64]string
		v = map[uint64]string{0: "a", math.MaxUint64: "b"}
		if err := encdec(v, &r, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
	}

	{
		var v, r map[uint64]bool
		v = map[uint64]bool{0: true, math.MaxUint64: false}
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext4
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r time.Time
		v = time.Unix(now.Unix(), int64(now.Nanosecond()))
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext8
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestMap(t *testing.T) {
	{
		var v, r map[int]int
		v = map[int]int{1: 2, 3: 4, 5: 6, 7: 8, 9: 10}
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, &r, func(code byte) bool {
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
		d, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Map16 {
			t.Error("code diffenrent")
		}
		err = msgpack.Decode(d, &r)
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
		d, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Map16 {
			t.Error("code diffenrent")
		}
		err = msgpack.Decode(d, &r)
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
		d, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Map16 {
			t.Error("code diffenrent")
		}
		err = msgpack.Decode(d, &r)
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
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Uint8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r *int
		d, err := msgpack.Encode(v)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Nil {
			t.Error("code diffenrent")
		}
		err = msgpack.Decode(d, &r)
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
		if err := encdec(v, r, func(code byte) bool {
			return code == def.Nil
		}); err == nil || !strings.Contains(err.Error(), "holder must set pointer value. but got:") {
			t.Error(err)
		}
	}
}

func TestUnsupported(t *testing.T) {
	{
		var v, r complex128
		v = 1i
		_, err := msgpack.Encode(v)
		if !strings.Contains(err.Error(), "type(complex128) is unsupported") {
			t.Error("test error")
		}
		err = msgpack.Decode([]byte{0xc0}, &r)
		if !strings.Contains(err.Error(), "type(complex128) is unsupported") {
			t.Error("test error")
		}
	}
}

/////////////////////////////////////////////////////////////////

func TestStruct(t *testing.T) {
	testSturctCode(t)
	testStructTag(t)
	testStructArray(t)

	testStructUseCase(t)
	msgpack.StructAsArray = true
	testStructUseCase(t)
}

func testStructTag(t *testing.T) {
	type vSt struct {
		One int     `msgpack:"Three"`
		Two string  `msgpack:"four"`
		Ten float32 `msgpack:"ignore"`
	}
	type rSt struct {
		Three int
		Four  string `msgpack:"four"`
		Ten   float32
	}

	msgpack.StructAsArray = false

	v := vSt{One: 1, Two: "2", Ten: 1.234}
	r := rSt{}

	d, err := msgpack.EncodeStructAsMap(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.FixMap+0x02 {
		t.Error("code different")
	}
	err = msgpack.DecodeStructAsMap(d, &r)
	if err != nil {
		t.Error(err)
	}
	if v.One != r.Three || v.Two != r.Four || r.Ten != 0 {
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

	d, err := msgpack.EncodeStructAsArray(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.FixArray+0x04 {
		t.Error("code different")
	}
	err = msgpack.DecodeStructAsArray(d, &r)
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
		if err := encdec(v1, &r1, func(code byte) bool {
			return def.FixMap <= code && code <= def.FixMap+0x0f
		}); err != nil {
			t.Error(err)
		}
		if err := encdec(v16, &r16, func(code byte) bool {
			return code == def.Map16
		}); err != nil {
			t.Error(err)
		}
	}
	msgpack.StructAsArray = true
	{
		var r1 st1
		var r16 st16
		if err := encdec(v1, &r1, func(code byte) bool {
			return def.FixArray <= code && code <= def.FixArray+0x0f
		}); err != nil {
			t.Error(err)
		}
		if err := encdec(v16, &r16, func(code byte) bool {
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
		Time        time.Time
		Duration    time.Duration
		Child       child
		Child3Array []child3
	}
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

func encSt(t *testing.T, in interface{}, isDebug bool) ([]byte, []byte, error) {
	var d1, d2 []byte
	var err error
	d1, err = msgpack.Encode(in)
	if err != nil {
		return nil, nil, err
	}

	if msgpack.StructAsArray {
		d2, err = msgpack.EncodeStructAsMap(in)
	} else {
		d2, err = msgpack.EncodeStructAsArray(in)
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
	if err := msgpack.Decode(d1, out1); err != nil {
		return err
	}

	if msgpack.StructAsArray {
		err := msgpack.DecodeStructAsMap(d2, out2)
		if err != nil {
			return err
		}
	} else {
		err := msgpack.DecodeStructAsArray(d2, out2)
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
		v := ExtInt{V: 321}
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
		v := ExtInt{V: 123}
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
}

type ExtStruct struct {
	V int
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
	if code == def.Fixext4 {
		t, _ := td.ReadSize1(offset, d)
		return int8(t) == td.Code()
	}
	return false
}

func (td *testDecoder) AsValue(offset int, k reflect.Kind, d *[]byte) (interface{}, int, error) {
	code, offset := td.ReadSize1(offset, d)

	switch code {
	case def.Fixext4:
		_, offset = td.ReadSize1(offset, d)
		bs, offset := td.ReadSize4(offset, d)
		return ExtInt{V: int(binary.BigEndian.Uint32(bs))}, offset, nil
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
	return def.Byte1 + def.Byte4, nil
}

func (s *testEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	t := value.Interface().(ExtInt)
	offset = s.SetByte1Int(def.Fixext4, offset, bytes)
	offset = s.SetByte1Int(int(s.Code()), offset, bytes)
	offset = s.SetByte4Int(t.V, offset, bytes)
	return offset
}

/////////////////////////////////////////////////////////

func encdec(v, r interface{}, j func(d byte) bool) error {
	d, err := msgpack.Encode(v)
	if err != nil {
		return err
	}
	if !j(d[0]) {
		return fmt.Errorf("different %s", hex.Dump(d))
	}
	if err := msgpack.Decode(d, r); err != nil {
		return err
	}
	if err := equalCheck(v, r); err != nil {
		return err
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
