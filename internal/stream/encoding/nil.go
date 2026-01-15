package encoding

import "github.com/shamaton/msgpack/v3/def"

func (e *encoder) writeNil() error {
	return e.setByte1Int(def.Nil)
}
