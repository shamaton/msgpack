package serialize

import (
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcFixedSlice(rv reflect.Value) (int, bool) {
	size := 0
	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			size += def.Byte1 + s.calcInt(int64(v))
		}
		return size, true
	}
	return size, false
}

func (s *serializer) writeSliceLength(l int, offset int) int {
	// format size
	if l <= 0x0f {
		offset = s.setByte1Int(def.FixArray+l, offset)
	} else if l <= math.MaxUint16 {
		offset = s.setByte1Int(def.Array16, offset)
		offset = s.setByte2Int(l, offset)
	} else if l <= math.MaxUint32 {
		offset = s.setByte1Int(def.Array32, offset)
		offset = s.setByte4Int(l, offset)
	}
	return offset
}

func (s *serializer) writeFixedSlice(rv reflect.Value, offset int) (int, bool) {
	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			offset = s.writeInt(int64(v), offset)
		}
		return offset, true
	case []uint:
	case []int8:
	case []int16:
	case []int32:
	case []int64:
	case []uint8:
	case []uint16:
	case []uint32:
	case []uint64:
	case []string:
	}

	return offset, false
}
