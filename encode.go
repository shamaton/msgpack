package msgpack

import (
	"github.com/shamaton/msgpack/internal/encoding"
)

// MarshalAsMap encodes data as map format.
// This is the same thing that StructAsArray sets false.
func MarshalAsMap(v interface{}) ([]byte, error) {
	return encoding.Encode(v, false)
}

// EncodeStructAsArray encodes data as array format.
// This is the same thing that StructAsArray sets true.
func MarshalAsArray(v interface{}) ([]byte, error) {
	return encoding.Encode(v, true)
}

// Deprecated: Use MarshalAsMap, this method will be deleted.
func EncodeStructAsMap(v interface{}) ([]byte, error) {
	return MarshalAsMap(v)
}

// Deprecated: Use MarshalAsArray, this method will be deleted.
func EncodeStructAsArray(v interface{}) ([]byte, error) {
	return MarshalAsArray(v)
}
