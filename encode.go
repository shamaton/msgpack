package msgpack

import (
	"github.com/shamaton/msgpack/internal/serialize"
)

func EncodeStructAsMap(v interface{}) ([]byte, error) {
	return serialize.Exec(v, false)
}

func EncodeStructAsArray(v interface{}) ([]byte, error) {
	return serialize.Exec(v, true)
}
