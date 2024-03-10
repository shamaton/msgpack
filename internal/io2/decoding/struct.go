package decoding

import (
	"encoding/binary"
	"github.com/shamaton/msgpack/v2/internal/common"
	"io"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v2/def"
)

type structCacheTypeMap struct {
	keys    [][]byte
	indexes []int
}

type structCacheTypeArray struct {
	m []int
}

// struct cache map
var mapSCTM = sync.Map{}
var mapSCTA = sync.Map{}

func setStruct(r io.Reader, code byte, rv reflect.Value, k reflect.Kind, asArray bool) error {
	if len(extCoders) > 0 {
		innerType, data, err := readIfExtType(r, code)
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

	if asArray {
		return setStructFromArray(r, code, rv, k)
	}
	return setStructFromMap(r, code, rv, k)
}

func setStructFromArray(r io.Reader, code byte, rv reflect.Value, k reflect.Kind) error {
	// get length
	l, err := sliceLength(r, code, k)
	if err != nil {
		return err
	}

	//if err = d.hasRequiredLeastSliceSize(o, l); err != nil {
	//	return err
	//}

	// find or create reference
	var scta *structCacheTypeArray
	cache, findCache := mapSCTA.Load(rv.Type())
	if !findCache {
		scta = &structCacheTypeArray{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, _ := common.CheckField(rv.Type().Field(i)); ok {
				scta.m = append(scta.m, i)
			}
		}
		mapSCTA.Store(rv.Type(), scta)
	} else {
		scta = cache.(*structCacheTypeArray)
	}
	// set value
	for i := 0; i < l; i++ {
		if i < len(scta.m) {
			err = decode(r, rv.Field(scta.m[i]), true)
			if err != nil {
				return err
			}
		} else {
			err = jumpOffset(r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setStructFromMap(r io.Reader, code byte, rv reflect.Value, k reflect.Kind) error {
	// get length
	l, err := mapLength(r, code, k)
	if err != nil {
		return err
	}

	//if err = d.hasRequiredLeastMapSize(o, l); err != nil {
	//	return 0, err
	//}

	var sctm *structCacheTypeMap
	cache, cacheFind := mapSCTM.Load(rv.Type())
	if !cacheFind {
		sctm = &structCacheTypeMap{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, name := common.CheckField(rv.Type().Field(i)); ok {
				sctm.keys = append(sctm.keys, []byte(name))
				sctm.indexes = append(sctm.indexes, i)
			}
		}
		mapSCTM.Store(rv.Type(), sctm)
	} else {
		sctm = cache.(*structCacheTypeMap)
	}

	for i := 0; i < l; i++ {
		dataKey, err := asStringByte(r, k)
		if err != nil {
			return err
		}

		fieldIndex := -1
		for keyIndex, keyBytes := range sctm.keys {
			if len(keyBytes) != len(dataKey) {
				continue
			}

			fieldIndex = sctm.indexes[keyIndex]
			for dataIndex := range dataKey {
				if dataKey[dataIndex] != keyBytes[dataIndex] {
					fieldIndex = -1
					break
				}
			}
			if fieldIndex >= 0 {
				break
			}
		}

		if fieldIndex >= 0 {
			err = decode(r, rv.Field(fieldIndex), false)
			if err != nil {
				return err
			}
		} else {
			err = jumpOffset(r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func jumpOffset(r io.Reader) error {
	code, err := readSize1(r)
	if err != nil {
		return err
	}

	switch {
	case code == def.True, code == def.False, code == def.Nil:
		// do nothing

	case isPositiveFixNum(code) || isNegativeFixNum(code):
		// do nothing
	case code == def.Uint8, code == def.Int8:
		_, err = readSize1(r)
		return err
	case code == def.Uint16, code == def.Int16:
		_, err = readSize2(r)
		return err
	case code == def.Uint32, code == def.Int32, code == def.Float32:
		_, err = readSize4(r)
		return err
	case code == def.Uint64, code == def.Int64, code == def.Float64:
		_, err = readSize8(r)
		return err

	case isFixString(code):
		_, err = readSizeN(r, int(code-def.FixStr))
		return err
	case code == def.Str8, code == def.Bin8:
		b, err := readSize1(r)
		if err != nil {
			return err
		}
		_, err = readSizeN(r, int(b))
		return err
	case code == def.Str16, code == def.Bin16:
		bs, err := readSize2(r)
		if err != nil {
			return err
		}
		_, err = readSizeN(r, int(binary.BigEndian.Uint16(bs)))
		return err
	case code == def.Str32, code == def.Bin32:
		bs, err := readSize4(r)
		if err != nil {
			return err
		}
		_, err = readSizeN(r, int(binary.BigEndian.Uint32(bs)))
		return err

	case isFixSlice(code):
		l := int(code - def.FixArray)
		for i := 0; i < l; i++ {
			if err = jumpOffset(r); err != nil {
				return err
			}
		}
	case code == def.Array16:
		bs, err := readSize2(r)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			if err = jumpOffset(r); err != nil {
				return err
			}
		}
	case code == def.Array32:
		bs, err := readSize4(r)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			if err = jumpOffset(r); err != nil {
				return err
			}
		}

	case isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			if err = jumpOffset(r); err != nil {
				return err
			}
		}
	case code == def.Map16:
		bs, err := readSize2(r)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			if err = jumpOffset(r); err != nil {
				return err
			}
		}
	case code == def.Map32:
		bs, err := readSize4(r)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			if err = jumpOffset(r); err != nil {
				return err
			}
		}

	case code == def.Fixext1:
		_, err = readSizeN(r, def.Byte1+def.Byte1)
	case code == def.Fixext2:
		_, err = readSizeN(r, def.Byte1+def.Byte2)
	case code == def.Fixext4:
		_, err = readSizeN(r, def.Byte1+def.Byte4)
	case code == def.Fixext8:
		_, err = readSizeN(r, def.Byte1+def.Byte8)
	case code == def.Fixext16:
		_, err = readSizeN(r, def.Byte1+def.Byte16)

	case code == def.Ext8:
		b, err := readSize1(r)
		if err != nil {
			return err
		}
		_, err = readSizeN(r, def.Byte1+int(b))
		return err
	case code == def.Ext16:
		bs, err := readSize2(r)
		if err != nil {
			return err
		}
		_, err = readSizeN(r, def.Byte1+int(binary.BigEndian.Uint16(bs)))
		return err
	case code == def.Ext32:
		bs, err := readSize4(r)
		if err != nil {
			return err
		}
		_, err = readSizeN(r, def.Byte1+int(binary.BigEndian.Uint32(bs)))
		return err
	}
	return nil
}
