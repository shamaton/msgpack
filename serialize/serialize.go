package serialize

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"time"

	"github.com/shamaton/msgpack/def"
)

type serializer struct {
	d       []byte
	asArray bool
}

var now = time.Now()

func Exec(v interface{}, asArray bool) (b []byte, err error) {
	s := serializer{asArray: asArray}
	/*
		defer func() {
			e := recover()
			if e != nil {
				b = nil
				err = fmt.Errorf("unexpected error!! \n%s", stackTrace())
			}
		}()
	*/

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
	last := s.create(rv, 0)
	if size != last {
		return nil, fmt.Errorf("failed serialization size=%d, lastIdx=%d", size, last)
	}
	return s.d, err
}

func stackTrace() string {
	msg := ""
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		msg += fmt.Sprintln(depth, ": ", file, ":", line)
	}
	return msg
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

	case reflect.Float32:
		ret += s.calcFloat32(0)

	case reflect.Float64:
		ret += s.calcFloat64(0)

	case reflect.String:
		ret += s.calcString(rv.String())

	case reflect.Bool:
		// do nothing

	case reflect.Slice:
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

	case reflect.Array:
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
		// todo : func
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
		ret = size

	case reflect.Interface:
		size, err := s.calcSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret = size

	case reflect.Invalid:
		// do nothing (return nil)

		// todo : default
	}

	return ret, nil
}

func (s *serializer) create(rv reflect.Value, offset int) int {

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		offset = s.writeUint(v, offset)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		offset = s.writeInt(v, offset)

	case reflect.Float32:
		offset = s.writeFloat32(rv.Float(), offset)

	case reflect.Float64:
		offset = s.writeFloat64(rv.Float(), offset)

	case reflect.Bool:
		offset = s.writeBool(rv.Bool(), offset)

	case reflect.String:
		offset = s.writeString(rv.String(), offset)

	case reflect.Slice:
		if rv.IsNil() {
			return s.writeNil(offset)
		}
		l := rv.Len()
		// bin format
		if s.isByteSlice(rv) {
			offset = s.writeByteSliceLength(l, offset)
			offset = s.setBytes(rv.Bytes(), offset)
			return offset
		}

		// format
		offset = s.writeSliceLength(l, offset)

		if offset, find := s.writeFixedSlice(rv, offset); find {
			return offset
		}

		// objects
		for i := 0; i < l; i++ {
			offset = s.create(rv.Index(i), offset)
		}

	case reflect.Array:
		l := rv.Len()
		// bin format
		if s.isByteSlice(rv) {
			offset = s.writeByteSliceLength(l, offset)
			// objects
			for i := 0; i < l; i++ {
				offset = s.setByte1Uint64(rv.Index(i).Uint(), offset)
			}
			return offset
		}

		// format
		offset = s.writeSliceLength(l, offset)

		// objects
		for i := 0; i < l; i++ {
			offset = s.create(rv.Index(i), offset)
		}

	case reflect.Map:
		if rv.IsNil() {
			return s.writeNil(offset)
		}

		l := rv.Len()
		offset = s.writeMapLength(l, offset)

		if offset, find := s.writeFixedMap(rv, offset); find {
			return offset
		}

		// key-value
		keys := rv.MapKeys()
		for _, k := range keys {
			offset = s.create(k, offset)
			offset = s.create(rv.MapIndex(k), offset)
		}

	case reflect.Struct:
		if isTime, tm := s.isDateTime(rv); isTime {
			return s.writeTime(tm, offset)
		}
		if s.asArray {
			offset = s.writeStructArray(rv, offset)
		} else {
			offset = s.writeStructMap(rv, offset)
		}

	case reflect.Ptr:
		if rv.IsNil() {
			return s.writeNil(offset)
		}

		offset = s.create(rv.Elem(), offset)

	case reflect.Interface:
		offset = s.create(rv.Elem(), offset)

	case reflect.Invalid:
		return s.writeNil(offset)

	}
	return offset
}
