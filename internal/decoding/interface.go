package decoding

import (
	"bufio"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) asInterface(reader *bufio.Reader, k reflect.Kind) (interface{}, error) {
	code, err := peekCode(reader)
	if err != nil {
		return nil, err
	}

	switch {
	case code == def.Nil:
		return nil, skipOne(reader)

	case code == def.True, code == def.False:
		v, err := d.asBool(reader, k)
		if err != nil {
			return nil, err
		}
		return v, nil

	case d.isPositiveFixNum(code), code == def.Uint8:
		v, err := d.asUint(reader, k)
		if err != nil {
			return nil, err
		}
		return uint8(v), err
	case code == def.Uint16:
		v, err := d.asUint(reader, k)
		if err != nil {
			return nil, err
		}
		return uint16(v), err
	case code == def.Uint32:
		v, err := d.asUint(reader, k)
		if err != nil {
			return nil, err
		}
		return uint32(v), err
	case code == def.Uint64:
		v, err := d.asUint(reader, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case d.isNegativeFixNum(code), code == def.Int8:
		v, err := d.asInt(reader, k)
		if err != nil {
			return nil, err
		}
		return int8(v), err
	case code == def.Int16:
		v, err := d.asInt(reader, k)
		if err != nil {
			return nil, err
		}
		return int16(v), err
	case code == def.Int32:
		v, err := d.asInt(reader, k)
		if err != nil {
			return nil, err
		}
		return int32(v), err
	case code == def.Int64:
		v, err := d.asInt(reader, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case code == def.Float32:
		v, err := d.asFloat32(reader, k)
		if err != nil {
			return nil, err
		}
		return v, err
	case code == def.Float64:
		v, err := d.asFloat64(reader, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case d.isFixString(code), code == def.Str8, code == def.Str16, code == def.Str32:
		v, err := d.asString(reader, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case code == def.Bin8, code == def.Bin16, code == def.Bin32:
		v, err := d.asBin(reader, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case d.isFixSlice(code), code == def.Array16, code == def.Array32:
		l, err := d.sliceLength(reader, k)
		if err != nil {
			return nil, err
		}

		v := make([]interface{}, l)
		for i := 0; i < l; i++ {
			vv, err := d.asInterface(reader, k)
			if err != nil {
				return nil, err
			}
			v[i] = vv
		}
		return v, nil

	case d.isFixMap(code), code == def.Map16, code == def.Map32:
		l, err := d.mapLength(reader, k)
		if err != nil {
			return nil, err
		}
		v := make(map[interface{}]interface{}, l)
		for i := 0; i < l; i++ {
			key, err := d.asInterface(reader, k)
			if err != nil {
				return nil, err
			}
			value, err := d.asInterface(reader, k)
			if err != nil {
				return nil, err
			}
			v[key] = value
		}
		return v, nil
	}

	/* use ext
	if d.isDateTime(offset) {
		v, err := d.asDateTime(reader, k)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	*/

	if code, data, err := d.readExt(reader); err == nil {
		for i := range extCoders {
			if extCoders[i].Code() == int8(code) {
				v, err := extCoders[i].AsValue(data, k)
				if err != nil {
					return nil, err
				}
				return v, nil
			}
		}
	}

	return nil, d.errorTemplate(code, k)
}
