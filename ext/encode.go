package ext

import (
	"io"
	"reflect"
)

type Encoder interface {
	Code() int8
	Type() reflect.Type
	CalcByteSize(value reflect.Value) (int, error)
	WriteToBytes(value reflect.Value, writer io.Writer) error
}

type EncoderCommon struct {
}

func (c *EncoderCommon) SetByte1Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{byte(value)})
	return err
}

func (c *EncoderCommon) SetByte2Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte4Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte8Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte1Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{byte(value)})
	return err
}

func (c *EncoderCommon) SetByte2Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte4Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte8Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte1Int(code int, writer io.Writer) error {
	_, err := writer.Write([]byte{byte(code)})
	return err
}

func (c *EncoderCommon) SetByte2Int(value int, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte4Int(value int, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetByte4Uint32(value uint32, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (c *EncoderCommon) SetBytes(bs []byte, writer io.Writer) error {
	_, err := writer.Write(bs)
	return err
}
