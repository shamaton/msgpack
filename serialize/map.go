package serialize

import (
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcFixedMap(rv reflect.Value) (int, bool) {
	size := 0
	switch m := rv.Interface().(type) {
	case map[int]int:
		for k, v := range m {
			size += def.Byte1 + s.calcInt(int64(k))
			size += def.Byte1 + s.calcInt(int64(v))
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
	case map[int]int:
		for k, v := range m {
			offset = s.writeInt(int64(k), offset)
			offset = s.writeInt(int64(v), offset)
		}
		return offset, true
	}
	return offset, false
}
