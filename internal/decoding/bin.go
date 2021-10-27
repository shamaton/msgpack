package decoding

import (
	"bufio"
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) isCodeBin(v byte) bool {
	switch v {
	case def.Bin8, def.Bin16, def.Bin32:
		return true
	}
	return false
}

func (d *decoder) asBin(reader *bufio.Reader, k reflect.Kind) ([]byte, error) {
	code, err := d.readSize1(reader)
	if err != nil {
		return nil, err
	}

	switch code {
	case def.Bin8:
		l, err := d.readSize1(reader)
		if err != nil {
			return nil, err
		}

		return d.readSizeN(reader, int(l))
	case def.Bin16:
		bs, err := d.readSize2(reader)
		if err != nil {
			return nil, err
		}

		return d.readSizeN(reader, int(binary.BigEndian.Uint16(bs[:])))
	case def.Bin32:
		bs, err := d.readSize4(reader)
		if err != nil {
			return nil, err
		}

		return d.readSizeN(reader, int(binary.BigEndian.Uint32(bs[:])))
	}

	return emptyBytes, d.errorTemplate(code, k)
}

func (d *decoder) asBinString(reader *bufio.Reader, k reflect.Kind) (string, error) {
	bs, err := d.asBin(reader, k)
	return string(bs), err
}
