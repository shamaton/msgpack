package decoding

import (
	"encoding/binary"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v3/def"
)

type structCacheTypeMap struct {
	keys [][]byte

	// fast path detection
	hasEmbedded bool

	// fast path (hasEmbedded == false): direct field access
	simpleIndexes []int

	// embedded path (hasEmbedded == true): path-based access
	indexes [][]int // field path (support for embedded structs)
}

type structCacheTypeArray struct {
	// fast path detection
	hasEmbedded bool

	// fast path (hasEmbedded == false): direct field access
	simpleIndexes []int

	// embedded path (hasEmbedded == true): path-based access
	indexes [][]int // field path (support for embedded structs)
}

// struct cache map
var (
	mapSCTM = sync.Map{}
	mapSCTA = sync.Map{}
)

// getFieldByPath returns the field value by following the path of indices.
// The bool indicates whether the path was reachable (no nil pointer in the path).
func getFieldByPath(rv reflect.Value, path []int, allowAlloc bool) (reflect.Value, bool) {
	for _, idx := range path {
		// Handle pointer indirection if needed
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				if !allowAlloc {
					return reflect.Value{}, false
				}
				// Allocate new value if pointer is nil
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		rv = rv.Field(idx)
	}
	return rv, true
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

	isExt, _, err := d.extEndOffset(offset)
	if err != nil {
		return 0, err
	}
	if isExt {
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
		fields := d.CollectFields(rv.Type(), nil)

		// detect embedded fields
		hasEmbedded := false
		for _, f := range fields {
			if len(f.Path) > 1 || len(f.OmitPaths) > 0 {
				hasEmbedded = true
				break
			}
		}
		scta.hasEmbedded = hasEmbedded

		for _, field := range fields {
			if hasEmbedded {
				scta.indexes = append(scta.indexes, field.Path)
			} else {
				scta.simpleIndexes = append(scta.simpleIndexes, field.Path[0])
			}
		}
		mapSCTA.Store(rv.Type(), scta)
	} else {
		scta = cache.(*structCacheTypeArray)
	}

	// set value
	if scta.hasEmbedded {
		for i := 0; i < l; i++ {
			if i < len(scta.indexes) {
				allowAlloc := !d.isCodeNil(d.data[o])
				fieldValue, ok := getFieldByPath(rv, scta.indexes[i], allowAlloc)
				if ok {
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
			} else {
				o, err = d.jumpOffset(o)
				if err != nil {
					return 0, err
				}
			}
		}
	} else {
		for i := 0; i < l; i++ {
			if i < len(scta.simpleIndexes) {
				o, err = d.decode(rv.Field(scta.simpleIndexes[i]), o)
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
		fields := d.CollectFields(rv.Type(), nil)

		// detect embedded fields
		hasEmbedded := false
		for _, f := range fields {
			if len(f.Path) > 1 || len(f.OmitPaths) > 0 {
				hasEmbedded = true
				break
			}
		}
		sctm.hasEmbedded = hasEmbedded

		for _, field := range fields {
			sctm.keys = append(sctm.keys, []byte(field.Name))
			if hasEmbedded {
				sctm.indexes = append(sctm.indexes, field.Path)
			} else {
				sctm.simpleIndexes = append(sctm.simpleIndexes, field.Path[0])
			}
		}
		mapSCTM.Store(rv.Type(), sctm)
	} else {
		sctm = cache.(*structCacheTypeMap)
	}

	if sctm.hasEmbedded {
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
				allowAlloc := !d.isCodeNil(d.data[o2])
				fieldValue, ok := getFieldByPath(rv, fieldPath, allowAlloc)
				if ok {
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
			} else {
				o2, err = d.jumpOffset(o2)
				if err != nil {
					return 0, err
				}
			}
			o = o2
		}
	} else {
		for i := 0; i < l; i++ {
			dataKey, o2, err := d.asStringByte(o, k)
			if err != nil {
				return 0, err
			}

			fieldIndex := -1
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
					fieldIndex = sctm.simpleIndexes[keyIndex]
					break
				}
			}

			if fieldIndex >= 0 {
				o2, err = d.decode(rv.Field(fieldIndex), o2)
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

	default:
		isExt, o, err := d.extEndOffsetWithCode(code, offset)
		if err != nil {
			return 0, err
		}
		if isExt {
			offset = o
		}

	}
	return offset, nil
}
