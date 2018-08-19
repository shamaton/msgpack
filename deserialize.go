package msgpack

import (
	"github.com/shamaton/msgpack/deserialize"
)

var defaultDeserializer = DeserializeAsMap

func Deserialize(data []byte, v interface{}) error {
	return defaultDeserializer(data, v)
}

func DeserializeAsArray(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, true)
}

func DeserializeAsMap(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, false)
}
