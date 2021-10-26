package encoding

import (
	"io"

	"github.com/shamaton/msgpack/v2/def"
)

func (e *encoder) writeNil(writer io.Writer) error {
	return e.setByte1Int(def.Nil, writer)
}
