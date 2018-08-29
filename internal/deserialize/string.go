package deserialize

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/shamaton/msgpack/def"
)

var emptyString = ""
var emptyBytes = []byte{}

func (d *deserializer) isCodeString(code byte) bool {
	// todo : refactor
	switch {
	case d.isFixString(code), code == def.Str8, code == def.Str16, code == def.Str32:
		return true
	}
	return false
}

func (d *deserializer) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *deserializer) stringByteLength(offset int, k reflect.Kind) (int, int, error) {
	code := d.data[offset]
	offset++

	if def.FixStr <= code && code <= def.FixStr+0x1f {
		l := int(code - def.FixStr)
		return l, offset, nil
	} else if code == def.Str8 {
		b, offset := d.readSize1(offset)
		return int(b), offset, nil
	} else if code == def.Str16 {
		b, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(b)), offset, nil
	} else if code == def.Str32 {
		b, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(b)), offset, nil
	} else if code == def.Nil {
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
}

func (d *deserializer) asString(offset int, k reflect.Kind) (string, int, error) {
	l, offset, err := d.stringByteLength(offset, k)
	if err != nil {
		return emptyString, 0, err
	}
	bs, offset := d.asStringByte(offset, l, k)
	return *(*string)(unsafe.Pointer(&bs)), offset, nil
}

func (d *deserializer) asStringByte(offset int, l int, k reflect.Kind) ([]byte, int) {
	if l < 1 {
		return emptyBytes, offset
	}

	return d.readSizeN(offset, l)
}
