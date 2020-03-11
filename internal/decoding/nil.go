package decoding

import "github.com/aucfan-yotsuya/msgpack/def"

func (d *decoder) isCodeNil(v byte) bool {
	return def.Nil == v
}
