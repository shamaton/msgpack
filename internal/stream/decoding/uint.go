package decoding

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/internal/common/decodingutil"
)

func (d *decoder) asUint(k reflect.Kind) (uint64, error) {
	code, err := d.readSize1()
	if err != nil {
		return 0, err
	}
	return d.asUintWithCode(code, k)
}

func (d *decoder) asUintWithCode(code byte, k reflect.Kind) (uint64, error) {
	switch {
	case d.isPositiveFixNum(code):
		return uint64(code), nil

	case d.isNegativeFixNum(code):
		return decodingutil.Uint64FromInt64(decodingutil.Int64FromInt8Byte(code), k)

	case code == def.Uint8:
		b, err := d.readSize1()
		if err != nil {
			return 0, err
		}
		return uint64(b), nil

	case code == def.Int8:
		b, err := d.readSize1()
		if err != nil {
			return 0, err
		}
		return decodingutil.Uint64FromInt64(decodingutil.Int64FromInt8Byte(b), k)

	case code == def.Uint16:
		bs, err := d.readSize2()
		if err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), nil

	case code == def.Int16:
		bs, err := d.readSize2()
		if err != nil {
			return 0, err
		}
		return decodingutil.Uint64FromInt64(decodingutil.Int64FromInt16Bits(binary.BigEndian.Uint16(bs)), k)

	case code == def.Uint32:
		bs, err := d.readSize4()
		if err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), nil

	case code == def.Int32:
		bs, err := d.readSize4()
		if err != nil {
			return 0, err
		}
		return decodingutil.Uint64FromInt64(decodingutil.Int64FromInt32Bits(binary.BigEndian.Uint32(bs)), k)

	case code == def.Uint64:
		bs, err := d.readSize8()
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(bs), nil

	case code == def.Int64:
		bs, err := d.readSize8()
		if err != nil {
			return 0, err
		}
		return decodingutil.Uint64FromInt64(decodingutil.Int64FromInt64Bits(binary.BigEndian.Uint64(bs)), k)

	case code == def.Nil:
		return 0, nil
	}

	return 0, d.errorTemplate(code, k)
}
