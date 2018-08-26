package deserialize

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

type structCache struct {
	m map[string]int
}

type structCache2 struct {
	m []int
}

// todo : rename
var cachemap = map[reflect.Type]*structCache{}
var cachemap2 = map[reflect.Type]*structCache2{}

// todo : change method name
func (d *deserializer) jumpByte(offset int) int {
	code, offset := d.readSize1(offset)
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
		b, offset := d.readSize1(offset)
		offset += int(b)
	case code == def.Str16, code == def.Bin16:
		bs, offset := d.readSize2(offset)
		offset += int(binary.BigEndian.Uint16(bs))
	case code == def.Str32, code == def.Bin32:
		bs, offset := d.readSize4(offset)
		offset += int(binary.BigEndian.Uint32(bs))

	case d.isFixSlice(code):
		l := int(code - def.FixStr)
		for i := 0; i < l; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Array16:
		bs, offset := d.readSize2(offset)
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Array32:
		bs, offset := d.readSize4(offset)
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			offset += d.jumpByte(offset)
		}

	case d.isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Map16:
		bs, offset := d.readSize2(offset)
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Map32:
		bs, offset := d.readSize4(offset)
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			offset += d.jumpByte(offset)
		}

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
		b, offset := d.readSize1(offset)
		offset += def.Byte1 + int(b)
	case code == def.Ext16:
		bs, offset := d.readSize2(offset)
		offset += def.Byte1 + int(binary.BigEndian.Uint16(bs))
	case code == def.Ext32:
		bs, offset := d.readSize4(offset)
		offset += def.Byte1 + int(binary.BigEndian.Uint32(bs))

	}
	return offset
}

// todo same method...
func (d *deserializer) checkField(field reflect.StructField) (bool, string) {
	// A to Z
	if d.isPublic(field.Name) {
		if tag := field.Tag.Get("msgpack"); tag == "ignore" {
			return false, ""
		} else if len(tag) > 0 {
			return true, tag
		}
		return true, field.Name
	}
	return false, ""
}

// todo same method...
func (d *deserializer) isPublic(name string) bool {
	return 0x41 <= name[0] && name[0] <= 0x5a
}
