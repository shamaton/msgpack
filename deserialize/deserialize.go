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

	case reflect.Array, reflect.Slice:
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
			bs, offset, err := d.asStringByte(offset, k)
			if err != nil {
				return 0, err
			}
			rv.SetBytes(bs)
			return offset, nil
		}

		l, o, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		// allocate interface
		switch rv.Interface().(type) {
		case []int8:
			sli := make([]int8, l) // allocate
			k = rv.Type().Elem().Kind()
			for i := range sli {
				v, _o, err := d.asInt(o, k)
				if err != nil {
					return 0, nil
				}
				sli[i] = int8(v)
				o = _o
			}
			rv.Set(reflect.ValueOf(sli)) // allocate
			return o, nil
		}

		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)

		// element type
		e := rv.Type().Elem()
		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			tmpSlice.Index(i).Set(v)
		}
		rv.Set(tmpSlice)

	case reflect.Ptr:

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

func (d *deserializer) errorTemplate(code byte, k reflect.Kind) error {
	return fmt.Errorf("msgpack : invalid code %x decoding %v", code, k)
}
