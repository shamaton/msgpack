package encoding

import (
	"fmt"
	"github.com/shamaton/msgpack/v2/def"
	"io"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/internal/common"
)

type encoder struct {
	w       io.Writer
	asArray bool
	buf     *common.Buffer
	common.Common
	mk map[uintptr][]reflect.Value
	mv map[uintptr][]reflect.Value
}

// Encode writes MessagePack-encoded byte array of v to writer.
func Encode(w io.Writer, v any, asArray bool) error {
	e := encoder{
		w:       w,
		buf:     common.GetBuffer(),
		asArray: asArray,
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}

	err := e.create(rv)
	if err == nil {
		err = e.buf.Flush(e.w)
	}
	common.PutBuffer(e.buf)
	return err
}

// Encode writes MessagePack-encoded byte array of v to writer.
func Encode2(v any, asArray bool) (data []byte, err error) {
	e := encoder{
		buf:     common.GetBuffer(),
		asArray: asArray,
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}

	size, err := e.calcSize(rv)
	if err != nil {
		return nil, err
	}
	e.buf.Grow(size)

	err = e.create(rv)
	if err == nil {
		data = make([]byte, size)
		copy(data, e.buf.Bytes())
	}
	common.PutBuffer(e.buf)
	return data, err
}

func (e *encoder) create(rv reflect.Value) error {

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		return e.writeUint(v)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		return e.writeInt(v)

	case reflect.Float32:
		return e.writeFloat32(rv.Float())

	case reflect.Float64:
		return e.writeFloat64(rv.Float())

	case reflect.Bool:
		return e.writeBool(rv.Bool())

	case reflect.String:
		return e.writeString(rv.String())

	case reflect.Complex64:
		return e.writeComplex64(complex64(rv.Complex()))

	case reflect.Complex128:
		return e.writeComplex128(rv.Complex())

	case reflect.Slice:
		if rv.IsNil() {
			return e.writeNil()
		}
		l := rv.Len()
		// bin format
		if e.isByteSlice(rv) {
			if err := e.writeByteSliceLength(l); err != nil {
				return err
			}
			return e.setBytes(rv.Bytes())
		}

		// format
		if err := e.writeSliceLength(l); err != nil {
			return err
		}

		if find, err := e.writeFixedSlice(rv); err != nil {
			return err
		} else if find {
			return nil
		}

		// func
		elem := rv.Type().Elem()
		var f structWriteFunc
		if elem.Kind() == reflect.Struct {
			f = e.getStructWriter(elem)
		} else {
			f = e.create
		}

		// objects
		for i := 0; i < l; i++ {
			if err := f(rv.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Array:
		l := rv.Len()
		// bin format
		if e.isByteSlice(rv) {
			if err := e.writeByteSliceLength(l); err != nil {
				return err
			}
			// objects
			for i := 0; i < l; i++ {
				if err := e.setByte1Uint64(rv.Index(i).Uint()); err != nil {
					return err
				}
			}
			return nil
		}

		// format
		if err := e.writeSliceLength(l); err != nil {
			return err
		}

		// func
		elem := rv.Type().Elem()
		var f structWriteFunc
		if elem.Kind() == reflect.Struct {
			f = e.getStructWriter(elem)
		} else {
			f = e.create
		}

		// objects
		for i := 0; i < l; i++ {
			if err := f(rv.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		if rv.IsNil() {
			return e.writeNil()
		}

		l := rv.Len()
		if err := e.writeMapLength(l); err != nil {
			return err
		}

		if find, err := e.writeFixedMap(rv); err != nil {
			return err
		} else if find {
			return nil
		}

		// key-value
		keys := rv.MapKeys()
		for _, k := range keys {
			if err := e.create(k); err != nil {
				return err
			}
			if err := e.create(rv.MapIndex(k)); err != nil {
				return err
			}
		}

	case reflect.Struct:
		return e.writeStruct(rv)

	case reflect.Ptr:
		if rv.IsNil() {
			return e.writeNil()
		}

		return e.create(rv.Elem())

	case reflect.Interface:
		return e.create(rv.Elem())

	case reflect.Invalid:
		return e.writeNil()
	default:
		return fmt.Errorf("type(%v) is unsupported", rv.Kind())
	}
	return nil
}

func (e *encoder) calcSize(rv reflect.Value) (int, error) {
	ret := def.Byte1

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		ret += e.calcUint(v)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		ret += e.calcInt(int64(v))

	case reflect.Float32:
		ret += e.calcFloat32(0)

	case reflect.Float64:
		ret += e.calcFloat64(0)

	case reflect.String:
		ret += e.calcString(rv.String())

	case reflect.Bool:
	// do nothing

	case reflect.Complex64:
		ret += e.calcComplex64()

	case reflect.Complex128:
		ret += e.calcComplex128()

	case reflect.Slice:
		if rv.IsNil() {
			return ret, nil
		}
		l := rv.Len()
		// bin format
		if e.isByteSlice(rv) {
			r, err := e.calcByteSlice(l)
			if err != nil {
				return 0, err
			}
			ret += r
			return ret, nil
		}

		// format size
		if l <= 0x0f {
			// format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}

		if size, find := e.calcFixedSlice(rv); find {
			ret += size
			return ret, nil
		}

		// func
		elem := rv.Type().Elem()
		var f structCalcFunc
		if elem.Kind() == reflect.Struct {
			f = e.getStructCalc(elem)
			ret += def.Byte1 * l
		} else {
			f = e.calcSize
		}

		// objects size
		for i := 0; i < l; i++ {
			size, err := f(rv.Index(i))
			if err != nil {
				return 0, err
			}
			ret += size
		}

	case reflect.Array:
		l := rv.Len()
		// bin format
		if e.isByteSlice(rv) {
			r, err := e.calcByteSlice(l)
			if err != nil {
				return 0, err
			}
			ret += r
			return ret, nil
		}

		// format size
		if l <= 0x0f {
			// format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}

		// func
		elem := rv.Type().Elem()
		var f structCalcFunc
		if elem.Kind() == reflect.Struct {
			f = e.getStructCalc(elem)
			ret += def.Byte1 * l
		} else {
			f = e.calcSize
		}

		// objects size
		for i := 0; i < l; i++ {
			size, err := f(rv.Index(i))
			if err != nil {
				return 0, err
			}
			ret += size
		}

	case reflect.Map:
		if rv.IsNil() {
			return ret, nil
		}

		l := rv.Len()
		// format
		if l <= 0x0f {
			// do nothing
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this map length : %d", l)
		}

		if size, find := e.calcFixedMap(rv); find {
			ret += size
			return ret, nil
		}

		if e.mk == nil {
			e.mk = map[uintptr][]reflect.Value{}
			e.mv = map[uintptr][]reflect.Value{}
		}

		// key-value
		keys := rv.MapKeys()
		mv := make([]reflect.Value, len(keys))
		i := 0
		for _, k := range keys {
			keySize, err := e.calcSize(k)
			if err != nil {
				return 0, err
			}
			value := rv.MapIndex(k)
			valueSize, err := e.calcSize(value)
			if err != nil {
				return 0, err
			}
			ret += keySize + valueSize
			mv[i] = value
			i++
		}
		e.mk[rv.Pointer()], e.mv[rv.Pointer()] = keys, mv

	case reflect.Struct:
		size, err := e.calcStruct(rv)
		if err != nil {
			return 0, err
		}
		ret += size

	case reflect.Ptr:
		if rv.IsNil() {
			return ret, nil
		}
		size, err := e.calcSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret = size

	case reflect.Interface:
		size, err := e.calcSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret = size

	case reflect.Invalid:
		// do nothing (return nil)

	default:
		return 0, fmt.Errorf("type(%v) is unsupported", rv.Kind())
	}

	return ret, nil
}
