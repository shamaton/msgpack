package decoding

import (
	"encoding/binary"
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func asUint(r io.Reader, k reflect.Kind) (uint64, error) {
	code, err := readSize1(r)
	if err != nil {
		return 0, err
	}
	return asUintWithCode(r, code, k)
}

func asUintWithCode(r io.Reader, code byte, k reflect.Kind) (uint64, error) {
	switch {
	case isPositiveFixNum(code):
		return uint64(code), nil

	case isNegativeFixNum(code):
		return uint64(int8(code)), nil

	case code == def.Uint8:
		b, err := readSize1(r)
		if err != nil {
			return 0, err
		}
		return uint64(b), nil

	case code == def.Int8:
		b, err := readSize1(r)
		if err != nil {
			return 0, err
		}
		return uint64(int8(b)), nil

	case code == def.Uint16:
		bs, err := readSize2(r)
		if err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), nil

	case code == def.Int16:
		bs, err := readSize2(r)
		if err != nil {
			return 0, err
		}
		v := int16(binary.BigEndian.Uint16(bs))
		return uint64(v), nil

	case code == def.Uint32:
		bs, err := readSize4(r)
		if err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), nil

	case code == def.Int32:
		bs, err := readSize4(r)
		if err != nil {
			return 0, err
		}
		v := int32(binary.BigEndian.Uint32(bs))
		return uint64(v), nil

	case code == def.Uint64:
		bs, err := readSize8(r)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(bs), nil

	case code == def.Int64:
		bs, err := readSize8(r)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(bs), nil

	case code == def.Nil:
		return 0, nil
	}

	return 0, errorTemplate(code, k)
}
