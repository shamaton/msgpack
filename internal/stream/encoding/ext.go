package encoding

import (
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/time"
)

type extEncoderRegistry struct {
	byKind      [reflect.UnsafePointer + 1]map[reflect.Type]ext.StreamEncoder
	hasCustom   [reflect.UnsafePointer + 1]bool
	customCount int
}

var (
	extEncoderRegistryMu      sync.Mutex
	currentExtEncoderRegistry atomic.Pointer[extEncoderRegistry]
	builtInTimeEncoder        = time.StreamEncoder
	builtInTimeType           = builtInTimeEncoder.Type()
)

func init() {
	registry := &extEncoderRegistry{}
	registry.byKind[reflect.Struct] = map[reflect.Type]ext.StreamEncoder{
		builtInTimeType: builtInTimeEncoder,
	}
	currentExtEncoderRegistry.Store(registry)
}

// AddExtEncoder adds an encoder for an exact type without replacing an existing encoder.
func AddExtEncoder(encoder ext.StreamEncoder) {
	typ := encoder.Type()
	if typ == builtInTimeType {
		return
	}

	extEncoderRegistryMu.Lock()
	defer extEncoderRegistryMu.Unlock()

	current := currentExtEncoderRegistry.Load()
	kind := typ.Kind()
	if _, exists := current.byKind[kind][typ]; exists {
		return
	}

	next := *current
	next.byKind[kind] = cloneExtEncoderMap(current.byKind[kind], 1)
	next.byKind[kind][typ] = encoder
	next.hasCustom[kind] = true
	next.customCount++
	currentExtEncoderRegistry.Store(&next)
}

// RemoveExtEncoder removes the encoder registered for the encoder's exact type.
func RemoveExtEncoder(encoder ext.StreamEncoder) {
	typ := encoder.Type()
	if typ == builtInTimeType {
		return
	}

	extEncoderRegistryMu.Lock()
	defer extEncoderRegistryMu.Unlock()

	current := currentExtEncoderRegistry.Load()
	kind := typ.Kind()
	currentKind := current.byKind[kind]
	if _, exists := currentKind[typ]; !exists {
		return
	}

	next := *current
	if len(currentKind) == 1 {
		next.byKind[kind] = nil
		next.hasCustom[kind] = false
	} else {
		next.byKind[kind] = cloneExtEncoderMap(currentKind, 0)
		delete(next.byKind[kind], typ)
		next.hasCustom[kind] = kind != reflect.Struct || len(next.byKind[kind]) > 1
	}
	next.customCount--
	currentExtEncoderRegistry.Store(&next)
}

func cloneExtEncoderMap(source map[reflect.Type]ext.StreamEncoder, extraCapacity int) map[reflect.Type]ext.StreamEncoder {
	cloned := make(map[reflect.Type]ext.StreamEncoder, len(source)+extraCapacity)
	for typ, encoder := range source {
		cloned[typ] = encoder
	}
	return cloned
}

func (e *encoder) extEncoderForValue(kind reflect.Kind, value reflect.Value) (ext.StreamEncoder, bool) {
	if e.extRegistry.customCount == 0 {
		if kind == reflect.Struct && value.Type() == builtInTimeType {
			return builtInTimeEncoder, true
		}
		return nil, false
	}
	if kind == reflect.Struct && !e.extRegistry.hasCustom[reflect.Struct] {
		if value.Type() == builtInTimeType {
			return builtInTimeEncoder, true
		}
		return nil, false
	}
	coders := e.extRegistry.byKind[kind]
	if coders == nil {
		return nil, false
	}
	encoder, ok := coders[value.Type()]
	return encoder, ok
}

func (e *encoder) extEncoderForType(typ reflect.Type) (ext.StreamEncoder, bool) {
	kind := typ.Kind()
	if e.extRegistry.customCount == 0 {
		if kind == reflect.Struct && typ == builtInTimeType {
			return builtInTimeEncoder, true
		}
		return nil, false
	}
	if kind == reflect.Struct && !e.extRegistry.hasCustom[reflect.Struct] {
		if typ == builtInTimeType {
			return builtInTimeEncoder, true
		}
		return nil, false
	}
	coders := e.extRegistry.byKind[kind]
	if coders == nil {
		return nil, false
	}
	encoder, ok := coders[typ]
	return encoder, ok
}

func (e *encoder) createBuiltIn(rv reflect.Value) error {
	kind := rv.Kind()
	switch kind {
	case reflect.Invalid:
		return e.writeNil()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return e.writeUint(rv.Uint())
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return e.writeInt(rv.Int())
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
		return e.createSliceBuiltIn(rv)
	case reflect.Struct:
		if rv.Type() == builtInTimeType {
			return e.writeExt(builtInTimeEncoder, rv)
		}
		return e.writeStruct(rv)
	case reflect.Ptr:
		if rv.IsNil() {
			return e.writeNil()
		}
		return e.createBuiltIn(rv.Elem())
	case reflect.Interface:
		return e.createBuiltIn(rv.Elem())
	default:
		return e.createByKind(rv, kind)
	}
}

func (e *encoder) createSliceBuiltIn(rv reflect.Value) error {
	if rv.IsNil() {
		return e.writeNil()
	}
	l := rv.Len()
	if e.isByteSlice(rv) {
		if err := e.writeByteSliceLength(l); err != nil {
			return err
		}
		return e.setBytes(rv.Bytes())
	}
	if err := e.writeSliceLength(l); err != nil {
		return err
	}
	if find, err := e.writeFixedSlice(rv); err != nil {
		return err
	} else if find {
		return nil
	}

	elementType := rv.Type().Elem()
	if elementType.Kind() == reflect.Struct {
		if elementType == builtInTimeType {
			for i := 0; i < l; i++ {
				if err := e.writeExt(builtInTimeEncoder, rv.Index(i)); err != nil {
					return err
				}
			}
			return nil
		}
		writeElement := e.getStructWriter(elementType)
		for i := 0; i < l; i++ {
			if err := writeElement(rv.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}

	for i := 0; i < l; i++ {
		if err := e.createBuiltIn(rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}
