package encoding

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/common"
)

type encoder struct {
	d           []byte
	asArray     bool
	extRegistry *extEncoderRegistry
	common.Common
	mk map[uintptr][]reflect.Value
	mv map[uintptr][]reflect.Value
}

// Encode returns the MessagePack-encoded byte array of v.
func Encode(v interface{}, asArray bool) (b []byte, err error) {
	e := encoder{
		asArray:     asArray,
		extRegistry: currentExtEncoderRegistry.Load(),
	}
	/*
		defer func() {
			e := recover()
			if e != nil {
				b = nil
				err = fmt.Errorf("unexpected error!! \n%s", stackTrace())
			}
		}()
	*/

	rv := reflect.ValueOf(v)
	noCustom := e.extRegistry.customCount == 0
	var size int
	if noCustom {
		size, err = e.calcSizeBuiltIn(rv)
	} else {
		size, err = e.calcSize(rv)
	}
	if err != nil {
		return nil, err
	}

	e.d = make([]byte, size)
	var last int
	if noCustom {
		last = e.createBuiltIn(rv, 0)
	} else {
		last = e.create(rv, 0)
	}
	if size != last {
		return nil, fmt.Errorf("%w size=%d, lastIdx=%d", def.ErrNotMatchLastIndex, size, last)
	}
	return e.d, err
}

//func stackTrace() string {
//	msg := ""
//	for depth := 0; ; depth++ {
//		_, file, line, ok := runtime.Caller(depth)
//		if !ok {
//			break
//		}
//		msg += fmt.Sprintln(depth, ": ", file, ":", line)
//	}
//	return msg
//}

func (e *encoder) calcSize(rv reflect.Value) (int, error) {
	kind := rv.Kind()
	if kind == reflect.Invalid {
		return def.Byte1, nil
	}
	if encoder, ok := e.extEncoderForValue(kind, rv); ok {
		return encoder.CalcByteSize(rv)
	}
	return e.calcSizeByKind(rv, kind)
}

func (e *encoder) calcSizeDefault(rv reflect.Value) (int, error) {
	return e.calcSizeByKind(rv, rv.Kind())
}

