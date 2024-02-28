package decoding

import (
	"encoding/binary"
	"io"
	"reflect"
	"unsafe"

	"github.com/shamaton/msgpack/v2/def"
)

func isCodeBin(v byte) bool {
	switch v {
	case def.Bin8, def.Bin16, def.Bin32:
		return true
	}
	return false
}

func asBin(r io.Reader, k reflect.Kind) ([]byte, error) {
	code, err := readSize1(r)
	if err != nil {
		return emptyBytes, err
	}
	return asBinWithCode(r, code, k)
}

func asBinWithCode(r io.Reader, code byte, k reflect.Kind) ([]byte, error) {

	switch code {
	case def.Bin8:
		l, err := readSize1(r)
		if err != nil {
			return emptyBytes, err
		}
		return readSizeN(r, int(l))

	case def.Bin16:
		bs, err := readSize2(r)
		if err != nil {
			return emptyBytes, err
		}
		return readSizeN(r, int(binary.BigEndian.Uint16(bs)))

	case def.Bin32:
		bs, err := readSize4(r)
		if err != nil {
			return emptyBytes, err
		}
		return readSizeN(r, int(binary.BigEndian.Uint32(bs)))
	}

	return emptyBytes, errorTemplate(code, k)
}

func asBinStringWithCode(r io.Reader, code byte, k reflect.Kind) (string, error) {
	bs, err := asBinWithCode(r, code, k)
	return *(*string)(unsafe.Pointer(&bs)), err
}
