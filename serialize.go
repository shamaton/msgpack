package msgpack

import (
	"github.com/shamaton/msgpack/serialize"
)

var defaultSerializer = SerializeAsMap

func Serialize(v interface{}) ([]byte, error) {
	return defaultSerializer(v)
}

func SerializeAsArray(v interface{}) ([]byte, error) {
	serialize.IsArray = true
	return serialize.AsArray(v)
}

func SerializeAsArray2(v interface{}) ([]byte, error) {
	// TODO : delete
	return serialize.AsArray2(v)
}

func SerializeAsMap(v interface{}) ([]byte, error) {
	serialize.IsArray = false
	return serialize.AsArray(v)
}
