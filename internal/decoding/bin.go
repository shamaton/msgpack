package decoding

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) isCodeBin(v byte) bool {
	switch v {
	case def.Bin8, def.Bin16, def.Bin32:
		return true
	}
	return false
}

func (d *decoder) asBin(offset int, k reflect.Kind) ([]byte, int, error) {
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

func (d *decoder) asBinString(offset int, k reflect.Kind) (string, int, error) {
	bs, offset, err := d.asBin(offset, k)
	return *(*string)(unsafe.Pointer(&bs)), offset, err
}
