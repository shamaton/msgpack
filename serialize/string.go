package serialize

import (
	"math"
	"unsafe"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcString(v string) int {
	// NOTE : unsafe
	strBytes := *(*[]byte)(unsafe.Pointer(&v))
	l := len(strBytes)
	if l < 32 {
		return l
	} else if l <= math.MaxUint8 {
		return def.Byte1 + l
	} else if l <= math.MaxUint16 {
		return def.Byte2 + l
	}
	return def.Byte4 + l
	// NOTE : length over uint32
}

func (s *serializer) writeString(str string, offset int) int {
	// NOTE : unsafe
	strBytes := *(*[]byte)(unsafe.Pointer(&str))
	l := len(strBytes)
	if l < 32 {
		offset = s.setByte1Int(def.FixStr+l, offset)
		offset = s.setBytes(strBytes, offset)
	} else if l <= math.MaxUint8 {
		offset = s.setByte1Int(def.Str8, offset)
		offset = s.setByte1Int(l, offset)
		offset = s.setBytes(strBytes, offset)
	} else if l <= math.MaxUint16 {
		offset = s.setByte1Int(def.Str16, offset)
		offset = s.setByte2Int(l, offset)
		offset = s.setBytes(strBytes, offset)
	} else {
		offset = s.setByte1Int(def.Str32, offset)
		offset = s.setByte4Int(l, offset)
		offset = s.setBytes(strBytes, offset)
	}
	return offset
}
