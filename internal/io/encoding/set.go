package encoding

func (e *encoder) setByte1Int64(value int64) error {
	return e.buf.Write(e.w, byte(value))
	//e.buf.B1[0] = byte(value)
	//_, err := e.w.Write(e.buf.B1)
	//return err
}

func (e *encoder) setByte2Int64(value int64) error {
	return e.buf.Write(e.w,
		byte(value>>8),
		byte(value),
	)
	//e.buf.B2[0] = byte(value >> 8)
	//e.buf.B2[1] = byte(value)
	//_, err := e.w.Write(e.buf.B2)
	//return err
}

func (e *encoder) setByte4Int64(value int64) error {
	return e.buf.Write(e.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//e.buf.B4[0] = byte(value >> 24)
	//e.buf.B4[1] = byte(value >> 16)
	//e.buf.B4[2] = byte(value >> 8)
	//e.buf.B4[3] = byte(value)
	//_, err := e.w.Write(e.buf.B4)
	//return err
}

func (e *encoder) setByte8Int64(value int64) error {
	return e.buf.Write(e.w,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//e.buf.B8[0] = byte(value >> 56)
	//e.buf.B8[1] = byte(value >> 48)
	//e.buf.B8[2] = byte(value >> 40)
	//e.buf.B8[3] = byte(value >> 32)
	//e.buf.B8[4] = byte(value >> 24)
	//e.buf.B8[5] = byte(value >> 16)
	//e.buf.B8[6] = byte(value >> 8)
	//e.buf.B8[7] = byte(value)
	//_, err := e.w.Write(e.buf.B8)
	//return err
}

func (e *encoder) setByte1Uint64(value uint64) error {
	return e.buf.Write(e.w, byte(value))
	//e.buf.B1[0] = byte(value)
	//_, err := e.w.Write(e.buf.B1)
	//return err
}

func (e *encoder) setByte2Uint64(value uint64) error {
	return e.buf.Write(e.w,
		byte(value>>8),
		byte(value),
	)
	//e.buf.B2[0] = byte(value >> 8)
	//e.buf.B2[1] = byte(value)
	//_, err := e.w.Write(e.buf.B2)
	//return err
}

func (e *encoder) setByte4Uint64(value uint64) error {
	return e.buf.Write(e.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//e.buf.B4[0] = byte(value >> 24)
	//e.buf.B4[1] = byte(value >> 16)
	//e.buf.B4[2] = byte(value >> 8)
	//e.buf.B4[3] = byte(value)
	//_, err := e.w.Write(e.buf.B4)
	//return err
}

func (e *encoder) setByte8Uint64(value uint64) error {
	return e.buf.Write(e.w,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//e.buf.B8[0] = byte(value >> 56)
	//e.buf.B8[1] = byte(value >> 48)
	//e.buf.B8[2] = byte(value >> 40)
	//e.buf.B8[3] = byte(value >> 32)
	//e.buf.B8[4] = byte(value >> 24)
	//e.buf.B8[5] = byte(value >> 16)
	//e.buf.B8[6] = byte(value >> 8)
	//e.buf.B8[7] = byte(value)
	//_, err := e.w.Write(e.buf.B8)
	//return err
}

func (e *encoder) setByte1Int(value int) error {
	return e.buf.Write(e.w,
		byte(value),
	)
	//e.buf.B1[0] = byte(value)
	//_, err := e.w.Write(e.buf.B1)
	//return err
}

func (e *encoder) setByte2Int(value int) error {
	return e.buf.Write(e.w,
		byte(value>>8),
		byte(value),
	)
	//e.buf.B2[0] = byte(value >> 8)
	//e.buf.B2[1] = byte(value)
	//_, err := e.w.Write(e.buf.B2)
	//return err
}

func (e *encoder) setByte4Int(value int) error {
	return e.buf.Write(e.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
	//e.buf.B4[0] = byte(value >> 24)
	//e.buf.B4[1] = byte(value >> 16)
	//e.buf.B4[2] = byte(value >> 8)
	//e.buf.B4[3] = byte(value)
	//_, err := e.w.Write(e.buf.B4)
	//return err
}

func (e *encoder) setBytes(bs []byte) error {
	return e.buf.Write(e.w, bs...)
	//_, err := e.w.Write(bs)
	//return err
}
