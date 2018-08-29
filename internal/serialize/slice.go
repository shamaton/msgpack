package serialize

import (
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcFixedSlice(rv reflect.Value) (int, bool) {
	size := 0
	// todo : add types
	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			size += def.Byte1 + s.calcInt(int64(v))
		}
		return size, true

	case []uint:
		for _, v := range sli {
			size += def.Byte1 + s.calcUint(uint64(v))
		}
		return size, true

	case []string:
		for _, v := range sli {
			size += def.Byte1 + s.calcString(v)
		}
		return size, true

	case []float32:
		for _, v := range sli {
			size += def.Byte1 + s.calcFloat32(float64(v))
		}
		return size, true

	case []float64:
		for _, v := range sli {
			size += def.Byte1 + s.calcFloat64(v)
		}
		return size, true

	case []bool:
		size += def.Byte1 * len(sli)
		return size, true

	case []int8:
	case []int16:
	case []int32:
	case []int64:
	case []uint8:
	case []uint16:
	case []uint32:
	case []uint64:
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
	// todo : add types
	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			offset = s.writeInt(int64(v), offset)
		}
		return offset, true

	case []uint:
		for _, v := range sli {
			offset = s.writeUint(uint64(v), offset)
		}
		return offset, true

	case []string:
		for _, v := range sli {
			offset = s.writeString(v, offset)
		}
		return offset, true

	case []float32:
		for _, v := range sli {
			offset = s.writeFloat32(float64(v), offset)
		}
		return offset, true

	case []float64:
		for _, v := range sli {
			offset = s.writeFloat64(float64(v), offset)
		}
		return offset, true

	case []bool:
		for _, v := range sli {
			offset = s.writeBool(v, offset)
		}
		return offset, true

	case []int8:
	case []int16:
	case []int32:
	case []int64:
	case []uint8:
	case []uint16:
	case []uint32:
	case []uint64:
	}

	return offset, false
}
