package msgpack

import "errors"

var defaultDeserializer = DeserializeAsMap

func Deserialize(data []byte, v interface{}) error {
	return defaultDeserializer(data, v)
}

func DeserializeAsArray(data []byte, v interface{}) error {
	return errors.New("not implement")
}

func DeserializeAsMap(data []byte, v interface{}) error {
	return errors.New("not implement")
}
