package msgpack

import "github.com/shamaton/msgpack/internal/decoding"

func DecodeStructAsMap(data []byte, v interface{}) error {
	return decoding.Decode(data, v, false)
}

func DecodeStructAsArray(data []byte, v interface{}) error {
	return decoding.Decode(data, v, true)
}
