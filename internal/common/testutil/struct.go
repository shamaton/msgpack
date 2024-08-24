package testutil

import (
	"math"
	"reflect"
	"strconv"

	"github.com/shamaton/msgpack/v2/def"
)

// CreateStruct returns a struct that is made dynamically and encoded bytes.
func CreateStruct(fieldNum int) (v any, asMapBytes []byte, asArrayBytes []byte) {
	if fieldNum < 0 {
		panic("negative field number")
	}

	fields := make([]reflect.StructField, 0, fieldNum)
	asMapBytes = make([]byte, 0, fieldNum*2)
	asArrayBytes = make([]byte, 0, fieldNum)

	for i := 0; i < fieldNum; i++ {
		// create struct field
		name := "A" + strconv.Itoa(i)
		typ := reflect.TypeOf(1)
		field := reflect.StructField{
			Name: name,
			Type: typ,
			Tag:  `json:"B"`,
		}
		fields = append(fields, field)

		// set encoded bytes
		if len(name) < 32 {
			asMapBytes = append(asMapBytes, def.FixStr+byte(len(name)))
		} else if len(name) < math.MaxUint8 {
			asMapBytes = append(asMapBytes, def.Str8)
			asMapBytes = append(asMapBytes, byte(len(name)))
		}
		for _, c := range name {
			asMapBytes = append(asMapBytes, byte(c))
		}
		value := byte(i % 0x7f)
		asMapBytes = append(asMapBytes, value)
		asArrayBytes = append(asArrayBytes, value)
	}
	t := reflect.StructOf(fields)

	// set field values
	v = reflect.New(t).Interface()
	rv := reflect.ValueOf(v)
	for i := 0; i < rv.Elem().NumField(); i++ {
		field := rv.Elem().Field(i)
		field.SetInt(int64(i % 0x7f))
	}
	return v, asMapBytes, asArrayBytes
}
