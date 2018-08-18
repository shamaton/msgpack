package serialize

import (
	"math"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcFloat32(v float64) int {
	return def.Byte4
}

func (s *serializer) calcFloat64(v float64) int {
	return def.Byte8
}

func (s *serializer) writeFloat32(v float64, offset int) int {
	offset = s.setByte1Int(def.Float32, offset)
	offset = s.setByte4Uint64(uint64(math.Float32bits(float32(v))), offset)
	return offset
}

func (s *serializer) writeFloat64(v float64, offset int) int {
	offset = s.setByte1Int(def.Float64, offset)
	offset = s.setByte8Uint64(math.Float64bits(v), offset)
	return offset
}
