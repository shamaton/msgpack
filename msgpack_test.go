package msgpack_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/shamaton/msgpack"
)

func TestMsgpack(t *testing.T) {
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
