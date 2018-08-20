package deserialize

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (d *deserializer) asUint(offset int, k reflect.Kind) (uint64, int, error) {

	code := d.data[offset]

	if d.isPositiveFixNum(code) {
		b, offset := d.readSize1(offset)
		return uint64(b), offset, nil
	} else if d.isNegativeFixNum(code) {
		b, offset := d.readSize1(offset)
		return uint64(int8(b)), offset, nil
	} else if code == def.Uint8 {
		offset++
		b, offset := d.readSize1(offset)
		return uint64(uint8(b)), offset, nil
	} else if code == def.Int8 {
		offset++
		b, offset := d.readSize1(offset)
		return uint64(int8(b)), offset, nil
	} else if code == def.Uint16 {
		offset++
		bs, offset := d.readSize2(offset)
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), offset, nil
	} else if code == def.Int16 {
		offset++
		bs, offset := d.readSize2(offset)
		v := int16(binary.BigEndian.Uint16(bs))
		return uint64(v), offset, nil
	} else if code == def.Uint32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), offset, nil
	} else if code == def.Int32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := int32(binary.BigEndian.Uint32(bs))
		return uint64(v), offset, nil
	} else if code == def.Uint64 {
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil
	} else if code == def.Int64 {
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil
	} /*else if code == def.Float64 {
		offset++
		bs, offset := d.readSize8(offset)
		v := math.Float64frombits(binary.BigEndian.Uint64(bs))
		return uint64(v), offset, nil
	}*/
	/*else if code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return uint64(v), offset, nil
	} */
	return 0, 0, fmt.Errorf("mismatch code : %x decoing %v", code, k)
}

func (d *deserializer) asUint8(rv reflect.Value, offset int) (uint8, int, error) {
	code := d.data[offset]

	if d.isFixNum(code) {
		b, offset := d.readSize1(offset)
		return uint8(b), offset, nil
	} else if code == def.Uint8 || code == def.Int8 {
		offset++
		b, offset := d.readSize1(offset)
		return uint8(b), offset, nil
	} else if code == def.Uint16 || code == def.Int16 {
		offset++
		bs, offset := d.readSize2(offset)
		return uint8(bs[1]), offset, nil
	} else if code == def.Uint32 || code == def.Int32 || code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		return uint8(bs[3]), offset, nil
	} else if code == def.Uint64 || code == def.Int64 || code == def.Float64 {
		offset++
		bs, offset := d.readSize8(offset)
		return uint8(bs[7]), offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x", code)
}

func (d *deserializer) asUint16(rv reflect.Value, offset int) (uint16, int, error) {
	code := d.data[offset]

	if d.isFixNum(code) {
		b, offset := d.readSize1(offset)
		return uint16(b), offset, nil
	} else if code == def.Uint8 || code == def.Int8 {
		offset++
		b, offset := d.readSize1(offset)
		return uint16(b), offset, nil
	} else if code == def.Uint16 || code == def.Int16 {
		offset++
		bs, offset := d.readSize2(offset)
		return binary.BigEndian.Uint16(bs), offset, nil
	} else if code == def.Uint32 || code == def.Int32 || code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		return binary.BigEndian.Uint16(bs[2:]), offset, nil
	} else if code == def.Uint64 || code == def.Int64 || code == def.Float64 {
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint16(bs[6:]), offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x", code)
}

func (d *deserializer) asUint32(rv reflect.Value, offset int) (uint32, int, error) {
	code := d.data[offset]

	if d.isFixNum(code) {
		b, offset := d.readSize1(offset)
		return uint32(b), offset, nil
	} else if code == def.Uint8 || code == def.Int8 {
		offset++
		b, offset := d.readSize1(offset)
		return uint32(b), offset, nil
	} else if code == def.Uint16 || code == def.Int16 {
		offset++
		bs, offset := d.readSize2(offset)
		return uint32(binary.BigEndian.Uint16(bs)), offset, nil
	} else if code == def.Uint32 || code == def.Int32 || code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		return binary.BigEndian.Uint32(bs), offset, nil
	} else if code == def.Uint64 || code == def.Int64 || code == def.Float64 {
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint32(bs[4:]), offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x", code)
}

func (d *deserializer) asUint64(rv reflect.Value, offset int) (uint64, int, error) {
	code := d.data[offset]

	if d.isFixNum(code) {
		b, offset := d.readSize1(offset)
		return uint64(b), offset, nil
	} else if code == def.Uint8 || code == def.Int8 {
		offset++
		b, offset := d.readSize1(offset)
		return uint64(b), offset, nil
	} else if code == def.Uint16 || code == def.Int16 {
		offset++
		bs, offset := d.readSize2(offset)
		return uint64(binary.BigEndian.Uint16(bs)), offset, nil
	} else if code == def.Uint32 || code == def.Int32 || code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		return uint64(binary.BigEndian.Uint32(bs)), offset, nil
	} else if code == def.Uint64 || code == def.Int64 || code == def.Float64 {
		offset++
		bs, offset := d.readSize8(offset)
		return binary.BigEndian.Uint64(bs), offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x", code)
}
