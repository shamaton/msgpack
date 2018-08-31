package msgpack_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shamaton/msgpack"
)

func TestMsgpack(t *testing.T) {
	testSimple(t)
	testStruct(t)
	msgpack.StructAsArray = true
	testSimple(t)
	testStruct(t)
}
func testSimple(t *testing.T) {
	uints := []uint64{1, math.MaxInt8, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}

	for _, v := range uints {
		var r uint64
		_, err := checkRoutine(t, v, &r, false)
		if err != nil || v != r {
			t.Error(err)
		}
	}

	ints := []int64{-1, -32, -33, math.MinInt8, math.MinInt16, math.MinInt32, math.MinInt64}

	for _, v := range ints {
		var r int64
		_, err := checkRoutine(t, v, &r, false)
		if err != nil || v != r {
			t.Error(err)
		}
	}

	{
		vs := []float64{0, -1, math.MaxFloat32, math.MaxFloat64, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat64}

		for _, v := range vs {
			var r float64
			_, err := checkRoutine(t, v, &r, true)
			if err != nil || v != r {
				t.Error(err)
			}
		}
	}
	{
		ls := []int{0, 32, math.MaxUint8, math.MaxUint16, math.MaxUint16 + 32}
		base := "abcdefghijklmnopqrstuvwxyz12345"
		for _, l := range ls {
			v := strings.Repeat(base, l/len(base))
			var r string
			_, err := checkRoutine(t, v, &r, false)
			if err != nil || v != r {
				t.Error(err)
			}
		}
	}

	{
		vs := []bool{true, false}
		for _, v := range vs {
			var r *bool
			_, err := checkRoutine(t, v, &r, false)
			if err != nil || v != *r {
				t.Error(err)
			}
		}
	}
	{
		v := time.Now()
		var r time.Time
		_, err := checkRoutine(t, v, &r, true)
		if err != nil {
			t.Error(err)
		}
		if err := equalCheck(v, r); err != nil {
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
		Time:   time.Now(),
		//Duration:   time.Duration(123 * time.Second),

		// child
		Child: child{
			IntArray:    []int{-1, -2, -3, -4, -5},
			UintArray:   []uint{1, 2, 3, 4, 5},
			FloatArray:  []float32{-1.2, -3.4, -5.6, -7.8},
			BoolArray:   []bool{true, true, false, false, true},
			StringArray: []string{"str", "ing", "arr", "ay"},
			TimeArray:   []time.Time{time.Now(), time.Now(), time.Now()},
			//DurationArray:   []time.Duration{time.Duration(1 * time.Nanosecond), time.Duration(2 * time.Nanosecond)},

			// childchild
			Child: child2{
				Int2Uint:   map[int]uint{-1: 2, -3: 4},
				Float2Bool: map[float32]bool{-1.1: true, -2.2: false},
				//Duration2Struct: map[time.Duration]child3{time.Duration(1 * time.Hour): child3{Int: 1}, time.Duration(2 * time.Hour): child3{Int: 2}},
			},
		},
	}
	rSt := st{}
	if _, err := checkRoutine(t, vSt, &rSt, false); err != nil {
		t.Error(err)
	}
	if err := equalCheck(vSt, rSt); err != nil {
		t.Error(err)
	}

	// pointer
	prSt := new(st)
	if _, err := checkRoutine(t, &vSt, &prSt, false); err != nil {
		t.Error(err)
	}
	if err := equalCheck(vSt, prSt); err != nil {
		t.Error(err)
	}
}

func TestMsgpackMultiByte(t *testing.T) {

}

func checkRoutine(t *testing.T, in interface{}, out interface{}, isDebug bool) ([]byte, error) {
	d, err := msgpack.Encode(in)
	if err != nil {
		return nil, err
	}

	if isDebug {
		t.Log(in, " -- to byte --> ", d)
	}

	if err := msgpack.Decode(d, out); err != nil {
		return nil, err
	}

	i := getValue(in)
	o := getValue(out)
	if isDebug {
		t.Log("value [in]:", i, " [out]:", o)
	}
	return d, nil
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
		return errors.New(fmt.Sprint("value different [in]:", in, " [out]:", out))
	}
	return nil
}
