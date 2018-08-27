package msgpack

var asArray = false

func SetDefaultToArray() {
	asArray = true
}

func SetDefaultToMap() {
	asArray = false
}
