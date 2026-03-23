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

func (d *decoder) setStruct(code byte, rv reflect.Value, k reflect.Kind) error {
	if len(extCoders) > 0 {
		innerType, data, err := d.readIfExtType(code)
		if err != nil {
			return err
		}
		if data != nil {
			for i := range extCoders {
				if extCoders[i].IsType(code, innerType, len(data)) {
					v, err := extCoders[i].ToValue(code, data, k)
					if err != nil {
						return err
					}

					// Validate that the receptacle is of the right value type.
					if rv.Type() == reflect.TypeOf(v) {
						rv.Set(reflect.ValueOf(v))
						return nil
					}
				}
			}
		}
	}

	if d.asArray {
		return d.setStructFromArray(code, rv, k)
	}
	return d.setStructFromMap(code, rv, k)
}

func (d *decoder) setStructFromArray(code byte, rv reflect.Value, k reflect.Kind) error {
	// get length
	l, err := d.sliceLength(code, k)
	if err != nil {
		return err
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
				code, err := d.readSize1()
				if err != nil {
					return err
				}
				allowAlloc := !d.isCodeNil(code)
				fieldValue, ok := getFieldByPath(rv, scta.indexes[i], allowAlloc)
				if ok {
					err = d.decodeWithCode(code, fieldValue)
					if err != nil {
						return err
					}
				} else if !d.isCodeNil(code) {
					return d.errorTemplate(code, k)
				}
			} else {
				err = d.jumpOffset()
				if err != nil {
					return err
				}
			}
		}
	} else {
		for i := 0; i < l; i++ {
			if i < len(scta.simpleIndexes) {
				err = d.decode(rv.Field(scta.simpleIndexes[i]))
				if err != nil {
					return err
				}
			} else {
				err = d.jumpOffset()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *decoder) setStructFromMap(code byte, rv reflect.Value, k reflect.Kind) error {
	// get length
	l, err := d.mapLength(code, k)
	if err != nil {
		return err
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
			dataKey, err := d.asStringByte(k)
			if err != nil {
				return err
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
				code, err := d.readSize1()
				if err != nil {
					return err
				}
				allowAlloc := !d.isCodeNil(code)
				fieldValue, ok := getFieldByPath(rv, fieldPath, allowAlloc)
				if ok {
					err = d.decodeWithCode(code, fieldValue)
					if err != nil {
						return err
					}
				} else if !d.isCodeNil(code) {
					return d.errorTemplate(code, k)
				}
			} else {
				err = d.jumpOffset()
				if err != nil {
					return err
				}
			}
		}
	} else {
		for i := 0; i < l; i++ {
			dataKey, err := d.asStringByte(k)
			if err != nil {
				return err
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
				err = d.decode(rv.Field(fieldIndex))
				if err != nil {
					return err
				}
			} else {
				err = d.jumpOffset()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *decoder) jumpOffset() error {
	code, err := d.readSize1()
	if err != nil {
		return err
	}

	switch {
	case code == def.True, code == def.False, code == def.Nil:
		// do nothing

	case d.isPositiveFixNum(code) || d.isNegativeFixNum(code):
		// do nothing
	case code == def.Uint8, code == def.Int8:
		_, err = d.readSize1()
		return err
	case code == def.Uint16, code == def.Int16:
		_, err = d.readSize2()
		return err
	case code == def.Uint32, code == def.Int32, code == def.Float32:
		_, err = d.readSize4()
		return err
	case code == def.Uint64, code == def.Int64, code == def.Float64:
		_, err = d.readSize8()
		return err

	case d.isFixString(code):
		_, err = d.readSizeN(int(code - def.FixStr))
		return err
	case code == def.Str8, code == def.Bin8:
		b, err := d.readSize1()
		if err != nil {
			return err
		}
		_, err = d.readSizeN(int(b))
		return err
	case code == def.Str16, code == def.Bin16:
		bs, err := d.readSize2()
		if err != nil {
			return err
		}
		_, err = d.readSizeN(int(binary.BigEndian.Uint16(bs)))
		return err
	case code == def.Str32, code == def.Bin32:
		bs, err := d.readSize4()
		if err != nil {
			return err
		}
		_, err = d.readSizeN(int(binary.BigEndian.Uint32(bs)))
		return err

	case d.isFixSlice(code):
		l := int(code - def.FixArray)
		for i := 0; i < l; i++ {
			if err = d.jumpOffset(); err != nil {
				return err
			}
		}
	case code == def.Array16:
		bs, err := d.readSize2()
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			if err = d.jumpOffset(); err != nil {
				return err
			}
		}
	case code == def.Array32:
		bs, err := d.readSize4()
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			if err = d.jumpOffset(); err != nil {
				return err
			}
		}

	case d.isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			if err = d.jumpOffset(); err != nil {
				return err
			}
		}
	case code == def.Map16:
		bs, err := d.readSize2()
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			if err = d.jumpOffset(); err != nil {
				return err
			}
		}
	case code == def.Map32:
		bs, err := d.readSize4()
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			if err = d.jumpOffset(); err != nil {
				return err
			}
		}

	case code == def.Fixext1:
		_, err = d.readSizeN(def.Byte1 + def.Byte1)
		return err
	case code == def.Fixext2:
		_, err = d.readSizeN(def.Byte1 + def.Byte2)
		return err
	case code == def.Fixext4:
		_, err = d.readSizeN(def.Byte1 + def.Byte4)
		return err
	case code == def.Fixext8:
		_, err = d.readSizeN(def.Byte1 + def.Byte8)
		return err
	case code == def.Fixext16:
		_, err = d.readSizeN(def.Byte1 + def.Byte16)
		return err

	case code == def.Ext8:
		b, err := d.readSize1()
		if err != nil {
			return err
		}
		_, err = d.readSizeN(def.Byte1 + int(b))
		return err
	case code == def.Ext16:
		bs, err := d.readSize2()
		if err != nil {
			return err
		}
		_, err = d.readSizeN(def.Byte1 + int(binary.BigEndian.Uint16(bs)))
		return err
	case code == def.Ext32:
		bs, err := d.readSize4()
		if err != nil {
			return err
		}
		_, err = d.readSizeN(def.Byte1 + int(binary.BigEndian.Uint32(bs)))
		return err
	}
	return nil
}
