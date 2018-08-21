package deserialize

import (
	"github.com/shamaton/msgpack/def"
)

func (d *deserializer) readSize1(index int) (byte, int) {
	rb := def.Byte1
	return d.data[index], index + rb
}

func (d *deserializer) readSize2(index int) ([]byte, int) {
	rb := def.Byte2
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) readSize4(index int) ([]byte, int) {
	rb := def.Byte4
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) readSize8(index int) ([]byte, int) {
	rb := def.Byte8
	return d.data[index : index+rb], index + rb
}

// TODO uint
func (d *deserializer) readSizeN(index, n int) ([]byte, int) {
	return d.data[index : index+n], index + n
}
