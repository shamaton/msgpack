package deserialize

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (d *deserializer) isFixMap(v byte) bool {
	return def.FixMap <= v && v <= def.FixMap+0x0f
}

func (d *deserializer) mapLength(offset int, k reflect.Kind) (int, int, error) {
	code, offset := d.readSize1(offset)

	switch {
	case d.isFixMap(code):
		return int(code - def.FixMap), offset, nil
	case code == def.Map16:
		bs, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Map32:
		bs, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
}
