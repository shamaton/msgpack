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

	case reflect.String:
		v, o, err := d.asString(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetString(v)
		offset = o

	case reflect.Bool:
		v, o, err := d.asBool(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetBool(v)
		offset = o

	case reflect.Slice:
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}
		// byte slice
		if d.isCodeBin(d.data[offset]) {
			bs, offset, err := d.asBin(offset, k)
			if err != nil {
				return 0, err
			}
			rv.SetBytes(bs)
			return offset, nil
		}
		// string to bytes
		if d.isCodeString(d.data[offset]) {
			l, offset, err := d.stringByteLength(offset, k)
			if err != nil {
				return 0, err
			}
			bs, offset := d.asStringByte(offset, l, k)
			rv.SetBytes(bs)
			return offset, nil
		}

		// get slice length
		l, offset, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		// check fixed type
		fixedOffset, found, err := d.asFixedSlice(rv, offset, l)
		if err != nil {
			return 0, err
		}
		if found {
			return fixedOffset, nil
		}

		// create slice dynamically
		e := rv.Type().Elem()
		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)
		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			offset, err = d.deserialize(v, offset)
			if err != nil {
				return 0, err
			}

			tmpSlice.Index(i).Set(v)
		}
		rv.Set(tmpSlice)

	case reflect.Array:
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}
		// byte slice
		if d.isCodeBin(d.data[offset]) {
			// todo : length check
			bs, offset, err := d.asBin(offset, k)
			if err != nil {
				return 0, err
			}
			rv.SetBytes(bs)
			return offset, nil
		}
		// string to bytes
		if d.isCodeString(d.data[offset]) {
			l, offset, err := d.stringByteLength(offset, k)
			if err != nil {
				return 0, err
			}
			if l > rv.Len() {
				return 0, fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), l)
			}
			bs, offset := d.asStringByte(offset, l, k)
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return offset, nil
		}

		// get slice length
		l, offset, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		if l > rv.Len() {
			return 0, fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), l)
		}

		// create array dynamically
		for i := 0; i < l; i++ {
			offset, err = d.deserialize(rv.Index(i), offset)
			if err != nil {
				return 0, err
			}
		}

	case reflect.Map:
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}

		// get map length
		l, offset, err := d.mapLength(offset, k)
		if err != nil {
			return 0, err
		}

		// check fixed type
		fixedOffset, found, err := d.asFixedMap(rv, offset, l)
		if err != nil {
			return 0, err
		}
		if found {
			return fixedOffset, nil
		}

		// create dynamically
		key := rv.Type().Key()
		value := rv.Type().Elem()
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		for i := 0; i < l; i++ {
			k := reflect.New(key).Elem()
			v := reflect.New(value).Elem()
			o, err := d.deserialize(k, offset)
			if err != nil {
				return 0, err
			}
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			rv.SetMapIndex(k, v)
			offset = o
		}

	case reflect.Struct:
	case reflect.Ptr:

	case reflect.Interface:
		// all type...

	default:
		return 0, d.errorTemplate(d.data[offset], k)
	}
	return offset, nil
}

func (d *deserializer) asFixedMap(rv reflect.Value, offset int, l int) (int, bool, error) {
	t := rv.Type()

	keyKind := rv.Type().Key().Kind()
	valueKind := rv.Type().Elem().Kind()
	switch t {
	case typeMapStringInt:
		m := make(map[string]int, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = int(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil
	}

	return offset, false, nil
}

func (d *deserializer) asFixedSlice(rv reflect.Value, offset int, l int) (int, bool, error) {
	t := rv.Type()
	k := t.Elem().Kind()
	switch t {
	case typeIntSlice:
		sli := make([]int, l)
		for i := range sli {
			v, o, err := d.asInt(offset, k)
			if err != nil {
				return 0, false, err
			}
			sli[i] = int(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(sli))
		return offset, true, nil

	case typeInt8Slice:
		sli := make([]int8, l)
		for i := range sli {
			v, o, err := d.asInt(offset, k)
			if err != nil {
				return 0, false, err
			}
			sli[i] = int8(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(sli))
		return offset, true, nil
	}

	return offset, false, nil
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

func (d *deserializer) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *deserializer) isCodeBin(v byte) bool {
	switch v {
	case def.Bin8, def.Bin16, def.Bin32:
		return true
	}
	return false
}

func (d *deserializer) isCodeNil(v byte) bool {
	return def.Nil == v
}

func (d *deserializer) errorTemplate(code byte, k reflect.Kind) error {
	return fmt.Errorf("msgpack : invalid code %x decoding %v", code, k)
}
