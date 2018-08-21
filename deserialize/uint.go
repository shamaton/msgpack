package deserialize

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"unsafe"

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
	}
	return 0, 0, fmt.Errorf("mismatch code : %x decoing %v", code, k)
}

func (d *deserializer) asInt(offset int, k reflect.Kind) (int64, int, error) {

	code := d.data[offset]

	if d.isPositiveFixNum(code) {
		b, offset := d.readSize1(offset)
		return int64(b), offset, nil
	} else if d.isNegativeFixNum(code) {
		b, offset := d.readSize1(offset)
		return int64(int8(b)), offset, nil
	} else if code == def.Uint8 {
		offset++
		b, offset := d.readSize1(offset)
		return int64(uint8(b)), offset, nil
	} else if code == def.Int8 {
		offset++
		b, offset := d.readSize1(offset)
		return int64(int8(b)), offset, nil
	} else if code == def.Uint16 {
		offset++
		bs, offset := d.readSize2(offset)
		v := binary.BigEndian.Uint16(bs)
		return int64(v), offset, nil
	} else if code == def.Int16 {
		offset++
		bs, offset := d.readSize2(offset)
		v := int16(binary.BigEndian.Uint16(bs))
		return int64(v), offset, nil
	} else if code == def.Uint32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := binary.BigEndian.Uint32(bs)
		return int64(v), offset, nil
	} else if code == def.Int32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := int32(binary.BigEndian.Uint32(bs))
		return int64(v), offset, nil
	} else if code == def.Uint64 {
		offset++
		bs, offset := d.readSize8(offset)
		return int64(binary.BigEndian.Uint64(bs)), offset, nil
	} else if code == def.Int64 {
		offset++
		bs, offset := d.readSize8(offset)
		return int64(binary.BigEndian.Uint64(bs)), offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x decoing %v", code, k)
}

func (d *deserializer) asFloat32(offset int, k reflect.Kind) (float32, int, error) {
	code := d.data[offset]
	if code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return v, offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x decoing %v", code, k)
}

func (d *deserializer) asFloat64(offset int, k reflect.Kind) (float64, int, error) {
	code := d.data[offset]
	if code == def.Float64 {
		offset++
		bs, offset := d.readSize8(offset)
		v := math.Float64frombits(binary.BigEndian.Uint64(bs))
		return v, offset, nil
	} else if code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return float64(v), offset, nil
	}
	return 0, 0, fmt.Errorf("mismatch code : %x decoing %v", code, k)
}

var emptyString = ""

func (d *deserializer) asString(offset int, k reflect.Kind) (string, int, error) {
	code := d.data[offset]
	offset++

	if def.FixStr <= code && code <= def.FixStr+0x1f {
		l := int(code - def.FixStr)
		bs, offset := d.readSizeN(offset, l)
		return *(*string)(unsafe.Pointer(&bs)), offset, nil
	} else if code == def.Str8 {
		b, offset := d.readSize1(offset)
		bs, offset := d.readSizeN(offset, int(b))
		return *(*string)(unsafe.Pointer(&bs)), offset, nil
	} else if code == def.Str16 {
		b, offset := d.readSize2(offset)
		bs, offset := d.readSizeN(offset, int(binary.BigEndian.Uint16(b)))
		return *(*string)(unsafe.Pointer(&bs)), offset, nil

	} else if code == def.Str32 {
		b, offset := d.readSize4(offset)
		bs, offset := d.readSizeN(offset, int(binary.BigEndian.Uint32(b)))
		return *(*string)(unsafe.Pointer(&bs)), offset, nil
	}
	return emptyString, 0, fmt.Errorf("mismatch code : %x decoing %v", code, k)
}
