package ext

import (
	"github.com/shamaton/msgpack/v2/internal/common"
	"io"
	"reflect"
)

type StreamEncoder interface {
	Code() int8
	Type() reflect.Type
	CalcByteSize(value reflect.Value) (int, error)
	WriteToBytes(w io.Writer, value reflect.Value, buf *common.Buffer) error
}

type StreamEncoderCommon struct{}

func (c *StreamEncoderCommon) SetByte1Int64(w io.Writer, value int64, buf *common.Buffer) error {
	buf.B1[0] = byte(value)
	_, err := w.Write(buf.B1)
	return err
}

func (c *StreamEncoderCommon) SetByte2Int64(w io.Writer, value int64, buf *common.Buffer) error {
	buf.B2[0] = byte(value >> 8)
	buf.B2[1] = byte(value)
	_, err := w.Write(buf.B2)
	return err
}

func (c *StreamEncoderCommon) SetByte4Int64(w io.Writer, value int64, buf *common.Buffer) error {
	buf.B4[0] = byte(value >> 24)
	buf.B4[1] = byte(value >> 16)
	buf.B4[2] = byte(value >> 8)
	buf.B4[3] = byte(value)
	_, err := w.Write(buf.B4)
	return err
}

func (c *StreamEncoderCommon) SetByte8Int64(w io.Writer, value int64, buf *common.Buffer) error {
	buf.B8[0] = byte(value >> 56)
	buf.B8[1] = byte(value >> 48)
	buf.B8[2] = byte(value >> 40)
	buf.B8[3] = byte(value >> 32)
	buf.B8[4] = byte(value >> 24)
	buf.B8[5] = byte(value >> 16)
	buf.B8[6] = byte(value >> 8)
	buf.B8[7] = byte(value)
	_, err := w.Write(buf.B8)
	return err
}

func (c *StreamEncoderCommon) SetByte1Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	buf.B1[0] = byte(value)
	_, err := w.Write(buf.B1)
	return err
}

func (c *StreamEncoderCommon) SetByte2Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	buf.B2[0] = byte(value >> 8)
	buf.B2[1] = byte(value)
	_, err := w.Write(buf.B2)
	return err
}

func (c *StreamEncoderCommon) SetByte4Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	buf.B4[0] = byte(value >> 24)
	buf.B4[1] = byte(value >> 16)
	buf.B4[2] = byte(value >> 8)
	buf.B4[3] = byte(value)
	_, err := w.Write(buf.B4)
	return err
}

func (c *StreamEncoderCommon) SetByte8Uint64(w io.Writer, value uint64, buf *common.Buffer) error {
	buf.B8[0] = byte(value >> 56)
	buf.B8[1] = byte(value >> 48)
	buf.B8[2] = byte(value >> 40)
	buf.B8[3] = byte(value >> 32)
	buf.B8[4] = byte(value >> 24)
	buf.B8[5] = byte(value >> 16)
	buf.B8[6] = byte(value >> 8)
	buf.B8[7] = byte(value)
	_, err := w.Write(buf.B8)
	return err
}

func (c *StreamEncoderCommon) SetByte1Int(w io.Writer, value int, buf *common.Buffer) error {
	buf.B1[0] = byte(value)
	_, err := w.Write(buf.B1)
	return err
}

func (c *StreamEncoderCommon) SetByte2Int(w io.Writer, value int, buf *common.Buffer) error {
	buf.B2[0] = byte(value >> 8)
	buf.B2[1] = byte(value)
	_, err := w.Write(buf.B2)
	return err
}

func (c *StreamEncoderCommon) SetByte4Int(w io.Writer, value int, buf *common.Buffer) error {
	buf.B4[0] = byte(value >> 24)
	buf.B4[1] = byte(value >> 16)
	buf.B4[2] = byte(value >> 8)
	buf.B4[3] = byte(value)
	_, err := w.Write(buf.B4)
	return err
}

func (c *StreamEncoderCommon) SetByte4Uint32(w io.Writer, value uint32, buf *common.Buffer) error {
	buf.B4[0] = byte(value >> 24)
	buf.B4[1] = byte(value >> 16)
	buf.B4[2] = byte(value >> 8)
	buf.B4[3] = byte(value)
	_, err := w.Write(buf.B4)
	return err
}

func (c *StreamEncoderCommon) SetBytes(w io.Writer, bs []byte) error {
	_, err := w.Write(bs)
	return err
}
