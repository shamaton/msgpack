package encoding

import (
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/rawext"
	"github.com/shamaton/msgpack/v3/internal/time"
)

type extEncoderRegistry struct {
	byKind      [reflect.UnsafePointer + 1]map[reflect.Type]ext.Encoder
	hasCustom   [reflect.UnsafePointer + 1]bool
	customCount int
}

var (
	extEncoderRegistryMu      sync.Mutex
	currentExtEncoderRegistry atomic.Pointer[extEncoderRegistry]
	builtInTimeType           = time.Encoder.Type()
	builtInRawType            = rawext.Encoder.Type()
)

// builtInExtEncoder returns the built-in encoder (time.Time or ext.Raw) for
// typ, if any. It centralizes the two built-in special cases so callers
// don't need to repeat both checks.
func builtInExtEncoder(typ reflect.Type) (ext.Encoder, bool) {
	switch typ {
	case builtInTimeType:
		return time.Encoder, true
	case builtInRawType:
		return rawext.Encoder, true
	default:
		return nil, false
	}
}

func init() {
	registry := &extEncoderRegistry{}
	registry.byKind[reflect.Struct] = map[reflect.Type]ext.Encoder{
		builtInTimeType: time.Encoder,
		builtInRawType:  rawext.Encoder,
	}
	currentExtEncoderRegistry.Store(registry)
}

// AddExtEncoder adds an encoder for an exact type without replacing an existing encoder.
func AddExtEncoder(encoder ext.Encoder) {
	typ := encoder.Type()
	if typ == builtInTimeType || typ == builtInRawType {
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
func RemoveExtEncoder(encoder ext.Encoder) {
	typ := encoder.Type()
	if typ == builtInTimeType || typ == builtInRawType {
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
		// Struct always retains two built-ins (time.Time and ext.Raw), so
		// more than two entries means a custom struct-kind encoder remains.
		next.hasCustom[kind] = kind != reflect.Struct || len(next.byKind[kind]) > 2
	}
	next.customCount--
	currentExtEncoderRegistry.Store(&next)
}

func cloneExtEncoderMap(source map[reflect.Type]ext.Encoder, extraCapacity int) map[reflect.Type]ext.Encoder {
	cloned := make(map[reflect.Type]ext.Encoder, len(source)+extraCapacity)
	for typ, encoder := range source {
		cloned[typ] = encoder
	}
	return cloned
}

func (e *encoder) extEncoderForValue(kind reflect.Kind, value reflect.Value) (ext.Encoder, bool) {
	if e.extRegistry.customCount == 0 {
		if kind == reflect.Struct {
			return builtInExtEncoder(value.Type())
		}
		return nil, false
	}
	if kind == reflect.Struct && !e.extRegistry.hasCustom[reflect.Struct] {
		return builtInExtEncoder(value.Type())
	}
	coders := e.extRegistry.byKind[kind]
	if coders == nil {
		return nil, false
	}
	encoder, ok := coders[value.Type()]
	return encoder, ok
}

func (e *encoder) extEncoderForType(typ reflect.Type) (ext.Encoder, bool) {
	kind := typ.Kind()
	if e.extRegistry.customCount == 0 {
		if kind == reflect.Struct {
			return builtInExtEncoder(typ)
		}
		return nil, false
	}
	if kind == reflect.Struct && !e.extRegistry.hasCustom[reflect.Struct] {
		return builtInExtEncoder(typ)
	}
	coders := e.extRegistry.byKind[kind]
	if coders == nil {
		return nil, false
	}
	encoder, ok := coders[typ]
	return encoder, ok
}

func (e *encoder) calcSizeBuiltIn(rv reflect.Value) (int, error) {
	kind := rv.Kind()
	switch kind {
	case reflect.Invalid:
		return def.Byte1, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return e.calcUint(rv.Uint()), nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return e.calcInt(rv.Int()), nil
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
	case reflect.Struct:
		if rv.Type() == builtInTimeType {
			return time.Encoder.CalcByteSize(rv)
		}
		if rv.Type() == builtInRawType {
			return rawext.Encoder.CalcByteSize(rv)
		}
		return e.calcStruct(rv)
	case reflect.Ptr:
		if rv.IsNil() {
			return def.Byte1, nil
		}
		return e.calcSizeBuiltIn(rv.Elem())
	case reflect.Interface:
		return e.calcSizeBuiltIn(rv.Elem())
	default:
		return e.calcSizeByKind(rv, kind)
	}
}

func (e *encoder) createBuiltIn(rv reflect.Value, offset int) int {
	kind := rv.Kind()
	switch kind {
	case reflect.Invalid:
		return e.writeNil(offset)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return e.writeUint(rv.Uint(), offset)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return e.writeInt(rv.Int(), offset)
	case reflect.Float32:
		return e.writeFloat32(rv.Float(), offset)
	case reflect.Float64:
		return e.writeFloat64(rv.Float(), offset)
	case reflect.Bool:
		return e.writeBool(rv.Bool(), offset)
	case reflect.String:
		return e.writeString(rv.String(), offset)
	case reflect.Complex64:
		return e.writeComplex64(complex64(rv.Complex()), offset)
	case reflect.Complex128:
		return e.writeComplex128(rv.Complex(), offset)
	case reflect.Struct:
		if rv.Type() == builtInTimeType {
			return e.writeExt(time.Encoder, rv, offset)
		}
		if rv.Type() == builtInRawType {
			return e.writeExt(rawext.Encoder, rv, offset)
		}
		return e.writeStruct(rv, offset)
	case reflect.Ptr:
		if rv.IsNil() {
			return e.writeNil(offset)
		}
		return e.createBuiltIn(rv.Elem(), offset)
	case reflect.Interface:
		return e.createBuiltIn(rv.Elem(), offset)
	default:
		return e.createByKind(rv, kind, offset)
	}
}
