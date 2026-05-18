package encoding

import "encoding/binary"

func (e *encoder) setUint64Bytes(value uint64, offset int, size int) int {
	switch size {
	case 1:
		e.d[offset] = byte(value) // #nosec G115 -- MessagePack writes the selected low-order byte.
	case 2:
		binary.BigEndian.PutUint16(e.d[offset:], uint16(value)) // #nosec G115 -- MessagePack writes the selected low-order bytes.
	case 4:
		binary.BigEndian.PutUint32(e.d[offset:], uint32(value)) // #nosec G115 -- MessagePack writes the selected low-order bytes.
	case 8:
		binary.BigEndian.PutUint64(e.d[offset:], value)
	default:
		panic("invalid uint64 byte size")
	}
	return offset + size
}

func (e *encoder) setByte1Int64(value int64, offset int) int {
	return e.setUint64Bytes(uint64(value), offset, 1) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte2Int64(value int64, offset int) int {
	return e.setUint64Bytes(uint64(value), offset, 2) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte4Int64(value int64, offset int) int {
	return e.setUint64Bytes(uint64(value), offset, 4) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte8Int64(value int64, offset int) int {
	return e.setUint64Bytes(uint64(value), offset, 8) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte1Uint64(value uint64, offset int) int {
	return e.setUint64Bytes(value, offset, 1)
}

func (e *encoder) setByte2Uint64(value uint64, offset int) int {
	return e.setUint64Bytes(value, offset, 2)
}

func (e *encoder) setByte4Uint64(value uint64, offset int) int {
	return e.setUint64Bytes(value, offset, 4)
}

func (e *encoder) setByte8Uint64(value uint64, offset int) int {
	return e.setUint64Bytes(value, offset, 8)
}

func (e *encoder) setByte1Int(code, offset int) int {
	return e.setUint64Bytes(uint64(code), offset, 1) // #nosec G115 -- callers pass bounded MessagePack code or length values.
}

func (e *encoder) setByte2Int(value int, offset int) int {
	return e.setUint64Bytes(uint64(value), offset, 2) // #nosec G115 -- callers pass bounded MessagePack length values.
}

func (e *encoder) setByte4Int(value int, offset int) int {
	return e.setUint64Bytes(uint64(value), offset, 4) // #nosec G115 -- callers pass bounded MessagePack length values.
}

func (e *encoder) setBytes(bs []byte, offset int) int {
	for i := range bs {
		e.d[offset+i] = bs[i]
	}
	return offset + len(bs)
}
