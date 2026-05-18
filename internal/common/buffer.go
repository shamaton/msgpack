package common

import (
	"encoding/binary"
	"io"
	"sync"
)

type Buffer struct {
	Data   []byte
	B1     []byte
	B2     []byte
	B4     []byte
	B8     []byte
	B16    []byte
	offset int
}

func (b *Buffer) Write(w io.Writer, vs ...byte) error {
	if err := b.ensure(w, len(vs)); err != nil {
		return err
	}
	copy(b.Data[b.offset:], vs)
	b.offset += len(vs)
	return nil
}

func (b *Buffer) WriteString(w io.Writer, s string) error {
	for len(s) > 0 {
		if b.offset == len(b.Data) {
			if err := b.flush(w); err != nil {
				return err
			}
		}
		n := copy(b.Data[b.offset:], s)
		b.offset += n
		s = s[n:]
	}
	return nil
}

func (b *Buffer) WriteUint64(w io.Writer, value uint64, size int) error {
	if err := b.ensure(w, size); err != nil {
		return err
	}
	switch size {
	case 1:
		b.Data[b.offset] = byte(value) // #nosec G115 -- MessagePack writes the selected low-order byte.
	case 2:
		binary.BigEndian.PutUint16(b.Data[b.offset:], uint16(value)) // #nosec G115 -- MessagePack writes the selected low-order bytes.
	case 4:
		binary.BigEndian.PutUint32(b.Data[b.offset:], uint32(value)) // #nosec G115 -- MessagePack writes the selected low-order bytes.
	case 8:
		binary.BigEndian.PutUint64(b.Data[b.offset:], value)
	default:
		panic("invalid uint64 byte size")
	}
	b.offset += size
	return nil
}

func (b *Buffer) Flush(w io.Writer) error {
	_, err := w.Write(b.Data[:b.offset])
	return err
}

func (b *Buffer) ensure(w io.Writer, size int) error {
	if len(b.Data) >= b.offset+size {
		return nil
	}
	if err := b.flush(w); err != nil {
		return err
	}
	if len(b.Data) < size {
		b.Data = append(b.Data, make([]byte, size-len(b.Data))...)
	}
	return nil
}

func (b *Buffer) flush(w io.Writer) error {
	_, err := w.Write(b.Data[:b.offset])
	b.offset = 0
	return err
}

var bufPool = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 64)
		return &Buffer{
			Data: data,
			B1:   data[:1],
			B2:   data[:2],
			B4:   data[:4],
			B8:   data[:8],
			B16:  data[:16],
		}
	},
}

func GetBuffer() *Buffer {
	buf := bufPool.Get().(*Buffer)
	buf.offset = 0
	return buf
}

func PutBuffer(buf *Buffer) {
	bufPool.Put(buf)
}
