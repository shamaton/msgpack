package serialize

import (
	"encoding/hex"
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

type serializer struct {
	common
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
	s.create(rv, 0)
	fmt.Println(def.Byte4, size, s)
	fmt.Println(hex.Dump(s.d))
}

func (s *serializer) calcSize(rv reflect.Value) int {
	ret := 0

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		ret = s.calcUint(v)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		if v >= 0 {
			s.calcUint(uint64(v))
		} else if s.isNegativeFixInt64(v) {
			ret = def.Byte1
		} else if v >= math.MinInt8 {
			ret = def.Byte1 + def.Byte1
		} else if v >= math.MinInt16 {
			ret = def.Byte1 + def.Byte2
		} else if v >= math.MinInt32 {
			ret = def.Byte1 + def.Byte4
		} else {
			ret = def.Byte1 + def.Byte8
		}
	}

	return ret
}

func (s *serializer) create(rv reflect.Value, offset int) {

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		offset = s.writeUint(v, offset)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		if v >= 0 {
			s.writeUint(uint64(v), offset)
		} else if s.isNegativeFixInt64(v) {
			offset = s.writeSize1Int64(v, offset)
		} else if v >= math.MinInt8 {
			offset = s.writeSize1Int(def.Int8, offset)
			offset = s.writeSize1Int64(v, offset)
		} else if v >= math.MinInt16 {
			offset = s.writeSize1Int(def.Int16, offset)
			offset = s.writeSize2Int64(v, offset)
		} else if v >= math.MinInt32 {
			offset = s.writeSize1Int(def.Int32, offset)
			offset = s.writeSize4Int64(v, offset)
		} else {
			offset = s.writeSize1Int(def.Int64, offset)
			offset = s.writeSize8Int64(v, offset)
		}
	}
}

func (s *serializer) calcUint(v uint64) int {
	size := def.Byte1
	if v <= math.MaxInt8 {
	} else if v <= math.MaxUint8 {
		size += def.Byte1
	} else if v <= math.MaxUint16 {
		size += def.Byte2
	} else if v <= math.MaxUint32 {
		size += def.Byte4
	} else {
		size += def.Byte8
	}
	return size
}

func (s *serializer) writeUint(v uint64, offset int) int {
	if v <= math.MaxInt8 {
		offset = s.writeSize1Uint64(v, offset)
	} else if v <= math.MaxUint8 {
		offset = s.writeSize1Int(def.Uint8, offset)
		offset = s.writeSize1Uint64(v, offset)
	} else if v <= math.MaxUint16 {
		offset = s.writeSize1Int(def.Uint16, offset)
		offset = s.writeSize2Uint64(v, offset)
	} else if v <= math.MaxUint32 {
		offset = s.writeSize1Int(def.Uint32, offset)
		offset = s.writeSize4Uint64(v, offset)
	} else {
		offset = s.writeSize1Int(def.Uint64, offset)
		offset = s.writeSize8Uint64(v, offset)
	}
	return offset
}
