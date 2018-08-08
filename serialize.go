package msgpack

import (
	"errors"

	"github.com/shamaton/msgpack/serialize"
)

var defaultSerializer = SerializeAsMap

func Serialize(v interface{}) ([]byte, error) {
	return defaultSerializer(v)
}

func SerializeAsArray(v interface{}) ([]byte, error) {
	serialize.AsArray(v)
	return []byte{}, errors.New("not implement")
}

func SerializeAsMap(v interface{}) ([]byte, error) {
	return []byte{}, errors.New("not implement")
}
