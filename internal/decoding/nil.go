package decoding

import "github.com/shamaton/msgpack/v3/def"

func (d *decoder) isCodeNil(v byte) bool {
	return def.Nil == v
}
