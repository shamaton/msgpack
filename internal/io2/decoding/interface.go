package decoding

import (
	"fmt"
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func asInterface(r io.Reader, k reflect.Kind) (interface{}, error) {
	code, err := readSize1(r)
	if err != nil {
		return 0, err
	}
	return asInterfaceWithCode(r, code, k)
}

func asInterfaceWithCode(r io.Reader, code byte, k reflect.Kind) (interface{}, error) {
	switch {
	case code == def.Nil:
		return nil, nil

	case code == def.True, code == def.False:
		v, err := asBoolWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, nil

	case isPositiveFixNum(code), code == def.Uint8:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return uint8(v), err
	case code == def.Uint16:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return uint16(v), err
	case code == def.Uint32:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return uint32(v), err
	case code == def.Uint64:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case isNegativeFixNum(code), code == def.Int8:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return int8(v), err
	case code == def.Int16:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return int16(v), err
	case code == def.Int32:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return int32(v), err
	case code == def.Int64:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case code == def.Float32:
		v, err := asFloat32WithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, err
	case code == def.Float64:
		v, err := asFloat64WithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case isFixString(code), code == def.Str8, code == def.Str16, code == def.Str32:
		v, err := asStringWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case code == def.Bin8, code == def.Bin16, code == def.Bin32:
		v, err := asBinWithCode(r, code, k)
		if err != nil {
			return nil, err
		}
		return v, err

	case isFixSlice(code), code == def.Array16, code == def.Array32:
		l, err := sliceLength(r, code, k)
		if err != nil {
			return nil, err
		}

		// todo : maybe enable to delete
		//if err =   hasRequiredLeastSliceSize(o, l); err != nil {
		//	return nil, err
		//}

		v := make([]interface{}, l)
		for i := 0; i < l; i++ {
			vv, err := asInterface(r, k)
			if err != nil {
				return nil, err
			}
			v[i] = vv
		}
		return v, nil

	case isFixMap(code), code == def.Map16, code == def.Map32:
		l, err := mapLength(r, code, k)
		if err != nil {
			return nil, err
		}
		// todo : maybe enable to delete
		//if err =   hasRequiredLeastMapSize(o, l); err != nil {
		//	return nil, err
		//}
		v := make(map[interface{}]interface{}, l)
		for i := 0; i < l; i++ {
			keyCode, err := readSize1(r)
			if err != nil {
				return 0, err
			}

			if canSetAsMapKey(keyCode) != nil {
				return nil, err
			}
			key, err := asInterfaceWithCode(r, keyCode, k)
			if err != nil {
				return nil, err
			}
			value, err := asInterface(r, k)
			if err != nil {
				return nil, err
			}
			v[key] = value
		}
		return v, nil
	}

	// ext
	extInnerType, extData, err := readIfExtType(r, code)
	if err != nil {
		return nil, err
	}
	for i := range extCoders {
		if extCoders[i].IsType(code, extInnerType, len(extData)) {
			v, err := extCoders[i].AsValue(code, extData, k)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
	}
	return nil, errorTemplate(code, k)
}

func canSetAsMapKey(code byte) error {
	switch {
	case isFixSlice(code), code == def.Array16, code == def.Array32:
		return fmt.Errorf("can not use slice code for map key/ code: %x", code)
	case isFixMap(code), code == def.Map16, code == def.Map32:
		return fmt.Errorf("can not use map code for map key/ code: %x", code)
	}
	return nil
}
