package decoding

import (
	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) readSize1() (byte, error) {
	bs, err := d.readSizeN(def.Byte1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func (d *decoder) readSize2() ([]byte, error) {
	return d.readSizeN(def.Byte2)
}

func (d *decoder) readSize4() ([]byte, error) {
	return d.readSizeN(def.Byte4)
}

func (d *decoder) readSize8() ([]byte, error) {
	return d.readSizeN(def.Byte8)
}

func (d *decoder) readSizeN(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := d.r.Read(b); err != nil {
		return emptyBytes, err
	}
	return b, nil
}
