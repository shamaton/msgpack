package msgpack

import (
	"github.com/shamaton/msgpack/internal/encoding"
)

// EncodeStructAsMap encodes data as map format.
// This is the same thing that StructAsArray sets false.
func EncodeStructAsMap(v interface{}) ([]byte, error) {
	return encoding.Encode(v, false)
}

// EncodeStructAsArray encodes data as array format.
// This is the same thing that StructAsArray sets true.
func EncodeStructAsArray(v interface{}) ([]byte, error) {
	return encoding.Encode(v, true)
}
