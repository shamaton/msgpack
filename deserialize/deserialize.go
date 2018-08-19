package deserialize

import (
	"fmt"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

type deserializer struct {
	data    []byte
	asArray bool
}

func Exec(data []byte, holder interface{}, asArray bool) error {
	d := deserializer{data: data, asArray: asArray}

	rv := reflect.ValueOf(holder)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", holder)
	}

	rv = rv.Elem()
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	_, err := d.deserialize(rv, 0)
	return err
}

func (d *deserializer) deserialize(rv reflect.Value, offset int) (int, error) {

	// TODO : offset use uint
	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if d.isFixNum(offset) {
			b, o := d.readSize1(offset)
			rv.SetUint(uint64(b))
			offset = o
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
	}
	return 0, nil
}

func (d *deserializer) isFixNum(offset int) bool {
	if def.PositiveFixIntMin <= d.data[offset] && d.data[offset] <= def.PositiveFixIntMax {
		return true
	} else if def.NegativeFixintMin <= int8(d.data[offset]) && int8(d.data[offset]) <= def.NegativeFixintMax {
		return true
	}
	return false
}
