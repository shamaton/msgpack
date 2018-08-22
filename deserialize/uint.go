package deserialize

import (
	"encoding/binary"
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
	} else if code == def.Nil {
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
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
	} else if code == def.Nil {
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
}

func (d *deserializer) asFloat32(offset int, k reflect.Kind) (float32, int, error) {
	code := d.data[offset]
	if code == def.Float32 {
		offset++
		bs, offset := d.readSize4(offset)
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return v, offset, nil
	} else if code == def.Nil {
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
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
	} else if code == def.Nil {
		offset++
		return 0, offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
}

var emptyString = ""
var emptyBytes = []byte{}

func (d *deserializer) asString(offset int, k reflect.Kind) (string, int, error) {
	l, offset, err := d.stringByteLength(offset, k)
	if err != nil {
		return emptyString, 0, err
	}
	bs, offset := d.asStringByte(offset, l, k)
	return *(*string)(unsafe.Pointer(&bs)), offset, nil
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

func (d *deserializer) asStringByte(offset int, l int, k reflect.Kind) ([]byte, int) {
	if l < 1 {
		return emptyBytes, offset
	}

	return d.readSizeN(offset, l)
}

func (d *deserializer) isCodeString(code byte) bool {
	switch {
	case d.isFixString(code), code == def.Str8, code == def.Str16, code == def.Str32:
		return true
	}
	return false
}

func (d *deserializer) asBool(offset int, k reflect.Kind) (bool, int, error) {
	code := d.data[offset]
	offset++

	if code == def.True {
		return true, 0, nil
	} else if code == def.False {
		return false, 0, nil
	}
	return false, 0, d.errorTemplate(code, k)
}

func (d *deserializer) asBin(offset int, k reflect.Kind) ([]byte, int, error) {
	code, offset := d.readSize1(offset)

	switch code {
	case def.Bin8:
		l, offset := d.readSize1(offset)
		o := offset + int(uint8(l))
		return d.data[offset:o], o, nil
	case def.Bin16:
		bs, offset := d.readSize2(offset)
		o := offset + int(binary.BigEndian.Uint16(bs))
		return d.data[offset:o], o, nil
	case def.Bin32:
		bs, offset := d.readSize4(offset)
		o := offset + int(binary.BigEndian.Uint32(bs))
		return d.data[offset:o], o, nil
	}

	return emptyBytes, 0, d.errorTemplate(code, k)
}

func (d *deserializer) sliceLength(offset int, k reflect.Kind) (int, int, error) {
	code, offset := d.readSize1(offset)

	switch {
	case d.isFixSlice(code):
		return int(code - def.FixArray), offset, nil
	case code == def.Array16:
		bs, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Array32:
		bs, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
}

func (d *deserializer) isFixSlice(v byte) bool {
	return def.FixArray <= v && v <= def.FixArray+0x0f
}
