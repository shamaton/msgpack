package decoding

import (
	"bufio"
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) asUint(reader *bufio.Reader, k reflect.Kind) (uint64, error) {
	code, err := peekCode(reader)
	if err != nil {
		return 0, err
	}

	switch {
	case d.isPositiveFixNum(code):
		b, err := d.readSize1(reader)
		if err != nil {
			return 0, err
		}
		return uint64(b), nil

	case d.isNegativeFixNum(code):
		b, err := d.readSize1(reader)
		if err != nil {
			return 0, err
		}
		return uint64(int8(b)), nil

	case code == def.Uint8:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		b, err := d.readSize1(reader)
		if err != nil {
			return 0, err
		}
		return uint64(uint8(b)), nil

	case code == def.Int8:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		b, err := d.readSize1(reader)
		if err != nil {
			return 0, err
		}
		return uint64(int8(b)), nil

	case code == def.Uint16:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize2(reader)
		if err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), nil

	case code == def.Int16:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize2(reader)
		if err != nil {
			return 0, err
		}
		v := int16(binary.BigEndian.Uint16(bs))
		return uint64(v), nil

	case code == def.Uint32:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize4(reader)
		if err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), nil

	case code == def.Int32:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize4(reader)
		if err != nil {
			return 0, err
		}
		v := int32(binary.BigEndian.Uint32(bs))
		return uint64(v), nil

	case code == def.Uint64:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize8(reader)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(bs), nil

	case code == def.Int64:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize8(reader)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(bs), nil

	case code == def.Nil:
		return 0, skipOne(reader)
	}

	return 0, d.errorTemplate(code, k)
}
