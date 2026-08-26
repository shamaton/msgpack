package decoding

import (
	"fmt"
	"io"
)

// The readSize* helpers use io.ReadFull so that readers which legitimately
// return fewer bytes than requested (n > 0, err == nil) are handled: the
// common case still completes with a single Read call, and only short reads
// pay for the completion loop.

func (d *decoder) readSize1() (byte, error) {
	if _, err := io.ReadFull(d.r, d.buf.B1); err != nil {
		return 0, err
	}
	return d.buf.B1[0], nil
}

func (d *decoder) readSize2() ([]byte, error) {
	if _, err := io.ReadFull(d.r, d.buf.B2); err != nil {
		return emptyBytes, err
	}
	return d.buf.B2, nil
}

func (d *decoder) readSize4() ([]byte, error) {
	if _, err := io.ReadFull(d.r, d.buf.B4); err != nil {
		return emptyBytes, err
	}
	return d.buf.B4, nil
}

func (d *decoder) readSize8() ([]byte, error) {
	if _, err := io.ReadFull(d.r, d.buf.B8); err != nil {
		return emptyBytes, err
	}
	return d.buf.B8, nil
}

func (d *decoder) readSize16() ([]byte, error) {
	if _, err := io.ReadFull(d.r, d.buf.B16); err != nil {
		return emptyBytes, err
	}
	return d.buf.B16, nil
}

func (d *decoder) readSizeN(n int) ([]byte, error) {
	if n < 0 {
		return emptyBytes, fmt.Errorf("invalid declared byte length %d", n)
	}
	if n <= len(d.buf.Data) {
		b := d.buf.Data[:n]
		if _, err := io.ReadFull(d.r, b); err != nil {
			return emptyBytes, err
		}
		return b, nil
	}
	if n <= maxPreallocBytes {
		b := make([]byte, n)
		if _, err := io.ReadFull(d.r, b); err != nil {
			return emptyBytes, err
		}
		return b, nil
	}
	return d.readSizeNGrowing(n)
}

func (d *decoder) readSizeNGrowing(n int) ([]byte, error) {
	b := make([]byte, 0, initialByteCap(n))
	chunk := make([]byte, min(readChunkSize, initialByteCap(n)))
	for len(b) < n {
		c := chunk
		if remaining := n - len(b); remaining < len(c) {
			c = c[:remaining]
		}
		if _, err := io.ReadFull(d.r, c); err != nil {
			return emptyBytes, err
		}
		b = append(b, c...)
	}
	return b, nil
}
