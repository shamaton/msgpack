package decoding

import "github.com/shamaton/msgpack/v2/def"

func isCodeNil(v byte) bool {
	return def.Nil == v
}
