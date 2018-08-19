package msgpack

import (
	"github.com/shamaton/msgpack/serialize"
)

var defaultSerializer = SerializeAsMap

func Serialize(v interface{}) ([]byte, error) {
	return defaultSerializer(v)
}

func SerializeAsArray(v interface{}) ([]byte, error) {
	return serialize.Exec(v, true)
}

func SerializeAsMap(v interface{}) ([]byte, error) {
	return serialize.Exec(v, false)
}
