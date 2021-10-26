package encoding

import (
	"io"
)

func (e *encoder) setByte1Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{byte(value)})
	return err
}

func (e *encoder) setByte2Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (e *encoder) setByte4Int64(value int64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (e *encoder) setByte8Int64(value int64, writer io.Writer) error {
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

func (e *encoder) setByte1Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{byte(value)})
	return err
}

func (e *encoder) setByte2Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (e *encoder) setByte4Uint64(value uint64, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (e *encoder) setByte8Uint64(value uint64, writer io.Writer) error {
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

func (e *encoder) setByte1Int(code int, writer io.Writer) error {
	_, err := writer.Write([]byte{byte(code)})
	return err
}

func (e *encoder) setByte2Int(value int, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (e *encoder) setByte4Int(value int, writer io.Writer) error {
	_, err := writer.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return err
}

func (e *encoder) setBytes(bs []byte, writer io.Writer) error {
	_, err := writer.Write(bs)
	return err
}
