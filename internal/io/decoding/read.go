package decoding

import (
	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) readSize1() (byte, error) {
	//bs, err := d.readSizeN(def.Byte1)
	//if err != nil {
	//	return 0, err
	//}
	//return bs[0], nil

	//b := d.b[:1]
	//if _, err := d.r.Read(b); err != nil {
	//	return 0, err
	//}
	//return b[0], nil

	//if _, err := d.r.Read(d.b1); err != nil {
	//	return 0, err
	//}
	//return d.b1[0], nil

	if _, err := d.r.Read(d.buf.B1); err != nil {
		return 0, err
	}
	return d.buf.B1[0], nil

	//b, err := d.readSizeNNoCheck(def.Byte1)
	//return b[0], err
}

func (d *decoder) readSize2() ([]byte, error) {
	if _, err := d.r.Read(d.buf.B2); err != nil {
		return emptyBytes, err
	}
	return d.buf.B2, nil

	//return d.readSizeNNoCheck(d.b2)

	//return d.readSizeNNoCheck(def.Byte2)
}

func (d *decoder) readSize4() ([]byte, error) {
	if _, err := d.r.Read(d.buf.B4); err != nil {
		return emptyBytes, err
	}
	return d.buf.B4, nil

	//return d.readSizeNNoCheck(d.b4)

	//return d.readSizeNNoCheck(def.Byte4)
}

func (d *decoder) readSize8() ([]byte, error) {
	if _, err := d.r.Read(d.buf.B8); err != nil {
		return emptyBytes, err
	}
	return d.buf.B8, nil

	//return d.readSizeNNoCheck(d.b8)

	//return d.readSizeNNoCheck(def.Byte8)
}

func (d *decoder) readSize16() ([]byte, error) {
	if _, err := d.r.Read(d.buf.B16); err != nil {
		return emptyBytes, err
	}
	return d.buf.B16, nil

	//return d.readSizeNNoCheck(d.b16)

	//return d.readSizeNNoCheck(def.Byte16)
}

func (d *decoder) readSizeNNoCheck(n int) ([]byte, error) {
	b := d.buf.Data[:n]
	if _, err := d.r.Read(b); err != nil {
		return emptyBytes, err
	}
	return b, nil
}

func (d *decoder) readSizeN(n int) ([]byte, error) {
	var b []byte
	if n <= def.Byte32 {
		b = d.buf.Data[:n]
	} else {
		b = make([]byte, n)
	}
	if _, err := d.r.Read(b); err != nil {
		return emptyBytes, err
	}
	return b, nil
}
