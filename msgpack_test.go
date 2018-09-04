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

func TestTime(t *testing.T) {
	now := time.Now()
	{
		var v, r time.Time
		v = time.Unix(int64(now.Second()), 0)
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext4
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r time.Time
		v = time.Unix(int64(now.Second()), int64(now.Nanosecond()))
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext8
		}); err != nil {
			t.Error(err)
		}
	}
	{
		var v, r time.Time
		v = time.Now()
		t.Log(v.Unix(), v.UnixNano(), now.Unix(), now.UnixNano())
		if err := encdec(v, &r, func(code byte) bool {
			return code == def.Fixext8
		}); err != nil {
			t.Error(err)
		}
	}
}

func TestMap(t *testing.T) {
}
func TestStruct(t *testing.T) {
	testStruct(t)
	msgpack.StructAsArray = true
	testStruct(t)
}
func TestMsgpack(t *testing.T) {
	testSimple(t)
	testStruct(t)
	msgpack.StructAsArray = true
	testSimple(t)
	testStruct(t)
}
func testSimple(t *testing.T) {
	{
		uints := []uint64{1, math.MaxInt8, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}

		for _, v := range uints {
			var r1, r2 uint64
			d1, d2, err := encode(t, v, false)
			if err != nil {
				t.Error(err)
			}
			err = decode(t, d1, d2, &r1, &r2, false)
			if err != nil || !(v == r1 && r1 == r2) {
				t.Error(err, v, r1, r2)
			}
		}
	}

	{
		ints := []int64{-1, -32, -33, math.MinInt8, math.MinInt16, math.MinInt32, math.MinInt64}

		for _, v := range ints {
			var r1, r2 int64
			d1, d2, err := encode(t, v, false)
			if err != nil {
				t.Error(err)
			}
			err = decode(t, d1, d2, &r1, &r2, false)
			if err != nil || !(v == r1 && r1 == r2) {
				t.Error(err)
			}
		}
	}

	{
		vs := []float64{0, -1, math.MaxFloat32, math.MaxFloat64, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat64}

		for _, v := range vs {
			var r1, r2 float64
			d1, d2, err := encode(t, v, false)
			if err != nil {
				t.Error(err)
			}
			err = decode(t, d1, d2, &r1, &r2, false)
			if err != nil || !(v == r1 && r1 == r2) {
				t.Error(err)
			}
		}
	}
	{
		ls := []int{0, 32, math.MaxUint8, math.MaxUint16, math.MaxUint16 + 32}
		base := "abcdefghijklmnopqrstuvwxyz12345"
		for _, l := range ls {
			v := strings.Repeat(base, l/len(base))
			var r1, r2 string
			d1, d2, err := encode(t, v, false)
			if err != nil {
				t.Error(err)
			}
			err = decode(t, d1, d2, &r1, &r2, false)
			if err != nil || !(v == r1 && r1 == r2) {
				t.Error(err)
			}
		}
	}

	{
		vs := []bool{true, false}
		for _, v := range vs {
			var r1, r2 *bool
			d1, d2, err := encode(t, v, false)
			if err != nil {
				t.Error(err)
			}
			err = decode(t, d1, d2, &r1, &r2, false)
			if err != nil || !(v == *r1 && *r1 == *r2) {
				t.Error(err)
			}
		}
	}
	{

		v := now
		var r1, r2 time.Time
		d1, d2, err := encode(t, v, false)
		if err != nil {
			t.Error(err)
		}
		err = decode(t, d1, d2, &r1, &r2, false)
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
	{
		v := 1
		var r1, r2 int
		d1, d2, err := encode(t, v, false)
		if err != nil {
			t.Error(err)
		}
		err = decode(t, d1, d2, r1, r2, false)
		if err == nil || !strings.Contains(err.Error(), "holder must set pointer value. but got:") {
			t.Error(err)
		}

	}
}

func testStruct(t *testing.T) {
	type child3 struct {
		Int int
	}
	type child2 struct {
		Int2Uint   map[int]uint
		Float2Bool map[float32]bool
		//Duration2Struct map[time.Duration]child3
	}
	type child struct {
		IntArray    []int
		UintArray   []uint
		FloatArray  []float32
		BoolArray   []bool
		StringArray []string
		TimeArray   []time.Time
		Child       child2
	}
	type st struct {
		Int8   int8
		Int16  int16
		Int32  int32
		Int64  int64
		Uint8  byte
		Uint16 uint16
		Uint32 uint32
		Uint64 uint64
		Float  float32
		Double float64
		Bool   bool
		String string
		Time   time.Time
		//Duration   time.Duration
		Child child
	}
	vSt := &st{
		Int32:  -32,
		Int8:   -8,
		Int16:  -16,
		Int64:  -64,
		Uint32: 32,
		Uint8:  8,
		Uint16: 16,
		Uint64: 64,
		Float:  1.23,
		Double: 2.3456,
		Bool:   true,
		String: "Parent",
		Time:   now,
		//Duration:   time.Duration(123 * time.Second),

		// child
		Child: child{
			IntArray:    []int{-1, -2, -3, -4, -5},
			UintArray:   []uint{1, 2, 3, 4, 5},
			FloatArray:  []float32{-1.2, -3.4, -5.6, -7.8},
			BoolArray:   []bool{true, true, false, false, true},
			StringArray: []string{"str", "ing", "arr", "ay"},
			TimeArray:   []time.Time{now, now, now},
			//DurationArray:   []time.Duration{time.Duration(1 * time.Nanosecond), time.Duration(2 * time.Nanosecond)},

			// childchild
			Child: child2{
				Int2Uint:   map[int]uint{-1: 2, -3: 4},
				Float2Bool: map[float32]bool{-1.1: true, -2.2: false},
				//Duration2Struct: map[time.Duration]child3{time.Duration(1 * time.Hour): child3{Int: 1}, time.Duration(2 * time.Hour): child3{Int: 2}},
			},
		},
	}
	{

		r1, r2 := st{}, st{}
		d1, d2, err := encode(t, vSt, false)
		if err != nil {
			t.Error(err)
		}
		err = decode(t, d1, d2, &r1, &r2, false)
		if err != nil {
			t.Error(err)
		}
		if err := equalCheck(vSt, r1); err != nil {
			t.Error(err)
		}
		if err := equalCheck(vSt, r2); err != nil {
			t.Error(err)
		}
	}

	// pointer
	{
		r1, r2 := new(st), new(st)
		d1, d2, err := encode(t, vSt, false)
		if err != nil {
			t.Error(err)
		}
		err = decode(t, d1, d2, &r1, &r2, false)
		if err != nil {
			t.Error(err)
		}
		if err := equalCheck(vSt, r1); err != nil {
			t.Error(err)
		}
		if err := equalCheck(vSt, r2); err != nil {
			t.Error(err)
		}
	}
}

func TestMsgpackExt(t *testing.T) {
	msgpack.AddExtCoder(encoder, decoder)

	{
		v := ExtInt{V: 321}
		var r1, r2 ExtInt
		d1, d2, err := encode(t, v, true)
		if err != nil {
			t.Error(err)
		}
		err = decode(t, d1, d2, &r1, &r2, false)
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
	msgpack.RemoveExtCoder(encoder, decoder)
	{
		v := ExtInt{V: 123}
		var r1, r2 ExtInt
		d1, d2, err := encode(t, v, true)
		if err != nil {
			t.Error(err)
		}
		err = decode(t, d1, d2, &r1, &r2, false)
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

func encode(t *testing.T, in interface{}, isDebug bool) ([]byte, []byte, error) {
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

func decode(t *testing.T, d1, d2 []byte, out1, out2 interface{}, isDebug bool) error {
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

///////
type ExtStruct struct {
	V int
}
type ExtInt ExtStruct

var decoder = new(testDecoder)

type testDecoder struct {
	ext.DecoderCommon
}

func (td *testDecoder) IsType(offset int, d *[]byte) bool {
	code, offset := td.ReadSize1(offset, d)
	fmt.Println("ok!!!!!!!!!!!!!")
	if code == def.Fixext4 {
		t, _ := td.ReadSize1(offset, d)
		return int8(t) == -2
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

func (s *testEncoder) IsType(value reflect.Value) bool {
	fmt.Println("ok2!!!!!!!!!!!!!")
	_, ok := value.Interface().(ExtInt)
	return ok
}

func (s *testEncoder) CalcByteSize(value reflect.Value) (int, error) {
	return def.Byte1 + def.Byte4, nil
}

func (s *testEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	t := value.Interface().(ExtInt)
	offset = s.SetByte1Int(def.Fixext4, offset, bytes)
	offset = s.SetByte1Int(-2, offset, bytes)
	offset = s.SetByte4Int(t.V, offset, bytes)
	return offset
}
