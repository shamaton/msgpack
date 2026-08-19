package encoding

import (
	"fmt"
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/common"
)

type encoder struct {
	w           io.Writer
	asArray     bool
	buf         *common.Buffer
	extRegistry *extEncoderRegistry
	common.Common
}

// Encode writes MessagePack-encoded byte array of v to writer.
func Encode(w io.Writer, v any, asArray bool) error {
	e := encoder{
		w:           w,
		buf:         common.GetBuffer(),
		asArray:     asArray,
		extRegistry: currentExtEncoderRegistry.Load(),
	}

	rv := reflect.ValueOf(v)

	var err error
	if e.extRegistry.customCount == 0 {
		err = e.createBuiltIn(rv)
	} else {
		err = e.create(rv)
	}
	if err == nil {
		err = e.buf.Flush(e.w)
	}
	common.PutBuffer(e.buf)
	return err
}

func (e *encoder) create(rv reflect.Value) error {
	kind := rv.Kind()
	if kind == reflect.Invalid {
		return e.writeNil()
	}
	if encoder, ok := e.extEncoderForValue(kind, rv); ok {
		return e.writeExt(encoder, rv)
	}
	return e.createByKind(rv, kind)
}

func (e *encoder) createDefault(rv reflect.Value) error {
	return e.createByKind(rv, rv.Kind())
}

func (e *encoder) writeExt(encoder ext.StreamEncoder, rv reflect.Value) error {
	w := ext.CreateStreamWriter(e.w, e.buf)
	return encoder.Write(w, rv)
}

func (e *encoder) createByKind(rv reflect.Value, kind reflect.Kind) error {
	switch kind {
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
		elementType := rv.Type().Elem()
		elementKind := elementType.Kind()
		var elementEncoder ext.StreamEncoder
		var elementHasExt bool
		if e.extRegistry.customCount == 0 {
			if elementKind == reflect.Struct && elementType == builtInTimeType {
				elementEncoder, elementHasExt = builtInTimeEncoder, true
			}
		} else {
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}
		// bin format
		if !elementHasExt && e.isByteSlice(rv) {
			if err := e.writeByteSliceLength(l); err != nil {
				return err
			}
			return e.setBytes(rv.Bytes())
		}

		// format
		if err := e.writeSliceLength(l); err != nil {
			return err
		}

		if !elementHasExt {
			if find, err := e.writeFixedSlice(rv); err != nil {
				return err
			} else if find {
				return nil
			}
		}

		if elementHasExt {
			for i := 0; i < l; i++ {
				if err := e.writeExt(elementEncoder, rv.Index(i)); err != nil {
					return err
				}
			}
			return nil
		}
		writeElement := e.createDefault
		if e.extRegistry.customCount == 0 {
			writeElement = e.createBuiltIn
		} else if elementKind == reflect.Struct {
			writeElement = e.getStructWriter(elementType)
		}
		for i := 0; i < l; i++ {
			if err := writeElement(rv.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Array:
		l := rv.Len()
		elementType := rv.Type().Elem()
		elementKind := elementType.Kind()
		var elementEncoder ext.StreamEncoder
		var elementHasExt bool
		if e.extRegistry.customCount == 0 {
			if elementKind == reflect.Struct && elementType == builtInTimeType {
				elementEncoder, elementHasExt = builtInTimeEncoder, true
			}
		} else {
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}
		// bin format
		if !elementHasExt && e.isByteSlice(rv) {
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

		if elementHasExt {
			for i := 0; i < l; i++ {
				if err := e.writeExt(elementEncoder, rv.Index(i)); err != nil {
					return err
				}
			}
			return nil
		}
		writeElement := e.createDefault
		if e.extRegistry.customCount == 0 {
			writeElement = e.createBuiltIn
		} else if elementKind == reflect.Struct {
			writeElement = e.getStructWriter(elementType)
		}
		for i := 0; i < l; i++ {
			if err := writeElement(rv.Index(i)); err != nil {
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

		keyType, valueType := rv.Type().Key(), rv.Type().Elem()
		var keyEncoder, valueEncoder ext.StreamEncoder
		var keyHasExt, valueHasExt bool
		if e.extRegistry.customCount == 0 {
			if keyType.Kind() == reflect.Struct && keyType == builtInTimeType {
				keyEncoder, keyHasExt = builtInTimeEncoder, true
			}
			if valueType.Kind() == reflect.Struct && valueType == builtInTimeType {
				valueEncoder, valueHasExt = builtInTimeEncoder, true
			}
		} else {
			keyEncoder, keyHasExt = e.extEncoderForType(keyType)
			valueEncoder, valueHasExt = e.extEncoderForType(valueType)
		}
		if !keyHasExt && !valueHasExt {
			if find, err := e.writeFixedMap(rv); err != nil {
				return err
			} else if find {
				return nil
			}
		}

		// key-value
		keys := rv.MapKeys()
		for _, k := range keys {
			if keyHasExt {
				if err := e.writeExt(keyEncoder, k); err != nil {
					return err
				}
			} else if e.extRegistry.customCount == 0 {
				if err := e.createBuiltIn(k); err != nil {
					return err
				}
			} else if err := e.createDefault(k); err != nil {
				return err
			}
			value := rv.MapIndex(k)
			if valueHasExt {
				if err := e.writeExt(valueEncoder, value); err != nil {
					return err
				}
			} else if e.extRegistry.customCount == 0 {
				if err := e.createBuiltIn(value); err != nil {
					return err
				}
			} else if err := e.createDefault(value); err != nil {
				return err
			}
		}

	case reflect.Struct:
		return e.writeStruct(rv)

	case reflect.Ptr:
		if rv.IsNil() {
			return e.writeNil()
		}

		if e.extRegistry.customCount == 0 {
			return e.createBuiltIn(rv.Elem())
		}
		return e.create(rv.Elem())

	case reflect.Interface:
		if e.extRegistry.customCount == 0 {
			return e.createBuiltIn(rv.Elem())
		}
		return e.create(rv.Elem())

	default:
		return fmt.Errorf("%v is %w type", kind, def.ErrUnsupportedType)
	}
	return nil
}
