package decoding

import (
	"encoding/binary"
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

var emptyString = ""
var emptyBytes = []byte{}

func isCodeString(code byte) bool {
	return isFixString(code) || code == def.Str8 || code == def.Str16 || code == def.Str32
}

func isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func stringByteLength(r io.Reader, code byte, k reflect.Kind) (int, error) {
	if def.FixStr <= code && code <= def.FixStr+0x1f {
		l := int(code - def.FixStr)
		return l, nil
	} else if code == def.Str8 {
		b, err := readSize1(r)
		if err != nil {
			return 0, err
		}
		return int(b), nil
	} else if code == def.Str16 {
		b, err := readSize2(r)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(b)), nil
	} else if code == def.Str32 {
		b, err := readSize4(r)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(b)), nil
	} else if code == def.Nil {
		return 0, nil
	}
	return 0, errorTemplate(code, k)
}

func asString(r io.Reader, k reflect.Kind) (string, error) {
	code, err := readSize1(r)
	if err != nil {
		return emptyString, err
	}
	return asStringWithCode(r, code, k)
}

func asStringWithCode(r io.Reader, code byte, k reflect.Kind) (string, error) {
	bs, err := asStringByteWithCode(r, code, k)
	if err != nil {
		return emptyString, err
	}
	return string(bs), nil
}

func asStringByte(r io.Reader, k reflect.Kind) ([]byte, error) {
	code, err := readSize1(r)
	if err != nil {
		return emptyBytes, err
	}
	return asStringByteWithCode(r, code, k)
}

func asStringByteWithCode(r io.Reader, code byte, k reflect.Kind) ([]byte, error) {
	l, err := stringByteLength(r, code, k)
	if err != nil {
		return emptyBytes, err
	}

	return asStringByteByLength(r, l, k)
}

func asStringByteByLength(r io.Reader, l int, _ reflect.Kind) ([]byte, error) {
	if l < 1 {
		return emptyBytes, nil
	}

	return readSizeN(r, l)
}
