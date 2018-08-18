package serialize

import "github.com/shamaton/msgpack/def"

func (s *serializer) calcBool() int {
	return 0
}

func (s *serializer) writeBool(v bool, offset int) int {
	if v {
		offset = s.setByte1Int(def.True, offset)
	} else {
		offset = s.setByte1Int(def.False, offset)
	}
	return offset
}
