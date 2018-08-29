package msgpack

import (
	"github.com/shamaton/msgpack/internal/encoding"
)

func EncodeStructAsMap(v interface{}) ([]byte, error) {
	return encoding.Encode(v, false)
}

func EncodeStructAsArray(v interface{}) ([]byte, error) {
	return encoding.Encode(v, true)
}
