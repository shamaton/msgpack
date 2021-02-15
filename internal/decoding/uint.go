package decoding

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) asUint(offset int, k reflect.Kind) (uint64, int, error) {

	code := d.data[offset]

	switch {
	case d.isPositiveFixNum(code):
		b, offset := d.readSize1(offset)
		return uint64(b), offset, nil

	case d.isNegativeFixNum(code):
		b, offset := d.readSize1(offset)
		return uint64(int8(b)), offset, nil

	case code == def.Uint8:
		offset++
		b, offset := d.readSize1(offset)
		return uint64(uint8(b)), offset, nil

	case code == def.Int8:
		offset++
		b, offset := d.readSize1(offset)
		return uint64(int8(b)), offset, nil

	case code == def.Uint16:
		offset++
		bs, offset := d.readSize2(offset)
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), offset, nil

	case code == def.Int16:
		offset++
		bs, offset := d.readSize2(offset)
		v := int16(binary.BigEndian.Uint16(bs))
		return uint64(v), offset, nil

	case code == def.Uint32:
		offset++
		bs, offset := d.readSize4(offset)
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), offset, nil

	case code == def.Int32:
		offset++
		bs, offset := d.readSize4(offset)
		v := int32(binary.BigEndian.Uint32(bs))
		return uint64(v), offset, nil

	case code == def.Uint64:
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil

	case code == def.Int64:
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil

	case code == def.Nil:
		offset++
		return 0, offset, nil
	}

	return 0, 0, d.errorTemplate(code, k)
}
