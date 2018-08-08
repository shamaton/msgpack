package serialize

import "github.com/shamaton/msgpack/def"

type common struct {
	d []byte
}

func (c common) isPositiveFixInt64(v int64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (c common) isPositiveFixUint64(v uint64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (c common) isNegativeFixInt64(v int64) bool {
	return def.NegativeFixintMin <= v && v <= def.NegativeFixintMax
}

func (c *common) writeSize1Int64(value int64, offset int) int {
	c.d[offset] = byte(value)
	return offset + 1
}

func (c *common) writeSize2Int64(value int64, offset int) int {
	c.d[offset+0] = byte(value >> 8)
	c.d[offset+1] = byte(value)
	return offset + 2
}

func (c *common) writeSize4Int64(value int64, offset int) int {
	c.d[offset+0] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *common) writeSize8Int64(value int64, offset int) int {
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

func (c *common) writeSize1Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value)
	return offset + 1
}

func (c *common) writeSize2Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value >> 8)
	c.d[offset+1] = byte(value)
	return offset + 2
}

func (c *common) writeSize4Uint64(value uint64, offset int) int {
	c.d[offset] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *common) writeSize8Uint64(value uint64, offset int) int {
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

func (c *common) writeSize1Int(code, offset int) int {
	c.d[offset] = byte(code)
	return offset + 1
}

func (c *common) writeSize4Int(value int, offset int) int {
	c.d[offset] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}

func (c *common) writeSize4Uint32(value uint32, offset int) int {
	c.d[offset] = byte(value >> 24)
	c.d[offset+1] = byte(value >> 16)
	c.d[offset+2] = byte(value >> 8)
	c.d[offset+3] = byte(value)
	return offset + 4
}
