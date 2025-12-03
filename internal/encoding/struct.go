package encoding

import (
	"math"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v2/def"
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

type structCalcFunc func(rv reflect.Value) (int, error)
type structWriteFunc func(rv reflect.Value, offset int) int

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

func (e *encoder) getStructCalc(typ reflect.Type) structCalcFunc {

	for j := range extCoders {
		if extCoders[j].Type() == typ {
			return extCoders[j].CalcByteSize
		}
	}
	if e.asArray {
		return e.calcStructArray
	}
	return e.calcStructMap

}

func (e *encoder) calcStruct(rv reflect.Value) (int, error) {

	//if isTime, tm := e.isDateTime(rv); isTime {
	//	size := e.calcTime(tm)
	//	return size, nil
	//}

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].CalcByteSize(rv)
		}
	}

	if e.asArray {
		return e.calcStructArray(rv)
	}
	return e.calcStructMap(rv)
}

func (e *encoder) calcStructArray(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	cache, find := cachemap.Load(t)
	var c *structCache
	if !find {
		c = &structCache{}
		fields := e.collectFields(t, nil)
		omitCount := 0
		for _, field := range fields {
			fieldValue := getFieldByPath(rv, field.path)
			if !fieldValue.IsValid() {
				continue
			}
			size, err := e.calcSize(fieldValue)
			if err != nil {
				return 0, err
			}
			ret += size
			c.indexes = append(c.indexes, field.path)
			c.names = append(c.names, field.name)
			c.omits = append(c.omits, field.omit)
			if field.omit {
				omitCount++
			}
		}
		c.noOmit = omitCount == 0
		cachemap.Store(t, c)
	} else {
		c = cache.(*structCache)
		for i := 0; i < len(c.indexes); i++ {
			fieldValue := getFieldByPath(rv, c.indexes[i])
			if !fieldValue.IsValid() {
				continue
			}
			size, err := e.calcSize(fieldValue)
			if err != nil {
				return 0, err
			}
			ret += size
		}
	}

	// format size
	size, err := e.calcLength(len(c.indexes))
	if err != nil {
		return 0, err
	}
	ret += size
	return ret, nil
}

func (e *encoder) calcStructMap(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	cache, find := cachemap.Load(t)
	var c *structCache
	var l int
	if !find {
		c = &structCache{}
		fields := e.collectFields(t, nil)
		omitCount := 0
		for _, field := range fields {
			fieldValue := getFieldByPath(rv, field.path)
			if !fieldValue.IsValid() {
				continue
			}
			size, err := e.calcSizeWithOmitEmpty(fieldValue, field.name, field.omit)
			if err != nil {
				return 0, err
			}
			ret += size
			c.indexes = append(c.indexes, field.path)
			c.names = append(c.names, field.name)
			c.omits = append(c.omits, field.omit)
			if field.omit {
				omitCount++
			}
			if size > 0 {
				l++
			}
		}
		c.noOmit = omitCount == 0
		cachemap.Store(t, c)
	} else {
		c = cache.(*structCache)
		for i := 0; i < len(c.indexes); i++ {
			fieldValue := getFieldByPath(rv, c.indexes[i])
			if !fieldValue.IsValid() {
				continue
			}
			size, err := e.calcSizeWithOmitEmpty(fieldValue, c.names[i], c.omits[i])
			if err != nil {
				return 0, err
			}
			ret += size
			if size > 0 {
				l++
			}
		}
	}

	// format size
	size, err := e.calcLength(len(c.indexes))
	if err != nil {
		return 0, err
	}
	ret += size
	return ret, nil
}

func (e *encoder) calcSizeWithOmitEmpty(rv reflect.Value, name string, omit bool) (int, error) {
	keySize := 0
	valueSize := 0
	if !omit || !rv.IsZero() {
		keySize = e.calcString(name)
		vSize, err := e.calcSize(rv)
		if err != nil {
			return 0, err
		}
		valueSize = vSize
	}
	return keySize + valueSize, nil
}

func (e *encoder) getStructWriter(typ reflect.Type) structWriteFunc {

	for i := range extCoders {
		if extCoders[i].Type() == typ {
			return func(rv reflect.Value, offset int) int {
				return extCoders[i].WriteToBytes(rv, offset, &e.d)
			}
		}
	}

	if e.asArray {
		return e.writeStructArray
	}
	return e.writeStructMap
}

func (e *encoder) writeStruct(rv reflect.Value, offset int) int {
	/*
		if isTime, tm := e.isDateTime(rv); isTime {
			return e.writeTime(tm, offset)
		}
	*/

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].WriteToBytes(rv, offset, &e.d)
		}
	}

	if e.asArray {
		return e.writeStructArray(rv, offset)
	}
	return e.writeStructMap(rv, offset)
}

func (e *encoder) writeStructArray(rv reflect.Value, offset int) int {

	cache, _ := cachemap.Load(rv.Type())
	c := cache.(*structCache)

	// write format
	num := len(c.indexes)
	if num <= 0x0f {
		offset = e.setByte1Int(def.FixArray+num, offset)
	} else if num <= math.MaxUint16 {
		offset = e.setByte1Int(def.Array16, offset)
		offset = e.setByte2Int(num, offset)
	} else if uint(num) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Array32, offset)
		offset = e.setByte4Int(num, offset)
	}

	for i := 0; i < num; i++ {
		fieldValue := getFieldByPath(rv, c.indexes[i])
		offset = e.create(fieldValue, offset)
	}
	return offset
}

func (e *encoder) writeStructMap(rv reflect.Value, offset int) int {

	cache, _ := cachemap.Load(rv.Type())
	c := cache.(*structCache)

	// format size
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

	if l <= 0x0f {
		offset = e.setByte1Int(def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(l, offset)
	}

	for i := 0; i < num; i++ {
		fieldValue := getFieldByPath(rv, c.indexes[i])
		if !c.omits[i] || !fieldValue.IsZero() {
			offset = e.writeString(c.names[i], offset)
			offset = e.create(fieldValue, offset)
		}
	}
	return offset
}