func (e *encoder) calcSizeByKind(rv reflect.Value, kind reflect.Kind) (int, error) {
	switch kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		return e.calcUint(v), nil

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		return e.calcInt(int64(v)), nil

	case reflect.Float32:
		return e.calcFloat32(0), nil

	case reflect.Float64:
		return e.calcFloat64(0), nil

	case reflect.String:
		return e.calcString(rv.String()), nil

	case reflect.Bool:
		return def.Byte1, nil

	case reflect.Complex64:
		return e.calcComplex64(), nil

	case reflect.Complex128:
		return e.calcComplex128(), nil

	case reflect.Slice:
		if rv.IsNil() {
			return def.Byte1, nil
		}
		var elementType reflect.Type
		var elementEncoder ext.Encoder
		var elementHasExt bool
		if e.extRegistry.customCount != 0 {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}
		// bin format
		if !elementHasExt && e.isByteSlice(rv) {
			size, err := e.calcByteSlice(rv.Len())
			if err != nil {
				return 0, err
			}
			return size, nil
		}

		if !elementHasExt {
			if size, find := e.calcFixedSlice(rv); find {
				return size, nil
			}
		}
		if elementType == nil {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}

		l := rv.Len()
		size, err := e.calcLength(l)
		if err != nil {
			return 0, err
		}

		if elementHasExt {
			for i := 0; i < l; i++ {
				s, err := elementEncoder.CalcByteSize(rv.Index(i))
				if err != nil {
					return 0, err
				}
				size += s
			}
			return size, nil
		}

		calcElement := e.calcSizeDefault
		if e.extRegistry.customCount == 0 {
			calcElement = e.calcSizeBuiltIn
		} else if elementType.Kind() == reflect.Struct {
			calcElement = e.getStructCalc(elementType)
		}
		for i := 0; i < l; i++ {
			s, err := calcElement(rv.Index(i))
			if err != nil {
				return 0, err
			}
			size += s
		}
		return size, nil

	case reflect.Array:
		var elementType reflect.Type
		var elementEncoder ext.Encoder
		var elementHasExt bool
		if e.extRegistry.customCount != 0 {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}
		// bin format
		if !elementHasExt && e.isByteSlice(rv) {
			size, err := e.calcByteSlice(rv.Len())
			if err != nil {
				return 0, err
			}
			return size, nil
		}
		if elementType == nil {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}

		l := rv.Len()
		size, err := e.calcLength(l)
		if err != nil {
			return 0, err
		}

		if elementHasExt {
			for i := 0; i < l; i++ {
				s, err := elementEncoder.CalcByteSize(rv.Index(i))
				if err != nil {
					return 0, err
				}
				size += s
			}
			return size, nil
		}

		calcElement := e.calcSizeDefault
		if e.extRegistry.customCount == 0 {
			calcElement = e.calcSizeBuiltIn
		} else if elementType.Kind() == reflect.Struct {
			calcElement = e.getStructCalc(elementType)
		}
		for i := 0; i < l; i++ {
			s, err := calcElement(rv.Index(i))
			if err != nil {
				return 0, err
			}
			size += s
		}
		return size, nil

	case reflect.Map:
		if rv.IsNil() {
			return def.Byte1, nil
		}

		var keyEncoder, valueEncoder ext.Encoder
		var keyHasExt, valueHasExt bool
		if e.extRegistry.customCount != 0 {
			keyEncoder, keyHasExt = e.extEncoderForType(rv.Type().Key())
			valueEncoder, valueHasExt = e.extEncoderForType(rv.Type().Elem())
		}
		if !keyHasExt && !valueHasExt {
			if size, find := e.calcFixedMap(rv); find {
				return size, nil
			}
		}
		if e.extRegistry.customCount == 0 {
			keyEncoder, keyHasExt = e.extEncoderForType(rv.Type().Key())
			valueEncoder, valueHasExt = e.extEncoderForType(rv.Type().Elem())
		}

		if e.mk == nil {
			e.mk = map[uintptr][]reflect.Value{}
			e.mv = map[uintptr][]reflect.Value{}
		}

		keys := rv.MapKeys()
		size, err := e.calcLength(len(keys))
		if err != nil {
			return 0, err
		}

		// key-value
		mv := make([]reflect.Value, len(keys))
		i := 0
		for _, k := range keys {
			var keySize int
			if keyHasExt {
				keySize, err = keyEncoder.CalcByteSize(k)
			} else if e.extRegistry.customCount == 0 {
				keySize, err = e.calcSizeBuiltIn(k)
			} else {
				keySize, err = e.calcSizeDefault(k)
			}
			if err != nil {
				return 0, err
			}
			value := rv.MapIndex(k)
			var valueSize int
			if valueHasExt {
				valueSize, err = valueEncoder.CalcByteSize(value)
			} else if e.extRegistry.customCount == 0 {
				valueSize, err = e.calcSizeBuiltIn(value)
			} else {
				valueSize, err = e.calcSizeDefault(value)
			}
			if err != nil {
				return 0, err
			}
			size += keySize + valueSize
			mv[i] = value
			i++
		}
		e.mk[rv.Pointer()], e.mv[rv.Pointer()] = keys, mv
		return size, nil

	case reflect.Struct:
		size, err := e.calcStruct(rv)
		if err != nil {
			return 0, err
		}
		return size, nil

	case reflect.Ptr:
		if rv.IsNil() {
			return def.Byte1, nil
		}
		calcElement := e.calcSize
		if e.extRegistry.customCount == 0 {
			calcElement = e.calcSizeBuiltIn
		}
		size, err := calcElement(rv.Elem())
		if err != nil {
			return 0, err
		}
		return size, nil

	case reflect.Interface:
		calcElement := e.calcSize
		if e.extRegistry.customCount == 0 {
			calcElement = e.calcSizeBuiltIn
		}
		size, err := calcElement(rv.Elem())
		if err != nil {
			return 0, err
		}
		return size, nil

	default:
		return 0, fmt.Errorf("%v is %w type", kind, def.ErrUnsupportedType)
	}
}

func (e *encoder) calcLength(l int) (int, error) {
	if l <= 0x0f {
		return def.Byte1, nil
	} else if l <= math.MaxUint16 {
		return def.Byte1 + def.Byte2, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte1 + def.Byte4, nil
	}
	// not supported error
	return 0, fmt.Errorf("array length %d is %w", l, def.ErrUnsupportedLength)
}

func (e *encoder) create(rv reflect.Value, offset int) int {
	kind := rv.Kind()
	if kind == reflect.Invalid {
		return e.writeNil(offset)
	}
	if encoder, ok := e.extEncoderForValue(kind, rv); ok {
		return e.writeExt(encoder, rv, offset)
	}
	return e.createByKind(rv, kind, offset)
}

func (e *encoder) createDefault(rv reflect.Value, offset int) int {
	return e.createByKind(rv, rv.Kind(), offset)
}

func (e *encoder) writeExt(encoder ext.Encoder, rv reflect.Value, offset int) int {
	data := e.d
	offset = encoder.WriteToBytes(rv, offset, &data)
	e.d = data
	return offset
}

