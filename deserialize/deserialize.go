package deserialize

import (
	"encoding/binary"
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
func (d *deserializer) _deserialize(rv reflect.Value, offset int) (int, error) {

	// TODO : offset use uint
	switch rv.Kind() {
	case reflect.Uint:
		if def.IsIntSize32 {
			v, o, err := d.asUint32(rv, offset)
			if err != nil {
				return 0, err
			}
			rv.SetUint(uint64(v))
			offset = o
		} else {
			v, o, err := d.asUint64(rv, offset)
			if err != nil {
				return 0, err
			}
			rv.SetUint(v)
			offset = o
		}

	case reflect.Uint8:
		v, o, err := d.asUint8(rv, offset)
		if err != nil {
			return 0, err
		}
		rv.SetUint(uint64(v))
		offset = o

	case reflect.Uint16:
		v, o, err := d.asUint16(rv, offset)
		if err != nil {
			return 0, err
		}
		rv.SetUint(uint64(v))
		offset = o

	case reflect.Uint32:
		v, o, err := d.asUint32(rv, offset)
		if err != nil {
			return 0, err
		}
		rv.SetUint(uint64(v))
		offset = o

	case reflect.Uint64:
		v, o, err := d.asUint64(rv, offset)
		if err != nil {
			return 0, err
		}
		rv.SetUint(v)
		offset = o

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
	}
	return offset, nil
}

func (d *deserializer) readAsUint(rv reflect.Value, offset int) (int, error) {
	code := d.data[offset]

	a := byte(0xff)
	b := int8(a)
	c := uint32(b)
	fmt.Println("c : ", c)

	if d.isFixNum(code) {
		b, o := d.readSize1(offset)
		rv.SetUint(uint64(b))
		offset = o
	} else if code == def.Uint8 {
		offset++
		b, o := d.readSize1(offset)
		rv.SetUint(uint64(b))
		offset = o
	} else if code == def.Uint16 {
		offset++
		b, o := d.readSize2(offset)
		rv.SetUint(uint64(binary.BigEndian.Uint16(b)))
		offset = o
	} else if code == def.Uint32 {
		offset++
		b, o := d.readSize4(offset)
		rv.SetUint(uint64(binary.BigEndian.Uint32(b)))
		offset = o
	} else if code == def.Uint64 {
		offset++
		b, o := d.readSize8(offset)
		rv.SetUint(binary.BigEndian.Uint64(b))
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

func (d *deserializer) isFixNum(v byte) bool {
	if def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax {
		return true
	} else if def.NegativeFixintMin <= int8(v) && int8(v) <= def.NegativeFixintMax {
		return true
	}
	return false
}
