package decoding

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/shamaton/msgpack/v2/internal/common"
)

type decoder struct {
	asArray bool
	common.Common
}

// Decode analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Decode(input io.Reader, v interface{}, asArray bool) error {
	d := decoder{asArray: asArray}

	if input == nil {
		return fmt.Errorf("data is nil")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", v)
	}

	rv = rv.Elem()

	// if input is already a bufio.Reader, just type assert it
	bufReader, ok := input.(*bufio.Reader)
	if !ok {
		// otherwise, wrap the input in a bufio reader
		bufReader = bufio.NewReader(input)
	}

	return d.decode(rv, bufReader)
}

// DecodeBytes analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func DecodeBytes(data []byte, v interface{}, asArray bool) error {
	d := decoder{asArray: asArray}

	if data == nil {
		return fmt.Errorf("data is nil")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", v)
	}

	rv = rv.Elem()

	input := bytes.NewReader(data)
	err := d.decode(rv, bufio.NewReader(input))
	if err != nil {
		return err
	}
	if input.Len() != 0 {
		return fmt.Errorf("failed deserialization size=%d, last=%d", len(data), len(data)-input.Len())
	}
	return err
}

func skipOne(reader *bufio.Reader) error {
	_, err := reader.ReadByte()
	return err
}

func skipN(reader *bufio.Reader, n int) error {
	_, err := reader.Discard(n)
	return err
}

func peekCode(reader *bufio.Reader) (byte, error) {
	code, err := reader.Peek(1)
	if err != nil {
		return 0, err
	}

	return code[0], nil
}

func (d *decoder) decode(rv reflect.Value, reader *bufio.Reader) error {
	k := rv.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := d.asInt(reader, k)
		if err != nil {
			return err
		}
		rv.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := d.asUint(reader, k)
		if err != nil {
			return err
		}
		rv.SetUint(v)

	case reflect.Float32:
		v, err := d.asFloat32(reader, k)
		if err != nil {
			return err
		}
		rv.SetFloat(float64(v))

	case reflect.Float64:
		v, err := d.asFloat64(reader, k)
		if err != nil {
			return err
		}
		rv.SetFloat(v)

	case reflect.String:
		code, err := reader.ReadByte()
		if err != nil {
			return err
		}

		// byte slice
		if d.isCodeBin(code) {
			v, err := d.asBinStringC(reader, code, k)
			if err != nil {
				return err
			}
			rv.SetString(v)
			return nil
		}
		v, err := d.asStringByteC(reader, code, nil, k)
		if err != nil {
			return err
		}
		rv.SetString(string(v))

	case reflect.Bool:
		v, err := d.asBool(reader, k)
		if err != nil {
			return err
		}
		rv.SetBool(v)

	case reflect.Slice:
		code, err := reader.ReadByte()
		if err != nil {
			return err
		}

		// nil
		if d.isCodeNil(code) {
			return nil
		}

		// byte slice
		if d.isCodeBin(code) {
			bs, err := d.asBinC(reader, code, k)
			if err != nil {
				return err
			}
			rv.SetBytes(bs)
			return nil
		}

		// string to bytes
		if d.isCodeString(code) {
			bs, err := d.asStringByteC(reader, code, nil, k)
			if err != nil {
				return err
			}
			rv.SetBytes(bs)
			return nil
		}

		// get slice length
		l, err := d.sliceLengthC(reader, code, k)
		if err != nil {
			return err
		}

		// check fixed type
		found, err := d.asFixedSlice(rv, reader, l)
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
				err = d.setStruct(v, reader, k)
			} else {
				err = d.decode(v, reader)
			}
			if err != nil {
				return err
			}
		}
		rv.Set(tmpSlice)

	case reflect.Complex64:
		v, err := d.asComplex64(reader, k)
		if err != nil {
			return err
		}
		rv.SetComplex(complex128(v))

	case reflect.Complex128:
		v, err := d.asComplex128(reader, k)
		if err != nil {
			return err
		}
		rv.SetComplex(v)

	case reflect.Array:
		code, err := reader.ReadByte()
		if err != nil {
			return err
		}

		// nil
		if d.isCodeNil(code) {
			return nil
		}

		// byte slice
		if d.isCodeBin(code) {
			bs, err := d.asBinC(reader, code, k)
			if err != nil {
				return err
			}
			if len(bs) > rv.Len() {
				return errors.New(rv.Type().String() + " len is " +
					strconv.FormatInt(int64(rv.Len()), 10) + ", but msgpack has "+
					strconv.FormatInt(int64(len(bs)), 10) +" elements")
			}
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return nil
		}
		// string to bytes
		if d.isCodeString(code) {
			bs, err := d.asStringByteC(reader, code, nil, k)
			if err != nil {
				return err
			}
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return nil
		}

		// get slice length
		l, err := d.sliceLengthC(reader, code, k)
		if err != nil {
			return err
		}

		if l > rv.Len() {
			return errors.New(rv.Type().String() + " len is " +
				strconv.FormatInt(int64(rv.Len()), 10) + ", but msgpack has "+
				strconv.FormatInt(int64(l), 10) +" elements")
		}

		// create array dynamically
		for i := 0; i < l; i++ {
			err = d.decode(rv.Index(i), reader)
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		code, err := reader.ReadByte()
		if err != nil {
			return err
		}

		// nil
		if d.isCodeNil(code) {
			return nil
		}

		// get map length
		l, err := d.mapLengthC(reader, code, k)
		if err != nil {
			return err
		}

		// check fixed type
		found, err := d.asFixedMap(rv, reader, l)
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
			err = d.decode(k, reader)
			if err != nil {
				return err
			}
			err = d.decode(v, reader)
			if err != nil {
				return err
			}

			rv.SetMapIndex(k, v)
		}

	case reflect.Struct:
		err := d.setStruct(rv, reader, k)
		if err != nil {
			return err
		}

	case reflect.Ptr:
		code, err := peekCode(reader)
		if err != nil {
			return err
		}

		// nil
		if d.isCodeNil(code) {
			_, err = reader.ReadByte()
			return err
		}

		if rv.Elem().Kind() == reflect.Invalid {
			n := reflect.New(rv.Type().Elem())
			rv.Set(n)
		}

		err = d.decode(rv.Elem(), reader)
		if err != nil {
			return err
		}

	case reflect.Interface:
		if rv.Elem().Kind() == reflect.Ptr {
			err := d.decode(rv.Elem(), reader)
			if err != nil {
				return err
			}
		} else {
			v, err := d.asInterface(reader, k)
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

func (d *decoder) errorTemplate(code byte, k reflect.Kind) error {
	return errors.New("msgpack : invalid code " + strconv.FormatInt(int64(code), 16) + " decoding " + k.String())
}
