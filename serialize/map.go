package serialize

import (
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcFixedMap(rv reflect.Value) (int, bool) {
	size := 0

	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			size += def.Byte1 + s.calcString(k)
			size += def.Byte1 + s.calcInt(int64(v))
		}
		return size, true

	case map[string]uint:
		for k, v := range m {
			size += def.Byte1 + s.calcString(k)
			size += def.Byte1 + s.calcUint(uint64(v))
		}
		return size, true

	case map[string]float32:
		for k := range m {
			size += def.Byte1 + s.calcString(k)
			size += def.Byte1 + s.calcFloat32(0)
		}
		return size, true

	case map[string]float64:
		for k := range m {
			size += def.Byte1 + s.calcString(k)
			size += def.Byte1 + s.calcFloat64(0)
		}
		return size, true

	case map[string]string:
		for k, v := range m {
			size += def.Byte1 + s.calcString(k)
			size += def.Byte1 + s.calcString(v)
		}
		return size, true

	case map[int]int:
		for k, v := range m {
			size += def.Byte1 + s.calcInt(int64(k))
			size += def.Byte1 + s.calcInt(int64(v))
		}
		return size, true

	case map[int]uint:
		for k, v := range m {
			size += def.Byte1 + s.calcInt(int64(k))
			size += def.Byte1 + s.calcUint(uint64(v))
		}
		return size, true

	case map[int]string:
		for k, v := range m {
			size += def.Byte1 + s.calcInt(int64(k))
			size += def.Byte1 + s.calcString(v)
		}
		return size, true
	}
	return size, false
}

func (s *serializer) writeMapLength(l int, offset int) int {

	// format
	if l <= 0x0f {
		offset = s.setByte1Int(def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = s.setByte1Int(def.Map16, offset)
		offset = s.setByte2Int(l, offset)
	} else if l <= math.MaxUint32 {
		offset = s.setByte1Int(def.Map32, offset)
		offset = s.setByte4Int(l, offset)
	}
	return offset
}

func (s *serializer) writeFixedMap(rv reflect.Value, offset int) (int, bool) {
	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			offset = s.writeString(k, offset)
			offset = s.writeInt(int64(v), offset)
		}
		return offset, true

	case map[string]uint:
		for k, v := range m {
			offset = s.writeString(k, offset)
			offset = s.writeUint(uint64(v), offset)
		}
		return offset, true

	case map[string]float32:
		for k, v := range m {
			offset = s.writeString(k, offset)
			offset = s.writeFloat32(float64(v), offset)
		}
		return offset, true

	case map[string]float64:
		for k, v := range m {
			offset = s.writeString(k, offset)
			offset = s.writeFloat64(v, offset)
		}
		return offset, true

	case map[string]string:
		for k, v := range m {
			offset = s.writeString(k, offset)
			offset = s.writeString(v, offset)
		}
		return offset, true

	case map[int]int:
		for k, v := range m {
			offset = s.writeInt(int64(k), offset)
			offset = s.writeInt(int64(v), offset)
		}
		return offset, true

	case map[int]uint:
		for k, v := range m {
			offset = s.writeInt(int64(k), offset)
			offset = s.writeUint(uint64(v), offset)
		}
		return offset, true

	case map[int]string:
		for k, v := range m {
			offset = s.writeInt(int64(k), offset)
			offset = s.writeString(v, offset)
		}
		return offset, true
	}
	return offset, false
}
