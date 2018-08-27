package ext

import (
	"reflect"

	"github.com/shamaton/msgpack/def"
)

type ExtSeri interface {
	IsType(value reflect.Value) bool
	CalcByteSize(value reflect.Value) (int, error)
	WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int
}

type ExtDeseri interface {
	IsType(offset int, d *[]byte) bool
	AsValue(offset int, k reflect.Kind, d *[]byte) (interface{}, int, error)
}

type CommonDeseri struct {
}

func (cd *CommonDeseri) ReadSize1(index int, d *[]byte) (byte, int) {
	rb := def.Byte1
	return (*d)[index], index + rb
}

func (cd *CommonDeseri) ReadSize2(index int, d *[]byte) ([]byte, int) {
	rb := def.Byte2
	return (*d)[index : index+rb], index + rb
}

func (cd *CommonDeseri) ReadSize4(index int, d *[]byte) ([]byte, int) {
	rb := def.Byte4
	return (*d)[index : index+rb], index + rb
}

func (cd *CommonDeseri) ReadSize8(index int, d *[]byte) ([]byte, int) {
	rb := def.Byte8
	return (*d)[index : index+rb], index + rb
}

func (cd *CommonDeseri) ReadSizeN(index, n int, d *[]byte) ([]byte, int) {
	return (*d)[index : index+n], index + n
}

type Common struct {
}

func (c *Common) SetByte1Int64(value int64, offset int, d *[]byte) int {
	(*d)[offset] = byte(value)
	return offset + 1
}

func (c *Common) SetByte2Int64(value int64, offset int, d *[]byte) int {
	(*d)[offset+0] = byte(value >> 8)
	(*d)[offset+1] = byte(value)
	return offset + 2
}

func (c *Common) SetByte4Int64(value int64, offset int, d *[]byte) int {
	(*d)[offset+0] = byte(value >> 24)
	(*d)[offset+1] = byte(value >> 16)
	(*d)[offset+2] = byte(value >> 8)
	(*d)[offset+3] = byte(value)
	return offset + 4
}

func (c *Common) SetByte8Int64(value int64, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 56)
	(*d)[offset+1] = byte(value >> 48)
	(*d)[offset+2] = byte(value >> 40)
	(*d)[offset+3] = byte(value >> 32)
	(*d)[offset+4] = byte(value >> 24)
	(*d)[offset+5] = byte(value >> 16)
	(*d)[offset+6] = byte(value >> 8)
	(*d)[offset+7] = byte(value)
	return offset + 8
}

func (c *Common) SetByte1Uint64(value uint64, offset int, d *[]byte) int {
	(*d)[offset] = byte(value)
	return offset + 1
}

func (c *Common) SetByte2Uint64(value uint64, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 8)
	(*d)[offset+1] = byte(value)
	return offset + 2
}

func (c *Common) SetByte4Uint64(value uint64, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 24)
	(*d)[offset+1] = byte(value >> 16)
	(*d)[offset+2] = byte(value >> 8)
	(*d)[offset+3] = byte(value)
	return offset + 4
}

func (c *Common) SetByte8Uint64(value uint64, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 56)
	(*d)[offset+1] = byte(value >> 48)
	(*d)[offset+2] = byte(value >> 40)
	(*d)[offset+3] = byte(value >> 32)
	(*d)[offset+4] = byte(value >> 24)
	(*d)[offset+5] = byte(value >> 16)
	(*d)[offset+6] = byte(value >> 8)
	(*d)[offset+7] = byte(value)
	return offset + 8
}

func (c *Common) SetByte1Int(code, offset int, d *[]byte) int {
	(*d)[offset] = byte(code)
	return offset + 1
}

func (c *Common) SetByte2Int(value int, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 8)
	(*d)[offset+1] = byte(value)
	return offset + 2
}

func (c *Common) SetByte4Int(value int, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 24)
	(*d)[offset+1] = byte(value >> 16)
	(*d)[offset+2] = byte(value >> 8)
	(*d)[offset+3] = byte(value)
	return offset + 4
}

func (c *Common) SetByte4Uint32(value uint32, offset int, d *[]byte) int {
	(*d)[offset] = byte(value >> 24)
	(*d)[offset+1] = byte(value >> 16)
	(*d)[offset+2] = byte(value >> 8)
	(*d)[offset+3] = byte(value)
	return offset + 4
}

func (c *Common) SetBytes(bs []byte, offset int, d *[]byte) int {
	for i := range bs {
		(*d)[offset+i] = bs[i]
	}
	return offset + len(bs)
}
