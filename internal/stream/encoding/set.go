package encoding

func (e *encoder) setByte1Int64(value int64) error {
	return e.buf.WriteUint64(e.w, uint64(value), 1) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte2Int64(value int64) error {
	return e.buf.WriteUint64(e.w, uint64(value), 2) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte4Int64(value int64) error {
	return e.buf.WriteUint64(e.w, uint64(value), 4) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte8Int64(value int64) error {
	return e.buf.WriteUint64(e.w, uint64(value), 8) // #nosec G115 -- MessagePack encodes signed integers as two's-complement bytes.
}

func (e *encoder) setByte1Uint64(value uint64) error {
	return e.buf.WriteUint64(e.w, value, 1)
}

func (e *encoder) setByte2Uint64(value uint64) error {
	return e.buf.WriteUint64(e.w, value, 2)
}

func (e *encoder) setByte4Uint64(value uint64) error {
	return e.buf.WriteUint64(e.w, value, 4)
}

func (e *encoder) setByte8Uint64(value uint64) error {
	return e.buf.WriteUint64(e.w, value, 8)
}

func (e *encoder) setByte1Int(value int) error {
	return e.buf.WriteUint64(e.w, uint64(value), 1) // #nosec G115 -- callers pass bounded MessagePack code or length values.
}

func (e *encoder) setByte2Int(value int) error {
	return e.buf.WriteUint64(e.w, uint64(value), 2) // #nosec G115 -- callers pass bounded MessagePack length values.
}

func (e *encoder) setByte4Int(value int) error {
	return e.buf.WriteUint64(e.w, uint64(value), 4) // #nosec G115 -- callers pass bounded MessagePack length values.
}

func (e *encoder) setBytes(bs []byte) error {
	return e.buf.Write(e.w, bs...)
}

func (e *encoder) setString(s string) error {
	return e.buf.WriteString(e.w, s)
}
