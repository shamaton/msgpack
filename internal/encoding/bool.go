package encoding

import (
	"github.com/shamaton/msgpack/v2/def"
)

//func (e *encoder) calcBool() int {
//	return 0
//}

func (e *encoder) writeBool(v bool, writer Writer) error {
	if v {
		return e.setByte1Int(def.True, writer)
	} else {
		return e.setByte1Int(def.False, writer)
	}
}
