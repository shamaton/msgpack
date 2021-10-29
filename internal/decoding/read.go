package decoding

import (
	"bufio"
	"errors"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) readSize1(reader *bufio.Reader) (byte, error) {
	return reader.ReadByte()
}

func (d *decoder) readSize2(reader *bufio.Reader) (p [def.Byte2]byte, err error) {
	return p, readFull(reader, p[:])
}

func (d *decoder) readSize4(reader *bufio.Reader) (p [def.Byte4]byte, err error) {
	return p, readFull(reader, p[:])
}

func (d *decoder) readSize8(reader *bufio.Reader) (p [def.Byte8]byte, err error) {
	return p, readFull(reader, p[:])
}

func (d *decoder) readSizeN(reader *bufio.Reader, n int) (p []byte, err error) {
	p = make([]byte, n)
	return p, readFull(reader, p)
}

func (d *decoder) readSizeNBuf(reader *bufio.Reader, buf []byte, n int) ([]byte, error) {
	if n > len(buf) {
		return d.readSizeN(reader, n)
	}

	buf = buf[:n]
	return buf, readFull(reader, buf)
}

func readFull(reader *bufio.Reader, buf []byte) (err error) {
	for i := 0; i < len(buf); {
		b, err := reader.Peek(len(buf) - i)
		if err != nil && !errors.Is(err, bufio.ErrBufferFull) {
			return err
		}

		copy(buf[i:], b)

		i += len(b)
		_, err = reader.Discard(len(b))
		if err != nil {
			return err
		}
	}

	return nil
}
