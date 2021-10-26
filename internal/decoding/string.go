package decoding

import (
	"bufio"
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

var emptyString = ""
var emptyBytes = []byte{}

func (d *decoder) isCodeString(code byte) bool {
	return d.isFixString(code) || code == def.Str8 || code == def.Str16 || code == def.Str32
}

func (d *decoder) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *decoder) stringByteLength(reader *bufio.Reader, k reflect.Kind) (int, error) {
	code, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if def.FixStr <= code && code <= def.FixStr+0x1f {
		l := int(code - def.FixStr)
		return l, nil
	} else if code == def.Str8 {
		b, err := d.readSize1(reader)
		if err != nil {
			return 0, err
		}

		return int(b), nil
	} else if code == def.Str16 {
		b, err := d.readSize2(reader)
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint16(b)), nil
	} else if code == def.Str32 {
		b, err := d.readSize4(reader)
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint32(b)), nil
	} else if code == def.Nil {
		return 0, nil
	}
	return 0, d.errorTemplate(code, k)
}

func (d *decoder) asString(reader *bufio.Reader, k reflect.Kind) (string, error) {
	bs, err := d.asStringByte(reader, k)
	if err != nil {
		return emptyString, err
	}
	return string(bs), nil
}

func (d *decoder) asStringByte(reader *bufio.Reader, k reflect.Kind) ([]byte, error) {
	l, err := d.stringByteLength(reader, k)
	if err != nil {
		return emptyBytes, err
	}

	return d.asStringByteByLength(reader, l, k)
}

func (d *decoder) asStringByteByLength(reader *bufio.Reader, l int, k reflect.Kind) ([]byte, error) {
	if l < 1 {
		return emptyBytes, nil
	}

	return d.readSizeN(reader, l)
}
