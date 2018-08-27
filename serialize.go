package msgpack

import (
	"github.com/shamaton/msgpack/serialize"
)

func Serialize(v interface{}) ([]byte, error) {
	return serialize.Exec(v, asArray)
}

func SerializeStructAsArray(v interface{}) ([]byte, error) {
	return serialize.Exec(v, true)
}

func SerializeStructAsMap(v interface{}) ([]byte, error) {
	return serialize.Exec(v, false)
}
