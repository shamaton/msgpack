package encoding

import (
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (e *encoder) calcFixedMap(rv reflect.Value) (int, bool) {
	size := 0

	// todo : add types
	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			size += def.Byte1 + e.calcString(k)
			size += def.Byte1 + e.calcInt(int64(v))
		}
		return size, true

	case map[string]uint:
		for k, v := range m {
			size += def.Byte1 + e.calcString(k)
			size += def.Byte1 + e.calcUint(uint64(v))
		}
		return size, true

	case map[string]float32:
		for k := range m {
			size += def.Byte1 + e.calcString(k)
			size += def.Byte1 + e.calcFloat32(0)
		}
		return size, true

	case map[string]float64:
		for k := range m {
			size += def.Byte1 + e.calcString(k)
			size += def.Byte1 + e.calcFloat64(0)
		}
		return size, true

	case map[string]bool:
		for k := range m {
			size += def.Byte1 + e.calcString(k)
			size += def.Byte1 /*+ e.calcBool()*/
		}
		return size, true

	case map[string]string:
		for k, v := range m {
			size += def.Byte1 + e.calcString(k)
			size += def.Byte1 + e.calcString(v)
		}
		return size, true

	case map[int]int:
		for k, v := range m {
			size += def.Byte1 + e.calcInt(int64(k))
			size += def.Byte1 + e.calcInt(int64(v))
		}
		return size, true

	case map[int]uint:
		for k, v := range m {
			size += def.Byte1 + e.calcInt(int64(k))
			size += def.Byte1 + e.calcUint(uint64(v))
		}
		return size, true

	case map[int]string:
		for k, v := range m {
			size += def.Byte1 + e.calcInt(int64(k))
			size += def.Byte1 + e.calcString(v)
		}
		return size, true
	}
	return size, false
}

func (e *encoder) writeMapLength(l int, offset int) int {

	// format
	if l <= 0x0f {
		offset = e.setByte1Int(def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(l, offset)
	} else if l <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(l, offset)
	}
	return offset
}

func (e *encoder) writeFixedMap(rv reflect.Value, offset int) (int, bool) {
	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true

	case map[string]uint:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true

	case map[string]float32:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeFloat32(float64(v), offset)
		}
		return offset, true

	case map[string]float64:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeFloat64(v, offset)
		}
		return offset, true

	case map[string]bool:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeBool(v, offset)
		}
		return offset, true

	case map[string]string:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeString(v, offset)
		}
		return offset, true

	case map[int]int:
		for k, v := range m {
			offset = e.writeInt(int64(k), offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true

	case map[int]uint:
		for k, v := range m {
			offset = e.writeInt(int64(k), offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true

	case map[int]string:
		for k, v := range m {
			offset = e.writeInt(int64(k), offset)
			offset = e.writeString(v, offset)
		}
		return offset, true
	}
	return offset, false
}
