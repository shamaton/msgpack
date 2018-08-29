package decoding

import "github.com/shamaton/msgpack/def"

func (d *decoder) isCodeNil(v byte) bool {
	return def.Nil == v
}
