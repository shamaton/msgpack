package msgpack

import (
	"github.com/shamaton/msgpack/internal/deserialize"
)

func Deserialize(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, asArray)
}

func DeserializeStructAsArray(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, true)
}

func DeserializeStructAsMap(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, false)
}
