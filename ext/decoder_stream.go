package ext

import (
	"github.com/shamaton/msgpack/v2/def"
	"io"
	"reflect"
)

var emptyBytes []byte

type StreamDecoder interface {
	Code() int8
	IsType(code byte, innerType int8, dataLength int) bool
	AsValue(code byte, data []byte, k reflect.Kind) (interface{}, error)
}

type DecoderStreamCommon struct {
}

func (d *DecoderStreamCommon) ReadSize1(r io.Reader) (byte, error) {
	bs, err := d.ReadSizeN(r, def.Byte1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func (d *DecoderStreamCommon) ReadSize2(r io.Reader) ([]byte, error) {
	return d.ReadSizeN(r, def.Byte2)
}

func (d *DecoderStreamCommon) ReadSize4(r io.Reader) ([]byte, error) {
	return d.ReadSizeN(r, def.Byte4)
}

func (d *DecoderStreamCommon) ReadSize8(r io.Reader) ([]byte, error) {
	return d.ReadSizeN(r, def.Byte8)
}

func (d *DecoderStreamCommon) ReadSizeN(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := r.Read(b); err != nil {
		return emptyBytes, err
	}
	return b, nil
}
