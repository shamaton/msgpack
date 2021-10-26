package decoding

import (
	"bufio"
	"encoding/binary"
	"fmt"
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

func (d *decoder) setStruct(rv reflect.Value, reader *bufio.Reader, k reflect.Kind) error {
	/*
		if d.isDateTime(reader) {
			dt, offset, err := d.asDateTime(offset, k)
			if err != nil {
				return 0, err
			}
			rv.Set(reflect.ValueOf(dt))
			return offset, nil
		}
	*/

	if code, data, err := d.readExt(reader); err == nil {
		for i := range extCoders {
			if extCoders[i].Code() == int8(code) {
				v, err := extCoders[i].AsValue(data, k)
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

	if d.asArray {
		return d.setStructFromArray(rv, reader, k)
	}
	return d.setStructFromMap(rv, reader, k)
}

func (d *decoder) setStructFromArray(rv reflect.Value, reader *bufio.Reader, k reflect.Kind) error {
	// get length
	l, err := d.sliceLength(reader, k)
	if err != nil {
		return err
	}

	// find or create reference
	var scta *structCacheTypeArray
	cache, findCache := mapSCTA.Load(rv.Type())
	if !findCache {
		scta = &structCacheTypeArray{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, _ := d.CheckField(rv.Type().Field(i)); ok {
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
			err = d.decode(rv.Field(scta.m[i]), reader)
			if err != nil {
				return err
			}
		} else {
			d.jumpOffset(reader)
		}
	}
	return nil
}

func (d *decoder) setStructFromMap(rv reflect.Value, reader *bufio.Reader, k reflect.Kind) error {
	// get length
	l, err := d.mapLength(reader, k)
	if err != nil {
		return err
	}

	var sctm *structCacheTypeMap
	cache, cacheFind := mapSCTM.Load(rv.Type())
	if !cacheFind {
		sctm = &structCacheTypeMap{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, name := d.CheckField(rv.Type().Field(i)); ok {
				sctm.keys = append(sctm.keys, []byte(name))
				sctm.indexes = append(sctm.indexes, i)
			}
		}
		mapSCTM.Store(rv.Type(), sctm)
	} else {
		sctm = cache.(*structCacheTypeMap)
	}

	for i := 0; i < l; i++ {
		dataKey, err := d.asStringByte(reader, k)
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
			err = d.decode(rv.Field(fieldIndex), reader)
			if err != nil {
				return err
			}
		} else {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *decoder) jumpOffset(reader *bufio.Reader) error {
	code, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch {
	case code == def.True, code == def.False, code == def.Nil:
		// do nothing

	case d.isPositiveFixNum(code) || d.isNegativeFixNum(code):
		// do nothing
	case code == def.Uint8, code == def.Int8:
		err = skipN(reader, def.Byte1)
		if err != nil {
			return err
		}
	case code == def.Uint16, code == def.Int16:
		err = skipN(reader, def.Byte2)
		if err != nil {
			return err
		}
	case code == def.Uint32, code == def.Int32, code == def.Float32:
		err = skipN(reader, def.Byte4)
		if err != nil {
			return err
		}
	case code == def.Uint64, code == def.Int64, code == def.Float64:
		err = skipN(reader, def.Byte8)
		if err != nil {
			return err
		}

	case d.isFixString(code):
		err = skipN(reader, int(code-def.FixStr))
		if err != nil {
			return err
		}
	case code == def.Str8, code == def.Bin8:
		b, err := d.readSize1(reader)
		if err != nil {
			return err
		}

		err = skipN(reader, int(b))
		if err != nil {
			return err
		}
	case code == def.Str16, code == def.Bin16:
		bs, err := d.readSize2(reader)
		if err != nil {
			return err
		}

		err = skipN(reader, int(binary.BigEndian.Uint16(bs)))
		if err != nil {
			return err
		}
	case code == def.Str32, code == def.Bin32:
		bs, err := d.readSize4(reader)
		if err != nil {
			return err
		}

		err = skipN(reader, int(binary.BigEndian.Uint32(bs)))
		if err != nil {
			return err
		}

	case d.isFixSlice(code):
		l := int(code - def.FixArray)
		for i := 0; i < l; i++ {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}
	case code == def.Array16:
		bs, err := d.readSize2(reader)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}
	case code == def.Array32:
		bs, err := d.readSize4(reader)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}

	case d.isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}
	case code == def.Map16:
		bs, err := d.readSize2(reader)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}
	case code == def.Map32:
		bs, err := d.readSize4(reader)
		if err != nil {
			return err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			err = d.jumpOffset(reader)
			if err != nil {
				return err
			}
		}

	case code == def.Fixext1:
		err = skipN(reader, def.Byte1+def.Byte1)
		if err != nil {
			return err
		}
	case code == def.Fixext2:
		err = skipN(reader, def.Byte1+def.Byte2)
		if err != nil {
			return err
		}
	case code == def.Fixext4:
		err = skipN(reader, def.Byte1+def.Byte4)
		if err != nil {
			return err
		}
	case code == def.Fixext8:
		err = skipN(reader, def.Byte1+def.Byte8)
		if err != nil {
			return err
		}
	case code == def.Fixext16:
		err = skipN(reader, def.Byte1+def.Byte16)
		if err != nil {
			return err
		}

	case code == def.Ext8:
		b, err := d.readSize1(reader)
		if err != nil {
			return err
		}

		err = skipN(reader, def.Byte1+int(b))
		if err != nil {
			return err
		}
	case code == def.Ext16:
		bs, err := d.readSize2(reader)
		if err != nil {
			return err
		}

		err = skipN(reader, def.Byte1+int(binary.BigEndian.Uint16(bs)))
		if err != nil {
			return err
		}
	case code == def.Ext32:
		bs, err := d.readSize4(reader)
		if err != nil {
			return err
		}

		err = skipN(reader, def.Byte1+int(binary.BigEndian.Uint32(bs)))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unrecognized code: %d", code)

	}

	return nil
}
