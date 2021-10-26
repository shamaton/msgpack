package decoding

import (
	"bufio"
	"io"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) readSize1(reader *bufio.Reader) (byte, error) {
	return reader.ReadByte()
}

func (d *decoder) readSize2(reader *bufio.Reader) ([]byte, error) {
	p := make([]byte, def.Byte2)
	_, err := io.ReadFull(reader, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *decoder) readSize4(reader *bufio.Reader) ([]byte, error) {
	p := make([]byte, def.Byte4)
	_, err := io.ReadFull(reader, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *decoder) readSize8(reader *bufio.Reader) ([]byte, error) {
	p := make([]byte, def.Byte8)
	_, err := io.ReadFull(reader, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *decoder) readSizeN(reader *bufio.Reader, n int) ([]byte, error) {
	p := make([]byte, n)
	_, err := io.ReadFull(reader, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
