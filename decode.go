package msgpack

import (
	"github.com/shamaton/msgpack/internal/deserialize"
)

func DecodeStructAsMap(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, false)
}

func DecodeStructAsArray(data []byte, v interface{}) error {
	return deserialize.Exec(data, v, true)
}
