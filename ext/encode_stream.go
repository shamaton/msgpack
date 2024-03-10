package ext

import (
	"github.com/shamaton/msgpack/v2/internal/common"
	"io"
	"reflect"
)

type StreamEncoder interface {
	Code() int8
	Type() reflect.Type
	Write(w io.Writer, value reflect.Value, buf *common.Buffer) error
}

type StreamEncoderCommon struct{}

func (c *StreamEncoderCommon) WriteByte1Int64(w io.Writer, value int64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value),
	)
	//buf.B1[0] = byte(value)
	//_, err := w.Write(buf.B1)
	//return err
}

func (c *StreamEncoderCommon) WriteByte2Int64(w io.Writer, value int64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>8),
		byte(value),
	)
	//buf.B2[0] = byte(value >> 8)
	//buf.B2[1] = byte(value)
	//_, err := w.Write(buf.B2)
	//return err
}

func (c *StreamEncoderCommon) WriteByte4Int64(w io.Writer, value int64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//buf.B4[0] = byte(value >> 24)
	//buf.B4[1] = byte(value >> 16)
	//buf.B4[2] = byte(value >> 8)
	//buf.B4[3] = byte(value)
	//_, err := w.Write(buf.B4)
	//return err
}

func (c *StreamEncoderCommon) WriteByte8Int64(w io.Writer, value int64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//buf.B8[0] = byte(value >> 56)
	//buf.B8[1] = byte(value >> 48)
	//buf.B8[2] = byte(value >> 40)
	//buf.B8[3] = byte(value >> 32)
	//buf.B8[4] = byte(value >> 24)
	//buf.B8[5] = byte(value >> 16)
	//buf.B8[6] = byte(value >> 8)
	//buf.B8[7] = byte(value)
	//_, err := w.Write(buf.B8)
	//return err
}

func (c *StreamEncoderCommon) WriteByte1Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value),
	)
	//buf.B1[0] = byte(value)
	//_, err := w.Write(buf.B1)
	//return err
}

func (c *StreamEncoderCommon) WriteByte2Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>8),
		byte(value),
	)
	//buf.B2[0] = byte(value >> 8)
	//buf.B2[1] = byte(value)
	//_, err := w.Write(buf.B2)
	//return err
}

func (c *StreamEncoderCommon) WriteByte4Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//buf.B4[0] = byte(value >> 24)
	//buf.B4[1] = byte(value >> 16)
	//buf.B4[2] = byte(value >> 8)
	//buf.B4[3] = byte(value)
	//_, err := w.Write(buf.B4)
	//return err
}

func (c *StreamEncoderCommon) WriteByte8Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//buf.B8[0] = byte(value >> 56)
	//buf.B8[1] = byte(value >> 48)
	//buf.B8[2] = byte(value >> 40)
	//buf.B8[3] = byte(value >> 32)
	//buf.B8[4] = byte(value >> 24)
	//buf.B8[5] = byte(value >> 16)
	//buf.B8[6] = byte(value >> 8)
	//buf.B8[7] = byte(value)
	//_, err := w.Write(buf.B8)
	//return err
}

func (c *StreamEncoderCommon) WriteByte1Int(w io.Writer, value int, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value),
	)
	//buf.B1[0] = byte(value)
	//_, err := w.Write(buf.B1)
	//return err
}

func (c *StreamEncoderCommon) WriteByte2Int(w io.Writer, value int, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>8),
		byte(value),
	)
	//buf.B2[0] = byte(value >> 8)
	//buf.B2[1] = byte(value)
	//_, err := w.Write(buf.B2)
	//return err
}

func (c *StreamEncoderCommon) WriteByte4Int(w io.Writer, value int, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//buf.B4[0] = byte(value >> 24)
	//buf.B4[1] = byte(value >> 16)
	//buf.B4[2] = byte(value >> 8)
	//buf.B4[3] = byte(value)
	//_, err := w.Write(buf.B4)
	//return err
}

func (c *StreamEncoderCommon) WriteByte4Uint32(w io.Writer, value uint32, buf *common.Buffer) error {
	return buf.Write(w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//buf.B4[0] = byte(value >> 24)
	//buf.B4[1] = byte(value >> 16)
	//buf.B4[2] = byte(value >> 8)
	//buf.B4[3] = byte(value)
	//_, err := w.Write(buf.B4)
	//return err
}

func (c *StreamEncoderCommon) WriteBytes(w io.Writer, bs []byte, buf *common.Buffer) error {
	return buf.Write(w, bs...)
	//_, err := w.Write(bs)
	//return err
}
