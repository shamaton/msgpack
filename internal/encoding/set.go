package encoding

func (e *encoder) setByte1Int64(value int64, writer Writer) error {
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte2Int64(value int64, writer Writer) error {
	err := writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte4Int64(value int64, writer Writer) error {
	err := writer.WriteByte(byte(value >> 24))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 16))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte8Int64(value int64, writer Writer) error {
	err := writer.WriteByte(byte(value >> 56))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 48))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 40))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 32))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 24))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 16))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte1Uint64(value uint64, writer Writer) error {
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte2Uint64(value uint64, writer Writer) error {
	err := writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte4Uint64(value uint64, writer Writer) error {
	err := writer.WriteByte(byte(value >> 24))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 16))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte8Uint64(value uint64, writer Writer) error {
	err := writer.WriteByte(byte(value >> 56))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 48))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 40))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 32))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 24))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 16))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte1Int(code int, writer Writer) error {
	return writer.WriteByte(byte(code))
}

func (e *encoder) setByte2Int(value int, writer Writer) error {
	err := writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setByte4Int(value int, writer Writer) error {
	err := writer.WriteByte(byte(value >> 24))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 16))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(value >> 8))
	if err != nil {
		return err
	}
	return writer.WriteByte(byte(value))
}

func (e *encoder) setBytes(bs []byte, writer Writer) error {
	_, err := writer.Write(bs)
	return err
}
