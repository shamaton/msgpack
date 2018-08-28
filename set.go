package msgpack

import (
	"github.com/shamaton/msgpack/ext"
	"github.com/shamaton/msgpack/internal/deserialize"
	"github.com/shamaton/msgpack/internal/serialize"
)

var asArray = false

func SetDefaultToArray() {
	asArray = true
}

func SetDefaultToMap() {
	asArray = false
}

func SetExtFunc(e ext.Encoder, d ext.Decoder) {
	serialize.SetExtFunc(e)
	deserialize.SetExtFunc(d)
}

func UnsetExtFunc(e ext.Encoder, d ext.Decoder) {
	serialize.UnsetExtFunc(e)
	deserialize.UnsetExtFunc(d)
}
