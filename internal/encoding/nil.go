package encoding

import (
	"github.com/shamaton/msgpack/v2/def"
)

func (e *encoder) writeNil(writer Writer) error {
	return e.setByte1Int(def.Nil, writer)
}
