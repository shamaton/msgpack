package msgpack_test

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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
			t.Error("error", err)
		}
	}
	{
		var r int32
		if err := encdec(int64(math.MinInt64+12345), &r, func(code byte) bool {
			return code == def.Int64
		}); err == nil || !strings.Contains(err.Error(), "value different") {
			t.Error("error", err)
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
	{
		var r int
		v := float32(2.345)
		b, err := msgpack.MarshalAsArray(v)
		if err != nil {
			t.Error(err)
		}

		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}

		if r != 2 {
			t.Error("different value", r)
		}
	}
	{
		var r int
		v := 6.789
		b, err := msgpack.MarshalAsArray(v)
		if err != nil {
			t.Error(err)
		}

		err = msgpack.UnmarshalAsArray(b, &r)
		if err != nil {
			t.Error(err)
		}

		if r != 6 {
			t.Error("different value", r)
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
			t.Error("error", err)
		}
	}
	{
		var v float64
		var r string
		v = math.MaxFloat64
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Float64
		}); err == nil || !strings.Contains(err.Error(), "invalid code cb decoding") {
			t.Error("error", err)
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
			t.Error("error", err)
		}
	}
}

func TestNil(t *testing.T) {
	{
		var r *map[interface{}]interface{}
		d, err := msgpack.Marshal(nil)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Nil {
			t.Error("not nil type")
		}
		err = msgpack.Unmarshal(d, &r)
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

func TestComplex(t *testing.T) {
	{
		var v, r complex64
		v = complex(1, 2)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext8
		}); err != nil {
			t.Error(err)
		}

		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if int8(b[1]) != def.ComplexTypeCode() {
			t.Errorf("complex type code is different %d, %d", int8(b[1]), def.ComplexTypeCode())
		}

		var rr complex128
		err = msgpack.Unmarshal(b, &rr)
		if err != nil {
			t.Error(err)
		}
		if imag(rr) == 0 || real(rr) == 0 {
			t.Errorf("somthing wrong %v", rr)
		}

		err = msgpack.Unmarshal([]byte{def.Nil}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if err != nil && !strings.Contains(err.Error(), "should not reach this line") {
			t.Error(err)
		}

		typeCode := int8(-99)
		msgpack.SetComplexTypeCode(typeCode)

		err = msgpack.Unmarshal(b, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if err != nil && !strings.Contains(err.Error(), "fixext8") {
			t.Error(err)
		}

		err = msgpack.Unmarshal(b, &rr)
		if err == nil {
			t.Errorf("error must occur")
		}
		if err != nil && !strings.Contains(err.Error(), "fixext8") {
			t.Error(err)
		}
	}

	typeCode := int8(-100)
	msgpack.SetComplexTypeCode(typeCode)
	{
		if def.ComplexTypeCode() != typeCode {
			t.Errorf("complex type code not set %d, %d", typeCode, def.ComplexTypeCode())
		}

		var v, r complex128
		v = complex(math.MaxFloat64, math.SmallestNonzeroFloat64)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext16
		}); err != nil {
			t.Error(err)
		}

		b, err := msgpack.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		if int8(b[1]) != def.ComplexTypeCode() {
			t.Errorf("complex type code is different %d, %d", int8(b[1]), def.ComplexTypeCode())
		}

		var rr complex64
		err = msgpack.Unmarshal(b, &rr)
		if err != nil {
			t.Error(err)
		}
		if imag(rr) != 0 || real(rr) == 0 {
			t.Errorf("somthing wrong %v", rr)
		}

		err = msgpack.Unmarshal([]byte{def.Nil}, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if err != nil && !strings.Contains(err.Error(), "should not reach this line") {
			t.Error(err)
		}

		typeCode := int8(-99)
		msgpack.SetComplexTypeCode(typeCode)

		err = msgpack.Unmarshal(b, &r)
		if err == nil {
			t.Errorf("error must occur")
		}
		if err != nil && !strings.Contains(err.Error(), "fixext16") {
			t.Error(err)
		}

		err = msgpack.Unmarshal(b, &rr)
		if err == nil {
			t.Errorf("error must occur")
		}
		if err != nil && !strings.Contains(err.Error(), "fixext16") {
			t.Error(err)
		}
	}
}

func TestInterface(t *testing.T) {
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
			t.Error(i, v, err)
		}
	}

	// error
	var r interface{}
	err := msgpack.Unmarshal([]byte{def.Ext32}, &r)
	if err == nil {
		t.Error("error must occur")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid code") {
		t.Error(err)
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
		v = nil
		if err := encdec(v, &r, func(code byte) bool {
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
		_, err = extTime.Decoder.AsValue(notReach, reflect.Bool)
		if err == nil || !strings.Contains(err.Error(), "should not reach this line") {
			t.Error("something wrong", err)
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
		if err := encdec(v, &r, func(code byte) bool {
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
		if err := encdec(v, r, func(code byte) bool {
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
	testStructCode(t)
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

func testStructCode(t *testing.T) {
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
			Byte2Int:    rand.Intn(math.MaxUint16) - math.MinInt16,
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
			Byte2Int:    rand.Intn(math.MaxUint16) - math.MinInt16,
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

type testDecoder struct{}

var extIntCode = int8(-2)

func (td *testDecoder) Code() int8 {
	return extIntCode
}

func (td *testDecoder) AsValue(data []byte, k reflect.Kind) (interface{}, error) {
	readN := func(n int) []byte {
		res := data[:n]
		data = data[n:]
		return res
	}

	switch len(data) {
	case 15 + 15 + 10 + 3:
		i8 := readN(1)[0]
		i16 := readN(2)
		i32 := readN(4)
		i64 := readN(8)
		u8 := readN(1)[0]
		u16 := readN(2)
		u32 := readN(4)
		u64 := readN(8)
		b16 := readN(2)
		b32 := readN(4)
		bu32 := readN(4)
		bs := readN(3)
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
		}, nil
	}

	return ExtInt{}, fmt.Errorf("should not reach this line!! data %x decoding %v", data, k)
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

func (s *testEncoder) WriteToBytes(value reflect.Value, writer io.Writer) error {
	t := value.Interface().(ExtInt)

	err := s.SetByte1Int(def.Ext8, writer)
	if err != nil {
		return err
	}

	err = s.SetByte1Int(15+15+10+len(t.Bytes), writer)
	if err != nil {
		return err
	}

	err = s.SetByte1Int(int(s.Code()), writer)
	if err != nil {
		return err
	}

	err = s.SetByte1Int64(int64(t.Int8), writer)
	if err != nil {
		return err
	}

	err = s.SetByte2Int64(int64(t.Int16), writer)
	if err != nil {
		return err
	}

	err = s.SetByte4Int64(int64(t.Int32), writer)
	if err != nil {
		return err
	}

	err = s.SetByte8Int64(t.Int64, writer)
	if err != nil {
		return err
	}

	err = s.SetByte1Uint64(uint64(t.Uint8), writer)
	if err != nil {
		return err
	}

	err = s.SetByte2Uint64(uint64(t.Uint16), writer)
	if err != nil {
		return err
	}

	err = s.SetByte4Uint64(uint64(t.Uint32), writer)
	if err != nil {
		return err
	}

	err = s.SetByte8Uint64(t.Uint64, writer)
	if err != nil {
		return err
	}

	err = s.SetByte2Int(t.Byte2Int, writer)
	if err != nil {
		return err
	}

	err = s.SetByte4Int(t.Byte4Int, writer)
	if err != nil {
		return err
	}

	err = s.SetByte4Uint32(t.Byte4Uint32, writer)
	if err != nil {
		return err
	}

	err = s.SetBytes(t.Bytes, writer)
	if err != nil {
		return err
	}

	return nil
}

/////////////////////////////////////////////////////////

type Ext2Struct struct {
	V int
}
type Ext2Int Ext2Struct

type testExt2Decoder struct{}

func (td *testExt2Decoder) Code() int8 {
	return 3
}

func (td *testExt2Decoder) AsValue(data []byte, k reflect.Kind) (interface{}, error) {
	return Ext2Int{}, fmt.Errorf("should not reach this line!! code %x data %x decoding %v", 3, data, k)
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

func (s *testExt2Encoder) WriteToBytes(value reflect.Value, writer io.Writer) error {
	return nil
}

/////////////////////////////////////////////////////////

func encdec(v, r interface{}, j func(d byte) bool) error {
	d, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}
	if !j(d[0]) {
		return fmt.Errorf("different %s", hex.Dump(d))
	}
	if err := msgpack.Unmarshal(d, r); err != nil {
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
