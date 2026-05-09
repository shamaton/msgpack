package encoding

import (
	"math"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/common"
)

type structCache struct {
	// common fields
	names  []string
	omits  []bool
	noOmit bool

	// fast path detection
	hasEmbedded bool

	// fast path (hasEmbedded == false): direct field access
	simpleIndexes []int

	// embedded path (hasEmbedded == true): path-based access
	indexes   [][]int   // field path (support for embedded structs)
	omitPaths [][][]int // embedded omitempty parent paths

	common.Common
}

var cachemap = sync.Map{}

type structWriteFunc func(rv reflect.Value) error

// getFieldByPath returns the field value by following the path of indices.
// The bool indicates whether the path was reachable (no nil pointer in the path).
func getFieldByPath(rv reflect.Value, path []int) (reflect.Value, bool) {
	for _, idx := range path {
		// Handle pointer indirection if needed
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				// Return invalid value if pointer is nil
				return reflect.Value{}, false
			}
			rv = rv.Elem()
		}
		rv = rv.Field(idx)
	}
	return rv, true
}

func shouldOmitByParent(rv reflect.Value, omitPaths [][]int) bool {
	for _, path := range omitPaths {
		parentValue, ok := getFieldByPath(rv, path)
		if !ok || parentValue.IsZero() {
			return true
		}
	}
	return false
}

func (e *encoder) getStructWriter(typ reflect.Type) structWriteFunc {
	for i := range extCoders {
		if extCoders[i].Type() == typ {
			return func(rv reflect.Value) error {
				w := ext.CreateStreamWriter(e.w, e.buf)
				return extCoders[i].Write(w, rv)
			}
		}
	}

	if e.asArray {
		return e.writeStructArray
	}
	return e.writeStructMap
}

func (e *encoder) writeStruct(rv reflect.Value) error {
	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			w := ext.CreateStreamWriter(e.w, e.buf)
			return extCoders[i].Write(w, rv)
		}
	}

	if e.asArray {
		return e.writeStructArray(rv)
	}
	return e.writeStructMap(rv)
}

func (e *encoder) writeStructArray(rv reflect.Value) error {
	c := e.getStructCache(rv)

	// write format
	var num int
	if c.hasEmbedded {
		num = len(c.indexes)
	} else {
		num = len(c.simpleIndexes)
	}

	if num <= 0x0f {
		if err := e.setByte1Int(def.FixArray + num); err != nil {
			return err
		}
	} else if num <= math.MaxUint16 {
		if err := e.setByte1Int(def.Array16); err != nil {
			return err
		}
		if err := e.setByte2Int(num); err != nil {
			return err
		}
	} else if uint(num) <= math.MaxUint32 {
		if err := e.setByte1Int(def.Array32); err != nil {
			return err
		}
		if err := e.setByte4Int(num); err != nil {
			return err
		}
	}

	if c.hasEmbedded {
		for i := 0; i < num; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				fieldValue = reflect.Value{}
			}
			if err := e.create(fieldValue); err != nil {
				return err
			}
		}
	} else {
		for i := 0; i < num; i++ {
			if err := e.create(rv.Field(c.simpleIndexes[i])); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *encoder) writeStructMap(rv reflect.Value) error {
	c := e.getStructCache(rv)

	l := 0
	if c.hasEmbedded {
		num := len(c.indexes)
		for i := 0; i < num; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				continue
			}
			if c.noOmit || !c.omits[i] || !fieldValue.IsZero() {
				l++
			}
		}
	} else {
		num := len(c.simpleIndexes)
		for i := 0; i < num; i++ {
			if c.noOmit || !c.omits[i] || !rv.Field(c.simpleIndexes[i]).IsZero() {
				l++
			}
		}
	}

	// format size
	if l <= 0x0f {
		if err := e.setByte1Int(def.FixMap + l); err != nil {
			return err
		}
	} else if l <= math.MaxUint16 {
		if err := e.setByte1Int(def.Map16); err != nil {
			return err
		}
		if err := e.setByte2Int(l); err != nil {
			return err
		}
	} else if uint(l) <= math.MaxUint32 {
		if err := e.setByte1Int(def.Map32); err != nil {
			return err
		}
		if err := e.setByte4Int(l); err != nil {
			return err
		}
	}

	if c.hasEmbedded {
		num := len(c.indexes)
		for i := 0; i < num; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				continue
			}
			if c.noOmit || !c.omits[i] || !fieldValue.IsZero() {
				if err := e.writeString(c.names[i]); err != nil {
					return err
				}
				if err := e.create(fieldValue); err != nil {
					return err
				}
			}
		}
	} else {
		num := len(c.simpleIndexes)
		for i := 0; i < num; i++ {
			fieldValue := rv.Field(c.simpleIndexes[i])
			if c.noOmit || !c.omits[i] || !fieldValue.IsZero() {
				if err := e.writeString(c.names[i]); err != nil {
					return err
				}
				if err := e.create(fieldValue); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *encoder) getStructCache(rv reflect.Value) *structCache {
	t := rv.Type()
	cache, find := cachemap.Load(t)
	if find {
		return cache.(*structCache)
	}

	c := &structCache{}
	fields := e.CollectFields(t, nil)

	// detect embedded fields
	hasEmbedded := false
	for _, f := range fields {
		if len(f.Path) > 1 || len(f.OmitPaths) > 0 {
			hasEmbedded = true
			break
		}
	}
	c.hasEmbedded = hasEmbedded

	omitCount := 0
	for _, field := range fields {
		c.names = append(c.names, field.Name)
		c.omits = append(c.omits, field.Omit)
		if hasEmbedded {
			c.indexes = append(c.indexes, field.Path)
			c.omitPaths = append(c.omitPaths, field.OmitPaths)
		} else {
			c.simpleIndexes = append(c.simpleIndexes, field.Path[0])
		}
		if field.Omit {
			omitCount++
		}
	}
	c.noOmit = omitCount == 0
	cachemap.Store(t, c)
	return c
}
