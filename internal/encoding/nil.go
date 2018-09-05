package encoding

import "github.com/shamaton/msgpack/def"

func (e *encoder) writeNil(offset int) int {
	offset = e.setByte1Int(def.Nil, offset)
	return offset
}
