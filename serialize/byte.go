package serialize

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

var typeByte = reflect.TypeOf(byte(0))

func (s *serializer) isByteSlice(rv reflect.Value) bool {
	return rv.Type().Elem() == typeByte
}

func (s *serializer) calcByteSlice(l int) (int, error) {
	if l <= math.MaxUint8 {
		return def.Byte1 + l, nil
	} else if l <= math.MaxUint16 {
		return def.Byte2 + l, nil
	} else if l <= math.MaxUint32 {
		return def.Byte4 + l, nil
	}
	// not supported error
	return 0, fmt.Errorf("not support this array length : %d", l)
}

func (s *serializer) writeByteSliceLength(l int, offset int) int {
	if l <= math.MaxUint8 {
		offset = s.setByte1Int(def.Bin8, offset)
		offset = s.setByte1Int(l, offset)
	} else if l <= math.MaxUint16 {
		offset = s.setByte1Int(def.Bin16, offset)
		offset = s.setByte2Int(l, offset)
	} else if l <= math.MaxUint32 {
		offset = s.setByte1Int(def.Bin32, offset)
		offset = s.setByte4Int(l, offset)
	}
	return offset
}
