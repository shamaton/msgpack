package msgpack

import (
	"github.com/shamaton/msgpack/ext"
	"github.com/shamaton/msgpack/internal/decoding"
	"github.com/shamaton/msgpack/internal/encoding"
)

var StructAsArray = false

func Encode(v interface{}) ([]byte, error) {
	return encoding.Encode(v, StructAsArray)
}

func Decode(data []byte, v interface{}) error {
	return decoding.Decode(data, v, StructAsArray)
}

func AddExtCoder(e ext.Encoder, d ext.Decoder) {
	encoding.AddExtEncoder(e)
	decoding.AddExtDecoder(d)
}

func RemoveExtCoder(e ext.Encoder, d ext.Decoder) {
	encoding.RemoveExtEncoder(e)
	decoding.RemoveExtDecoder(d)
}
