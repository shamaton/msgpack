package serialize

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

type serializer struct {
	common
	d []byte
}

func AsArray(v interface{}) /*([]byte, error)*/ {
	s := serializer{}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}

	size := s.calcSize(rv)
	s.d = make([]byte, size)
	s.create(rv)
	fmt.Println(def.Byte1, size, s)
	fmt.Println(hex.Dump(s.d))
}

func (s *serializer) calcSize(rv reflect.Value) int {
	ret := 0

	switch rv.Kind() {
	case reflect.Int8:
		val := rv.Int()
		if s.isPositiveFixInt64(val) || s.isNegativeFixInt64(val) {
			ret = def.Byte1
		} else {
			ret = def.Byte1 + def.Byte1
		}

	case reflect.Int16:
		val := rv.Int()
		if s.isPositiveFixInt64(val) || s.isNegativeFixInt64(val) {
			ret = def.Byte1
		} else {
			ret = def.Byte1 + def.Byte2
		}

	case reflect.Int32:
		val := rv.Int()
		if s.isPositiveFixInt64(val) || s.isNegativeFixInt64(val) {
			ret = def.Byte1
		} else {
			ret = def.Byte1 + def.Byte4
		}

	case reflect.Int64:
		val := rv.Int()
		if s.isPositiveFixInt64(val) || s.isNegativeFixInt64(val) {
			ret = def.Byte1
		} else {
			ret = def.Byte1 + def.Byte8
		}

	case reflect.Int:
		val := rv.Int()
		if s.isPositiveFixInt64(val) || s.isNegativeFixInt64(val) {
			ret = def.Byte1
		} else {
			if def.IntSize == 32 {
				ret = def.Byte1 + def.Byte4
			} else {
				ret = def.Byte1 + def.Byte8
			}
		}
	}

	return ret
}

func (s *serializer) create(rv reflect.Value) {

	switch rv.Kind() {
	case reflect.Int8:
		val := rv.Int()
		if s.isPositiveFixInt64(val) || s.isNegativeFixInt64(val) {
			s.d[0] = byte(val)
		} else {
			s.d[0] = def.Int8
			s.d[1] = byte(val)
		}

	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Int:
	}
}
