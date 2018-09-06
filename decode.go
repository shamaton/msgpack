package msgpack

import "github.com/shamaton/msgpack/internal/decoding"

// DecodeStructAsMap decodes data that is encoded as map format.
// This is the same thing that StructAsArray sets false.
func DecodeStructAsMap(data []byte, v interface{}) error {
	return decoding.Decode(data, v, false)
}

// DecodeStructAsArray decodes data that is encoded as array format.
// This is the same thing that StructAsArray sets true.
func DecodeStructAsArray(data []byte, v interface{}) error {
	return decoding.Decode(data, v, true)
}
