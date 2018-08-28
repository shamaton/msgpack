package serialize

import (
	"math"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) isPositiveFixInt64(v int64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (s *serializer) isNegativeFixInt64(v int64) bool {
	return def.NegativeFixintMin <= v && v <= def.NegativeFixintMax
}

func (s *serializer) calcInt(v int64) int {
	if v >= 0 {
		return s.calcUint(uint64(v))
	} else if s.isNegativeFixInt64(v) {
		// format code only
		return 0
	} else if v >= math.MinInt8 {
		return def.Byte1
	} else if v >= math.MinInt16 {
		return def.Byte2
	} else if v >= math.MinInt32 {
		return def.Byte4
	}
	return def.Byte8
}

func (s *serializer) writeInt(v int64, offset int) int {
	if v >= 0 {
		offset = s.writeUint(uint64(v), offset)
	} else if s.isNegativeFixInt64(v) {
		offset = s.setByte1Int64(v, offset)
	} else if v >= math.MinInt8 {
		offset = s.setByte1Int(def.Int8, offset)
		offset = s.setByte1Int64(v, offset)
	} else if v >= math.MinInt16 {
		offset = s.setByte1Int(def.Int16, offset)
		offset = s.setByte2Int64(v, offset)
	} else if v >= math.MinInt32 {
		offset = s.setByte1Int(def.Int32, offset)
		offset = s.setByte4Int64(v, offset)
	} else {
		offset = s.setByte1Int(def.Int64, offset)
		offset = s.setByte8Int64(v, offset)
	}
	return offset
}