func (e *encoder) createByKind(rv reflect.Value, kind reflect.Kind, offset int) int {
	switch kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		offset = e.writeUint(v, offset)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		offset = e.writeInt(v, offset)

	case reflect.Float32:
		offset = e.writeFloat32(rv.Float(), offset)

	case reflect.Float64:
		offset = e.writeFloat64(rv.Float(), offset)

	case reflect.Bool:
		offset = e.writeBool(rv.Bool(), offset)

	case reflect.String:
		offset = e.writeString(rv.String(), offset)

	case reflect.Complex64:
		offset = e.writeComplex64(complex64(rv.Complex()), offset)

	case reflect.Complex128:
		offset = e.writeComplex128(rv.Complex(), offset)

	case reflect.Slice:
		if rv.IsNil() {
			return e.writeNil(offset)
		}
		var elementType reflect.Type
		var elementEncoder ext.Encoder
		var elementHasExt bool
		if e.extRegistry.customCount != 0 {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}

		// bin format
		if !elementHasExt && e.isByteSlice(rv) {
			offset = e.writeByteSliceLength(rv.Len(), offset)
			offset = e.setBytes(rv.Bytes(), offset)
			return offset
		}

		if !elementHasExt {
			if offset, find := e.writeFixedSlice(rv, offset); find {
				return offset
			}
		}
		if elementType == nil {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}

		// objects
		l := rv.Len()
		offset = e.writeSliceLength(l, offset)
		if elementHasExt {
			data := e.d
			for i := 0; i < l; i++ {
				offset = elementEncoder.WriteToBytes(rv.Index(i), offset, &data)
			}
			e.d = data
			return offset
		}
		writeElement := e.createDefault
		if e.extRegistry.customCount == 0 {
			writeElement = e.createBuiltIn
		} else if elementType.Kind() == reflect.Struct {
			writeElement = e.getStructWriter(elementType)
		}
		for i := 0; i < l; i++ {
			offset = writeElement(rv.Index(i), offset)
		}

	case reflect.Array:
		l := rv.Len()
		var elementType reflect.Type
		var elementEncoder ext.Encoder
		var elementHasExt bool
		if e.extRegistry.customCount != 0 {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}
		// bin format
		if !elementHasExt && e.isByteSlice(rv) {
			offset = e.writeByteSliceLength(l, offset)
			// objects
			for i := 0; i < l; i++ {
				offset = e.setByte1Uint64(rv.Index(i).Uint(), offset)
			}
			return offset
		}
		if elementType == nil {
			elementType = rv.Type().Elem()
			elementEncoder, elementHasExt = e.extEncoderForType(elementType)
		}

		// format
		offset = e.writeSliceLength(l, offset)

		if elementHasExt {
			data := e.d
			for i := 0; i < l; i++ {
				offset = elementEncoder.WriteToBytes(rv.Index(i), offset, &data)
			}
			e.d = data
			return offset
		}
		writeElement := e.createDefault
		if e.extRegistry.customCount == 0 {
			writeElement = e.createBuiltIn
		} else if elementType.Kind() == reflect.Struct {
			writeElement = e.getStructWriter(elementType)
		}
		for i := 0; i < l; i++ {
			offset = writeElement(rv.Index(i), offset)
		}

	case reflect.Map:
		if rv.IsNil() {
			return e.writeNil(offset)
		}

		l := rv.Len()
		offset = e.writeMapLength(l, offset)

		var keyEncoder, valueEncoder ext.Encoder
		var keyHasExt, valueHasExt bool
		if e.extRegistry.customCount != 0 {
			keyEncoder, keyHasExt = e.extEncoderForType(rv.Type().Key())
			valueEncoder, valueHasExt = e.extEncoderForType(rv.Type().Elem())
		}
		if !keyHasExt && !valueHasExt {
			if offset, find := e.writeFixedMap(rv, offset); find {
				return offset
			}
		}
		if e.extRegistry.customCount == 0 {
			keyEncoder, keyHasExt = e.extEncoderForType(rv.Type().Key())
			valueEncoder, valueHasExt = e.extEncoderForType(rv.Type().Elem())
		}

		// key-value
		p := rv.Pointer()
		data := e.d
		for i := range e.mk[p] {
			if keyHasExt {
				offset = keyEncoder.WriteToBytes(e.mk[p][i], offset, &data)
				e.d = data
			} else if e.extRegistry.customCount == 0 {
				offset = e.createBuiltIn(e.mk[p][i], offset)
				data = e.d
			} else {
				offset = e.createDefault(e.mk[p][i], offset)
				data = e.d
			}
			if valueHasExt {
				offset = valueEncoder.WriteToBytes(e.mv[p][i], offset, &data)
				e.d = data
			} else if e.extRegistry.customCount == 0 {
				offset = e.createBuiltIn(e.mv[p][i], offset)
				data = e.d
			} else {
				offset = e.createDefault(e.mv[p][i], offset)
				data = e.d
			}
		}
		e.d = data

	case reflect.Struct:
		offset = e.writeStruct(rv, offset)

	case reflect.Ptr:
		if rv.IsNil() {
			return e.writeNil(offset)
		}

		if e.extRegistry.customCount == 0 {
			offset = e.createBuiltIn(rv.Elem(), offset)
		} else {
			offset = e.create(rv.Elem(), offset)
		}

	case reflect.Interface:
		if e.extRegistry.customCount == 0 {
			offset = e.createBuiltIn(rv.Elem(), offset)
		} else {
			offset = e.create(rv.Elem(), offset)
		}

	}
	return offset
}
