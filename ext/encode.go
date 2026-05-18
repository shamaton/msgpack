package ext

import (
	"encoding/binary"
	"reflect"
)

// Encoder defines an interface for encoding values into bytes.
// It provides methods to get the encoding type, calculate the byte size of a value,
// and write the encoded value into a byte slice.
type Encoder interface {
	// Code returns the unique code representing the encoder type.
	Code() int8

	// Type returns the reflect.Type of the value that the encoder handles.
	Type() reflect.Type

	// CalcByteSize calculates the number of bytes required to encode the given value.
	// Returns the size and an error if the calculation fails.
	CalcByteSize(value reflect.Value) (int, error)

	// WriteToBytes encodes the given value into a byte slice starting at the specified offset.
	// Returns the new offset after writing the bytes.
	WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int
}

// EncoderCommon provides utility methods for encoding various types of values into bytes.
// It includes methods to encode integers and unsigned integers of different sizes,
// as well as methods to write raw byte slices into a target byte slice.
type EncoderCommon struct{}

func setUint64Bytes(value uint64, offset int, size int, d *[]byte) int {
	switch size {
	case 1:
		(*d)[offset] = byte(value) // #nosec G115 -- MessagePack writes the selected low-order byte.
	case 2:
		binary.BigEndian.PutUint16((*d)[offset:], uint16(value)) // #nosec G115 -- MessagePack writes the selected low-order bytes.
	case 4:
		binary.BigEndian.PutUint32((*d)[offset:], uint32(value)) // #nosec G115 -- MessagePack writes the selected low-order bytes.
	case 8:
		binary.BigEndian.PutUint64((*d)[offset:], value)
	default:
		panic("invalid uint64 byte size")
	}
	return offset + size
}

// SetByte1Int64 encodes a single byte from the given int64 value into the byte slice at the specified offset.
// Returns the new offset after writing the byte.
func (c *EncoderCommon) SetByte1Int64(value int64, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 1, d) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

// SetByte2Int64 encodes the lower two bytes of the given int64 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte2Int64(value int64, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 2, d) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

// SetByte4Int64 encodes the lower four bytes of the given int64 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte4Int64(value int64, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 4, d) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

// SetByte8Int64 encodes all eight bytes of the given int64 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte8Int64(value int64, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 8, d) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

// SetByte1Uint64 encodes a single byte from the given uint64 value into the byte slice at the specified offset.
// Returns the new offset after writing the byte.
func (c *EncoderCommon) SetByte1Uint64(value uint64, offset int, d *[]byte) int {
	return setUint64Bytes(value, offset, 1, d)
}

// SetByte2Uint64 encodes the lower two bytes of the given uint64 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte2Uint64(value uint64, offset int, d *[]byte) int {
	return setUint64Bytes(value, offset, 2, d)
}

// SetByte4Uint64 encodes the lower four bytes of the given uint64 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte4Uint64(value uint64, offset int, d *[]byte) int {
	return setUint64Bytes(value, offset, 4, d)
}

// SetByte8Uint64 encodes all eight bytes of the given uint64 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte8Uint64(value uint64, offset int, d *[]byte) int {
	return setUint64Bytes(value, offset, 8, d)
}

// SetByte1Int encodes a single byte from the given int value into the byte slice at the specified offset.
// Returns the new offset after writing the byte.
func (c *EncoderCommon) SetByte1Int(code, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(code), offset, 1, d) // #nosec G115 -- callers pass bounded MessagePack code or length values.
}

// SetByte2Int encodes the lower two bytes of the given int value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte2Int(value int, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 2, d) // #nosec G115 -- callers pass bounded MessagePack length values.
}

// SetByte4Int encodes the lower four bytes of the given int value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte4Int(value int, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 4, d) // #nosec G115 -- callers pass bounded MessagePack length values.
}

// SetByte4Uint32 encodes the lower four bytes of the given uint32 value into the byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetByte4Uint32(value uint32, offset int, d *[]byte) int {
	return setUint64Bytes(uint64(value), offset, 4, d)
}

// SetBytes writes the given byte slice `bs` into the target byte slice at the specified offset.
// Returns the new offset after writing the bytes.
func (c *EncoderCommon) SetBytes(bs []byte, offset int, d *[]byte) int {
	for i := range bs {
		(*d)[offset+i] = bs[i]
	}
	return offset + len(bs)
}
