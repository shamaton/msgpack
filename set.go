package msgpack

func SetDefaultToArray() {
	defaultSerializer = SerializeAsArray
	defaultDeserializer = DeserializeAsArray
}

func SetDefaultToMap() {
	defaultSerializer = SerializeAsMap
	defaultDeserializer = DeserializeAsMap
}
