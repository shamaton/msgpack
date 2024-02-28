package decoding

import (
	"fmt"
	"io"
	"reflect"
)

// Decode analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Decode(r io.Reader, v interface{}, asArray bool) error {

	if r == nil {
		return fmt.Errorf("reader is nil")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", v)
	}

	rv = rv.Elem()

	err := decode(r, rv, asArray)
	if err != nil {
		return err
	}
	// todo : maybe enable to delete
	//if len(data) != last {
	//	return fmt.Errorf("failed deserialization size=%d, last=%d", len(data), last)
	//}
	return err
}

func decode(r io.Reader, rv reflect.Value, asArray bool) error {
	code, err := readSize1(r)
	if err != nil {
		return err
	}
	return decodeWithCode(r, code, rv, asArray)
}

func decodeWithCode(r io.Reader, code byte, rv reflect.Value, asArray bool) error {
	k := rv.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			return err
		}
		rv.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			return err
		}
		rv.SetUint(v)

	case reflect.Float32:
		v, err := asFloat32WithCode(r, code, k)
		if err != nil {
			return err
		}
		rv.SetFloat(float64(v))

	case reflect.Float64:
		v, err := asFloat64WithCode(r, code, k)
		if err != nil {
			return err
		}
		rv.SetFloat(v)

	case reflect.String:
		// byte slice
		if isCodeBin(code) {
			v, err := asBinStringWithCode(r, code, k)
			if err != nil {
				return err
			}
			rv.SetString(v)
			return nil
		}
		v, err := asStringWithCode(r, code, k)
		if err != nil {
			return err
		}
		rv.SetString(v)

	case reflect.Bool:
		v, err := asBoolWithCode(r, code, k)
		if err != nil {
			return err
		}
		rv.SetBool(v)

	case reflect.Slice:
		// nil
		if isCodeNil(code) {
			return nil
		}
		// byte slice
		if isCodeBin(code) {
			bs, err := asBinWithCode(r, code, k)
			if err != nil {
				return err
			}
			rv.SetBytes(bs)
			return nil
		}
		// string to bytes
		if isCodeString(code) {
			l, err := stringByteLength(r, code, k)
			if err != nil {
				return err
			}
			bs, err := asStringByteByLength(r, l, k)
			if err != nil {
				return err
			}
			rv.SetBytes(bs)
			return nil
		}

		// get slice length
		l, err := sliceLength(r, code, k)
		if err != nil {
			return err
		}

		//if err = d.hasRequiredLeastSliceSize(o, l); err != nil {
		//	return err
		//}

		// check fixed type
		found, err := asFixedSlice(r, rv, l)
		if err != nil {
			return err
		}
		if found {
			return nil
		}

		// create slice dynamically
		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)
		for i := 0; i < l; i++ {
			v := tmpSlice.Index(i)
			if v.Kind() == reflect.Struct {
				structCode, err := readSize1(r)
				if err != nil {
					return err
				}
				err = setStruct(r, structCode, v, k, asArray)
			} else {
				err = decode(r, v, asArray)
			}
			if err != nil {
				return err
			}
		}
		rv.Set(tmpSlice)

	case reflect.Complex64:
		v, err := asComplex64(r, code, k)
		if err != nil {
			return err
		}
		rv.SetComplex(complex128(v))

	case reflect.Complex128:
		v, err := asComplex128(r, code, k)
		if err != nil {
			return err
		}
		rv.SetComplex(v)

	case reflect.Array:
		// nil
		if isCodeNil(code) {
			return nil
		}
		// byte slice
		if isCodeBin(code) {
			bs, err := asBinWithCode(r, code, k)
			if err != nil {
				return err
			}
			if len(bs) > rv.Len() {
				return fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), len(bs))
			}
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return nil
		}
		// string to bytes
		if isCodeString(code) {
			l, err := stringByteLength(r, code, k)
			if err != nil {
				return err
			}
			if l > rv.Len() {
				return fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), l)
			}
			bs, err := asStringByteByLength(r, l, k)
			if err != nil {
				return err
			}
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return nil
		}

		// get slice length
		l, err := sliceLength(r, code, k)
		if err != nil {
			return err
		}

		if l > rv.Len() {
			return fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), l)
		}

		// todo : maybe enable to delete
		//if err = d.hasRequiredLeastSliceSize(o, l); err != nil {
		//	return err
		//}

		// create array dynamically
		for i := 0; i < l; i++ {
			err = decode(r, rv.Index(i), asArray)
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		// nil
		if isCodeNil(code) {
			return nil
		}

		// get map length
		l, err := mapLength(r, code, k)
		if err != nil {
			return err
		}

		// todo : maybe enable to delete
		//if err =   hasRequiredLeastMapSize(o, l); err != nil {
		//	return err
		//}

		// check fixed type
		found, err := asFixedMap(r, rv, l)
		if err != nil {
			return err
		}
		if found {
			return nil
		}

		// create dynamically
		key := rv.Type().Key()
		value := rv.Type().Elem()
		if rv.IsNil() {
			rv.Set(reflect.MakeMapWithSize(rv.Type(), l))
		}
		for i := 0; i < l; i++ {
			k := reflect.New(key).Elem()
			v := reflect.New(value).Elem()
			err = decode(r, k, asArray)
			if err != nil {
				return err
			}
			err = decode(r, v, asArray)
			if err != nil {
				return err
			}

			rv.SetMapIndex(k, v)
		}

	case reflect.Struct:
		err := setStruct(r, code, rv, k, asArray)
		if err != nil {
			return err
		}

	case reflect.Ptr:
		// nil
		if isCodeNil(code) {
			return nil
		}

		if rv.Elem().Kind() == reflect.Invalid {
			n := reflect.New(rv.Type().Elem())
			rv.Set(n)
		}

		err := decodeWithCode(r, code, rv.Elem(), asArray)
		if err != nil {
			return err
		}

	case reflect.Interface:
		if rv.Elem().Kind() == reflect.Ptr {
			err := decode(r, rv.Elem(), asArray)
			if err != nil {
				return err
			}
		} else {
			v, err := asInterfaceWithCode(r, code, k)
			if err != nil {
				return err
			}
			if v != nil {
				rv.Set(reflect.ValueOf(v))
			}
		}

	default:
		return fmt.Errorf("type(%v) is unsupported", rv.Kind())
	}
	return nil
}

func errorTemplate(code byte, k reflect.Kind) error {
	return fmt.Errorf("msgpack : invalid code %x decoding %v", code, k)
}
