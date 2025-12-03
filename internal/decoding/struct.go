package decoding

import (
	"encoding/binary"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v2/def"
)

type structCacheTypeMap struct {
	keys    [][]byte
	indexes [][]int // field path (support for embedded structs)
}

type structCacheTypeArray struct {
	m [][]int // field path (support for embedded structs)
}

// struct cache map
var mapSCTM = sync.Map{}
var mapSCTA = sync.Map{}

// fieldInfo holds information about a struct field including its path for embedded structs
type fieldInfo struct {
	path  []int  // path to reach this field (indices for embedded structs)
	name  string // field name or tag
	index int    // field index in the struct
}

// collectFields collects all fields from a struct, expanding embedded structs
// following the same rules as encoding/json
func (d *decoder) collectFields(t reflect.Type, path []int) []fieldInfo {
	var fields []fieldInfo
	var embedded []fieldInfo // embedded fields to process later (lower priority)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check field visibility
		public, _, name := d.CheckField(field)
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
				embeddedFields := d.collectFields(fieldType, newPath)
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

func (d *decoder) isPublic(name string) bool {
	return len(name) > 0 && 0x41 <= name[0] && name[0] <= 0x5a
}

// getFieldByPath returns the field value by following the path of indices
func getFieldByPath(rv reflect.Value, path []int) reflect.Value {
	for _, idx := range path {
		// Handle pointer indirection if needed
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				// Allocate new value if pointer is nil
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		rv = rv.Field(idx)
	}
	return rv
}

func (d *decoder) setStruct(rv reflect.Value, offset int, k reflect.Kind) (int, error) {
	/*
		if d.isDateTime(offset) {
			dt, offset, err := d.asDateTime(offset, k)
			if err != nil {
				return 0, err
			}
			rv.Set(reflect.ValueOf(dt))
			return offset, nil
		}
	*/

	for i := range extCoders {
		if extCoders[i].IsType(offset, &d.data) {
			v, offset, err := extCoders[i].AsValue(offset, k, &d.data)
			if err != nil {
				return 0, err
			}

			// Validate that the receptacle is of the right value type.
			if rv.Type() == reflect.TypeOf(v) {
				rv.Set(reflect.ValueOf(v))
				return offset, nil
			}
		}
	}

	if d.asArray {
		return d.setStructFromArray(rv, offset, k)
	}
	return d.setStructFromMap(rv, offset, k)
}

func (d *decoder) setStructFromArray(rv reflect.Value, offset int, k reflect.Kind) (int, error) {
	// get length
	l, o, err := d.sliceLength(offset, k)
	if err != nil {
		return 0, err
	}

	if err = d.hasRequiredLeastSliceSize(o, l); err != nil {
		return 0, err
	}

	// find or create reference
	var scta *structCacheTypeArray
	cache, findCache := mapSCTA.Load(rv.Type())
	if !findCache {
		scta = &structCacheTypeArray{}
		fields := d.collectFields(rv.Type(), nil)
		for _, field := range fields {
			scta.m = append(scta.m, field.path)
		}
		mapSCTA.Store(rv.Type(), scta)
	} else {
		scta = cache.(*structCacheTypeArray)
	}
	// set value
	for i := 0; i < l; i++ {
		if i < len(scta.m) {
			fieldValue := getFieldByPath(rv, scta.m[i])
			o, err = d.decode(fieldValue, o)
			if err != nil {
				return 0, err
			}
		} else {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
	}
	return o, nil
}

func (d *decoder) setStructFromMap(rv reflect.Value, offset int, k reflect.Kind) (int, error) {
	// get length
	l, o, err := d.mapLength(offset, k)
	if err != nil {
		return 0, err
	}

	if err = d.hasRequiredLeastMapSize(o, l); err != nil {
		return 0, err
	}

	var sctm *structCacheTypeMap
	cache, cacheFind := mapSCTM.Load(rv.Type())
	if !cacheFind {
		sctm = &structCacheTypeMap{}
		fields := d.collectFields(rv.Type(), nil)
		for _, field := range fields {
			sctm.keys = append(sctm.keys, []byte(field.name))
			sctm.indexes = append(sctm.indexes, field.path)
		}
		mapSCTM.Store(rv.Type(), sctm)
	} else {
		sctm = cache.(*structCacheTypeMap)
	}

	for i := 0; i < l; i++ {
		dataKey, o2, err := d.asStringByte(o, k)
		if err != nil {
			return 0, err
		}

		fieldPath := []int(nil)
		for keyIndex, keyBytes := range sctm.keys {
			if len(keyBytes) != len(dataKey) {
				continue
			}

			found := true
			for dataIndex := range dataKey {
				if dataKey[dataIndex] != keyBytes[dataIndex] {
					found = false
					break
				}
			}
			if found {
				fieldPath = sctm.indexes[keyIndex]
				break
			}
		}

		if fieldPath != nil {
			fieldValue := getFieldByPath(rv, fieldPath)
			o2, err = d.decode(fieldValue, o2)
			if err != nil {
				return 0, err
			}
		} else {
			o2, err = d.jumpOffset(o2)
			if err != nil {
				return 0, err
			}
		}
		o = o2
	}
	return o, nil
}

func (d *decoder) jumpOffset(offset int) (int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return 0, err
	}

	switch {
	case code == def.True, code == def.False, code == def.Nil:
		// do nothing

	case d.isPositiveFixNum(code) || d.isNegativeFixNum(code):
		// do nothing
	case code == def.Uint8, code == def.Int8:
		offset += def.Byte1
	case code == def.Uint16, code == def.Int16:
		offset += def.Byte2
	case code == def.Uint32, code == def.Int32, code == def.Float32:
		offset += def.Byte4
	case code == def.Uint64, code == def.Int64, code == def.Float64:
		offset += def.Byte8

	case d.isFixString(code):
		offset += int(code - def.FixStr)
	case code == def.Str8, code == def.Bin8:
		b, o, err := d.readSize1(offset)
		if err != nil {
			return 0, err
		}
		o += int(b)
		offset = o
	case code == def.Str16, code == def.Bin16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		o += int(binary.BigEndian.Uint16(bs))
		offset = o
	case code == def.Str32, code == def.Bin32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		o += int(binary.BigEndian.Uint32(bs))
		offset = o

	case d.isFixSlice(code):
		l := int(code - def.FixArray)
		for i := 0; i < l; i++ {
			offset, err = d.jumpOffset(offset)
			if err != nil {
				return 0, err
			}
		}
	case code == def.Array16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o
	case code == def.Array32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o

	case d.isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			offset, err = d.jumpOffset(offset)
			if err != nil {
				return 0, err
			}
		}
	case code == def.Map16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o
	case code == def.Map32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o

	case code == def.Fixext1:
		offset += def.Byte1 + def.Byte1
	case code == def.Fixext2:
		offset += def.Byte1 + def.Byte2
	case code == def.Fixext4:
		offset += def.Byte1 + def.Byte4
	case code == def.Fixext8:
		offset += def.Byte1 + def.Byte8
	case code == def.Fixext16:
		offset += def.Byte1 + def.Byte16

	case code == def.Ext8:
		b, o, err := d.readSize1(offset)
		if err != nil {
			return 0, err
		}
		o += def.Byte1 + int(b)
		offset = o
	case code == def.Ext16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		o += def.Byte1 + int(binary.BigEndian.Uint16(bs))
		offset = o
	case code == def.Ext32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		o += def.Byte1 + int(binary.BigEndian.Uint32(bs))
		offset = o

	}
	return offset, nil
}
