package encoding

import (
	"math"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	"github.com/shamaton/msgpack/v2/internal/common"
)

type structCache struct {
	indexes [][]int // field path (support for embedded structs)
	names   []string
	omits   []bool
	noOmit  bool
	common.Common
}

var cachemap = sync.Map{}

type structWriteFunc func(rv reflect.Value) error

// fieldInfo holds information about a struct field including its path for embedded structs
type fieldInfo struct {
	path  []int  // path to reach this field (indices for embedded structs)
	name  string // field name or tag
	index int    // field index in the struct
	omit  bool   // omitempty flag
}

// collectFields collects all fields from a struct, expanding embedded structs
// following the same rules as encoding/json
func (e *encoder) collectFields(t reflect.Type, path []int) []fieldInfo {
	var fields []fieldInfo
	var embedded []fieldInfo // embedded fields to process later (lower priority)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check field visibility and get omitempty flag
		public, omit, name := e.CheckField(field)
		if !public {
			continue
		}

		// Get tag to check if embedded
		tag := field.Tag.Get("msgpack")
		// Extract just the name part (before comma if any)
		tagName := tag
		if idx := len(tag); idx > 0 {
			for j, c := range tag {
				if c == ',' {
					tagName = tag[:j]
					break
				}
			}
		}

		// Check if this is an embedded struct
		isEmbedded := field.Anonymous && (tag == "" || tagName == "")

		if isEmbedded {
			// Get the actual type (dereference pointer if needed)
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			// If it's a struct, expand its fields
			if fieldType.Kind() == reflect.Struct {
				newPath := append(append([]int{}, path...), i)
				embeddedFields := e.collectFields(fieldType, newPath)
				embedded = append(embedded, embeddedFields...)
				continue
			}
		}

		// Regular field or embedded non-struct
		newPath := append(append([]int{}, path...), i)
		fields = append(fields, fieldInfo{
			path:  newPath,
			name:  name,
			index: i,
			omit:  omit,
		})
	}

	// Add embedded fields after regular fields (they have lower priority)
	fields = append(fields, embedded...)

	// Remove duplicates and handle ambiguous fields
	// Group fields by name and depth, preserving order
	type fieldAtDepth struct {
		field fieldInfo
		depth int
	}
	fieldsByName := make(map[string][]fieldAtDepth)
	var seenNames []string // To preserve order

	for _, f := range fields {
		if _, seen := fieldsByName[f.name]; !seen {
			seenNames = append(seenNames, f.name)
		}
		fieldsByName[f.name] = append(fieldsByName[f.name], fieldAtDepth{
			field: f,
			depth: len(f.path),
		})
	}

	var result []fieldInfo
	for _, name := range seenNames {
		fieldsWithDepth := fieldsByName[name]

		// Find minimum depth
		minDepth := fieldsWithDepth[0].depth
		for _, fd := range fieldsWithDepth {
			if fd.depth < minDepth {
				minDepth = fd.depth
			}
		}

		// Count fields at minimum depth
		var fieldsAtMinDepth []fieldInfo
		for _, fd := range fieldsWithDepth {
			if fd.depth == minDepth {
				fieldsAtMinDepth = append(fieldsAtMinDepth, fd.field)
			}
		}

		// If there's exactly one field at minimum depth, use it
		// If there are multiple fields at the same minimum depth, it's ambiguous - skip it
		if len(fieldsAtMinDepth) == 1 {
			result = append(result, fieldsAtMinDepth[0])
		}
		// else: ambiguous field, skip it (following encoding/json behavior)
	}

	return result
}

func (e *encoder) isPublic(name string) bool {
	return len(name) > 0 && 0x41 <= name[0] && name[0] <= 0x5a
}

// getFieldByPath returns the field value by following the path of indices
func getFieldByPath(rv reflect.Value, path []int) reflect.Value {
	for _, idx := range path {
		// Handle pointer indirection if needed
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				// Return zero value if pointer is nil
				return reflect.Value{}
			}
			rv = rv.Elem()
		}
		rv = rv.Field(idx)
	}
	return rv
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
	num := len(c.indexes)
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

	for i := 0; i < num; i++ {
		fieldValue := getFieldByPath(rv, c.indexes[i])
		if err := e.create(fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) writeStructMap(rv reflect.Value) error {
	c := e.getStructCache(rv)

	num := len(c.indexes)
	l := 0
	if c.noOmit {
		l = num
	} else {
		for i := 0; i < num; i++ {
			fieldValue := getFieldByPath(rv, c.indexes[i])
			if !c.omits[i] || !fieldValue.IsZero() {
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

	for i := 0; i < num; i++ {
		fieldValue := getFieldByPath(rv, c.indexes[i])
		if !c.omits[i] || !fieldValue.IsZero() {
			if err := e.writeString(c.names[i]); err != nil {
				return err
			}
			if err := e.create(fieldValue); err != nil {
				return err
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
	fields := e.collectFields(t, nil)
	omitCount := 0
	for _, field := range fields {
		c.indexes = append(c.indexes, field.path)
		c.names = append(c.names, field.name)
		c.omits = append(c.omits, field.omit)
		if field.omit {
			omitCount++
		}
	}
	c.noOmit = omitCount == 0
	cachemap.Store(t, c)
	return c
}
