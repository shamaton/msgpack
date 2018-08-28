package serialize

func (c *serializer) setByte1Int64(value int64, offset int) int {
	c.d[offset] = byte(value)
	return offset + 1
}

func (c *serializer) setByte2Int64(value int64, offset int) int {
	c.d[offset+0] = byte(value >> 8)
	c.d[offset+1] = byte(value)
	return offset + 2
}

func (c *serializer) setByte4Int64(value int64, offset int) int {
	c.d[offset+0] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *serializer) setByte8Int64(value int64, offset int) int {
	c.d[offset] = byte(value >> 56)
	c.d[offset+1] = byte(value >> 48)
	c.d[offset+2] = byte(value >> 40)
	c.d[offset+3] = byte(value >> 32)
	c.d[offset+4] = byte(value >> 24)
	c.d[offset+5] = byte(value >> 16)
	c.d[offset+6] = byte(value >> 8)
	c.d[offset+7] = byte(value)
	return offset + 8
}

func (c *serializer) setByte1Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value)
	return offset + 1
}

func (c *serializer) setByte2Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value >> 8)
	c.d[offset+1] = byte(value)
	return offset + 2
}

func (c *serializer) setByte4Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *serializer) setByte8Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value >> 56)
	c.d[offset+1] = byte(value >> 48)
	c.d[offset+2] = byte(value >> 40)
	c.d[offset+3] = byte(value >> 32)
	c.d[offset+4] = byte(value >> 24)
	c.d[offset+5] = byte(value >> 16)
	c.d[offset+6] = byte(value >> 8)
	c.d[offset+7] = byte(value)
	return offset + 8
}

func (c *serializer) setByte1Int(code, offset int) int {
	c.d[offset] = byte(code)
	return offset + 1
}

func (c *serializer) setByte2Int(value int, offset int) int {
	c.d[offset] = byte(value >> 8)
	c.d[offset+1] = byte(value)
	return offset + 2
}

func (c *serializer) setByte4Int(value int, offset int) int {
	c.d[offset] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *serializer) setByte4Uint32(value uint32, offset int) int {
	c.d[offset] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *serializer) setBytes(bs []byte, offset int) int {
	for i := range bs {
		c.d[offset+i] = bs[i]
	}
	return offset + len(bs)
}
