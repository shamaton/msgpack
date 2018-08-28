package serialize

import (
	"math"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) isPositiveFixUint64(v uint64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (s *serializer) calcUint(v uint64) int {
	if v <= math.MaxInt8 {
		// format code only
		return 0
	} else if v <= math.MaxUint8 {
		return def.Byte1
	} else if v <= math.MaxUint16 {
		return def.Byte2
	} else if v <= math.MaxUint32 {
		return def.Byte4
	}
	return def.Byte8
}

func (s *serializer) writeUint(v uint64, offset int) int {
	if v <= math.MaxInt8 {
		offset = s.setByte1Uint64(v, offset)
	} else if v <= math.MaxUint8 {
		offset = s.setByte1Int(def.Uint8, offset)
		offset = s.setByte1Uint64(v, offset)
	} else if v <= math.MaxUint16 {
		offset = s.setByte1Int(def.Uint16, offset)
		offset = s.setByte2Uint64(v, offset)
	} else if v <= math.MaxUint32 {
		offset = s.setByte1Int(def.Uint32, offset)
		offset = s.setByte4Uint64(v, offset)
	} else {
		offset = s.setByte1Int(def.Uint64, offset)
		offset = s.setByte8Uint64(v, offset)
	}
	return offset
}
