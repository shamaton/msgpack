package decoding

import (
	"github.com/shamaton/msgpack/v2/def"
	"io"
)

func readSize1(r io.Reader) (byte, error) {
	bs, err := readSizeN(r, def.Byte1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func readSize2(r io.Reader) ([]byte, error) {
	return readSizeN(r, def.Byte2)
}

func readSize4(r io.Reader) ([]byte, error) {
	return readSizeN(r, def.Byte4)
}

func readSize8(r io.Reader) ([]byte, error) {
	return readSizeN(r, def.Byte8)
}

func readSizeN(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := r.Read(b); err != nil {
		return emptyBytes, err
	}
	return b, nil
}
