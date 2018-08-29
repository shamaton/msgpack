package msgpack

import (
	"github.com/shamaton/msgpack/ext"
	"github.com/shamaton/msgpack/internal/deserialize"
	"github.com/shamaton/msgpack/internal/serialize"
)

var StructAsArray = false

func Encode(v interface{}) ([]byte, error) {
	return serialize.Exec(v, StructAsArray)
}

func Decode(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, StructAsArray)
}

func AddExtCoder(e ext.Encoder, d ext.Decoder) {
	serialize.SetExtFunc(e)
	deserialize.SetExtFunc(d)
}

func RemoveExtCoder(e ext.Encoder, d ext.Decoder) {
	serialize.UnsetExtFunc(e)
	deserialize.UnsetExtFunc(d)
}
