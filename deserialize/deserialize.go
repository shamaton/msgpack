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
	k := rv.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, o, err := d.asInt(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetInt(v)
		offset = o

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, o, err := d.asUint(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetUint(v)
		offset = o

	case reflect.Float32:
		v, o, err := d.asFloat32(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetFloat(float64(v))
		offset = o

	case reflect.Float64:
		v, o, err := d.asFloat64(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetFloat(v)
		offset = o

	}
	return offset, nil
}

func (d *deserializer) isPositiveFixNum(v byte) bool {
	if def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax {
		return true
	}
	return false
}

func (d *deserializer) isNegativeFixNum(v byte) bool {
	if def.NegativeFixintMin <= int8(v) && int8(v) <= def.NegativeFixintMax {
		return true
	}
	return false
}
