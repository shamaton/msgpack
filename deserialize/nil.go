package deserialize

import "github.com/shamaton/msgpack/def"

func (d *deserializer) isCodeNil(v byte) bool {
	return def.Nil == v
}
