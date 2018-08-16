package serialize

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"time"

	"github.com/shamaton/msgpack/def"
)

type serializer struct {
	common
	asArray bool
}

var now = time.Now()

func AsArray(v interface{}, asArray bool) ([]byte, error) {
	s := serializer{asArray: asArray}
	defer s.recover()

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}
	size, err := s.calcSize(rv)
	if err != nil {
		return nil, err
	}

	s.d = make([]byte, size)
	last, err := s.create(rv, 0)
	if err != nil {
		return nil, err
	}
	if size != last {
		return nil, fmt.Errorf("failed serialization size=%d, lastIdx=%d", size, last)
	}
	return s.d, err
}

func (s *serializer) recover() {
	err := recover()
	if err != nil {
		for depth := 0; ; depth++ {
			_, file, line, ok := runtime.Caller(depth)
			if !ok {
				break
			}
			log.Printf("======> %d: %v:%d", depth, file, line)
		}
	}
}

func (s *serializer) calcSize(rv reflect.Value) (int, error) {
	ret := def.Byte1

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		ret += s.calcUint(v)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		ret += s.calcInt(int64(v))

	case reflect.Float32, reflect.Float64:
		v := rv.Float()
		if math.SmallestNonzeroFloat32 <= v && v <= math.MaxFloat32 {
			ret += def.Byte4
		} else {
			ret += def.Byte8
		}

	case reflect.String:
		ret += s.calcString(rv.String())

	case reflect.Bool:
		// do nothing

	case reflect.Array, reflect.Slice:
		if rv.IsNil() {
			return ret, nil
		}
		l := rv.Len()
		// bin format
		if s.isByteSlice(rv) {
			r, err := s.calcByteSlice(l)
			if err != nil {
				return 0, err
			}
			ret += r
			return ret, nil
		}

		// format size
		if l <= 0x0f {
			// format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if l <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}

		if size, find := s.calcFixedSlice(rv); find {
			ret += size
			return ret, nil
		}

		// objects size
		for i := 0; i < l; i++ {
			s, err := s.calcSize(rv.Index(i))
			if err != nil {
				return 0, err
			}
			ret += s
		}

	case reflect.Map:
		if rv.IsNil() {
			return ret, nil
		}

		l := rv.Len()
		// format
		if l <= 0x0f {
			// do nothing
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if l <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this map length : %d", l)
		}

		if size, find := s.calcFixedMap(rv); find {
			ret += size
			return ret, nil
		}

		// key-value
		keys := rv.MapKeys()
		for _, k := range keys {
			keySize, err := s.calcSize(k)
			if err != nil {
				return 0, err
			}
			valueSize, err := s.calcSize(rv.MapIndex(k))
			if err != nil {
				return 0, err
			}
			ret += keySize + valueSize
		}

	case reflect.Struct:
		if isTime, tm := s.isDateTime(rv); isTime {
			size := s.calcTime(tm)
			ret += size
			return ret, nil
		}

		var size int
		var err error
		if s.asArray {
			size, err = s.calcStructArray(rv)
		} else {
			size, err = s.calcStructMap(rv)
		}
		if err != nil {
			return 0, err
		}
		ret += size

	case reflect.Ptr:
		if rv.IsNil() {
			return ret, nil
		}
		size, err := s.calcSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret += size

	case reflect.Invalid:
		// do nothing (return nil)
	}

	return ret, nil
}

func (s *serializer) create(rv reflect.Value, offset int) (int, error) {

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		offset = s.writeUint(v, offset)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		offset = s.writeInt(v, offset)

	case reflect.Float32, reflect.Float64:
		v := rv.Float()
		if math.SmallestNonzeroFloat32 <= v && v <= math.MaxFloat32 {
			offset = s.setByte1Int(def.Float32, offset)
			offset = s.setByte4Uint64(uint64(math.Float32bits(float32(v))), offset)
		} else {
			offset = s.setByte1Int(def.Float64, offset)
			offset = s.setByte8Uint64(math.Float64bits(v), offset)
		}

	case reflect.Bool:
		if rv.Bool() {
			offset = s.setByte1Int(def.True, offset)
		} else {
			offset = s.setByte1Int(def.False, offset)
		}

	case reflect.String:
		offset = s.writeString(rv.String(), offset)

	case reflect.Array, reflect.Slice:
		if rv.IsNil() {
			return s.writeNil(offset)
		}
		l := rv.Len()
		// bin format
		if s.isByteSlice(rv) {
			offset = s.writeByteSliceLength(l, offset)
			offset = s.setBytes(rv.Bytes(), offset)
			return offset, nil
		}

		// format
		offset = s.writeSliceLength(l, offset)

		if offset, find := s.writeFixedSlice(rv, offset); find {
			return offset, nil
		}

		// objects
		for i := 0; i < l; i++ {
			var err error
			offset, err = s.create(rv.Index(i), offset)
			if err != nil {
				return 0, err
			}
		}

	case reflect.Map:
		if rv.IsNil() {
			return s.writeNil(offset)
		}

		l := rv.Len()
		offset = s.writeMapLength(l, offset)

		if offset, find := s.writeFixedMap(rv, offset); find {
			return offset, nil
		}

		// key-value
		keys := rv.MapKeys()
		for _, k := range keys {
			o, err := s.create(k, offset)
			if err != nil {
				return 0, err
			}
			o, err = s.create(rv.MapIndex(k), o)
			if err != nil {
				return 0, err
			}
			offset = o
		}

	case reflect.Struct:
		if isTime, tm := s.isDateTime(rv); isTime {
			return s.writeTime(tm, offset)
		}
		var o int
		var err error
		if s.asArray {
			o, err = s.writeStructArray(rv, offset)
		} else {
			o, err = s.writeStructMap(rv, offset)
		}
		if err != nil {
			return 0, err
		}
		offset = o

	case reflect.Ptr:
		if rv.IsNil() {
			return s.writeNil(offset)
		}

		o, err := s.create(rv.Elem(), offset)
		if err != nil {
			return 0, err
		}
		offset = o

	case reflect.Invalid:
		return s.writeNil(offset)

	}
	return offset, nil
}
