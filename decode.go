package msgpack

import "github.com/shamaton/msgpack/internal/decoding"

// UnmarshalAsMap decodes data that is encoded as map format.
// This is the same thing that StructAsArray sets false.
func UnmarshalAsMap(data []byte, v interface{}) error {
	return decoding.Decode(data, v, false)
}

// UnmarshalAsArray decodes data that is encoded as array format.
// This is the same thing that StructAsArray sets true.
func UnmarshalAsArray(data []byte, v interface{}) error {
	return decoding.Decode(data, v, true)
}

// Deprecated: Use UnmarshalAsMap, this method will be deleted.
func DecodeStructAsMap(data []byte, v interface{}) error {
	return UnmarshalAsMap(data, v)
}

// Deprecated: Use UnmarshalAsArray, this method will be deleted.
func DecodeStructAsArray(data []byte, v interface{}) error {
	return UnmarshalAsArray(data, v)
}
