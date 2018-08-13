package serialize

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/shamaton/msgpack/def"
)

func AsArray2(v interface{}) ([]byte, error) {
	s := serializer{}

	size, err := s.calcSize2(v)
	if err != nil {
		return nil, err
	}

	s.d = make([]byte, size)
	last, err := s.create2(v, 0)
	if err != nil {
		return nil, err
	}
	if size != last {
		return nil, fmt.Errorf("failed serialization size=%d, lastIdx=%d", size, last)
	}
	return s.d, err
}

func (s *serializer) calcSize2(v interface{}) (int, error) {
	ret := def.Byte1

	switch v := v.(type) {
	case int:
		ret += s.calcInt(int64(v))
	case uint:
		ret += s.calcUint(uint64(v))

	case string:
		ret += s.calcString(v)

	case bool, nil:
		// do nothing

	case []int:
		lSize, err := s.calcSliceLength(len(v))
		if err != nil {
			return 0, err
		}
		ret += lSize

		// objects size
		for _, vv := range v {
			size, err := s.calcSize2(vv)
			if err != nil {
				return 0, err
			}
			ret += size
		}

	case []uint:
		lSize, err := s.calcSliceLength(len(v))
		if err != nil {
			return 0, err
		}
		ret += lSize

		// objects size
		for _, vv := range v {
			size, err := s.calcSize2(vv)
			if err != nil {
				return 0, err
			}
			ret += size
		}

	case int8:
		ret += s.calcInt(int64(v))
	case int16:
		ret += s.calcInt(int64(v))
	case int32:
		ret += s.calcInt(int64(v))
	case int64:
		ret += s.calcInt(v)

	case uint8:
		ret += s.calcUint(uint64(v))
	case uint16:
		ret += s.calcUint(uint64(v))
	case uint32:
		ret += s.calcUint(uint64(v))
	case uint64:
		ret += s.calcUint(v)

	}

	// analyze value

	return ret, nil
}

func (s *serializer) create2(v interface{}, offset int) (int, error) {

	switch v := v.(type) {
	case int:
		offset = s.writeInt(int64(v), offset)
	case uint:
		offset = s.writeUint(uint64(v), offset)

	case string:
		// NOTE : unsafe
		strBytes := *(*[]byte)(unsafe.Pointer(&v))
		l := len(strBytes)
		if l < 32 {
			offset = s.writeSize1Int(def.FixStr+l, offset)
			offset = s.writeBytes(strBytes, offset)
		} else if l <= math.MaxUint8 {
			offset = s.writeSize1Int(def.Str8, offset)
			offset = s.writeSize1Int(l, offset)
			offset = s.writeBytes(strBytes, offset)
		} else if l <= math.MaxUint16 {
			offset = s.writeSize1Int(def.Str16, offset)
			offset = s.writeSize2Int(l, offset)
			offset = s.writeBytes(strBytes, offset)
		} else {
			offset = s.writeSize1Int(def.Str32, offset)
			offset = s.writeSize4Int(l, offset)
			offset = s.writeBytes(strBytes, offset)
		}

	case bool:
		if v {
			offset = s.writeSize1Int(def.True, offset)
		} else {
			offset = s.writeSize1Int(def.False, offset)
		}

	case []int:
		offset = s.writeSliceLength(len(v), offset)

		// objects
		for _, vv := range v {
			o, err := s.create2(vv, offset)
			if err != nil {
				return 0, err
			}
			offset = o
		}

	case []uint:

	case nil:
		offset = s.writeSize1Int(def.Nil, offset)

	case int8:
		offset = s.writeInt(int64(v), offset)
	case int16:
		offset = s.writeInt(int64(v), offset)
	case int32:
		offset = s.writeInt(int64(v), offset)
	case int64:
		offset = s.writeInt(v, offset)

	case uint8:
		offset = s.writeUint(uint64(v), offset)
	case uint16:
		offset = s.writeUint(uint64(v), offset)
	case uint32:
		offset = s.writeUint(uint64(v), offset)
	case uint64:
		offset = s.writeUint(uint64(v), offset)
	}
	return offset, nil
}
