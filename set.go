package msgpack

import (
	"github.com/shamaton/msgpack/deserialize"
	"github.com/shamaton/msgpack/ext"
	"github.com/shamaton/msgpack/serialize"
)

var asArray = false

func SetDefaultToArray() {
	asArray = true
}

func SetDefaultToMap() {
	asArray = false
}

func SetExtFunc(f ext.ExtSeri, f2 ext.ExtDeseri) {
	serialize.SetExtFunc(f)
	deserialize.SetExtFunc(f2)
}

func UnsetExtFunc(f ext.ExtSeri, f2 ext.ExtDeseri) {
	serialize.UnsetExtFunc(f)
	deserialize.UnsetExtFunc(f2)
}
