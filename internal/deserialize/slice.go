package deserialize

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (d *deserializer) isFixSlice(v byte) bool {
	return def.FixArray <= v && v <= def.FixArray+0x0f
}

func (d *deserializer) sliceLength(offset int, k reflect.Kind) (int, int, error) {
	code, offset := d.readSize1(offset)

	switch {
	case d.isFixSlice(code):
		return int(code - def.FixArray), offset, nil
	case code == def.Array16:
		bs, offset := d.readSize2(offset)
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Array32:
		bs, offset := d.readSize4(offset)
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}
	return 0, 0, d.errorTemplate(code, k)
}

func (d *deserializer) asFixedSlice(rv reflect.Value, offset int, l int) (int, bool, error) {
	t := rv.Type()
	k := t.Elem().Kind()
	switch t {
	case typeIntSlice:
		sli := make([]int, l)
		for i := range sli {
			v, o, err := d.asInt(offset, k)
			if err != nil {
				return 0, false, err
			}
			sli[i] = int(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(sli))
		return offset, true, nil

	case typeInt8Slice:
		sli := make([]int8, l)
		for i := range sli {
			v, o, err := d.asInt(offset, k)
			if err != nil {
				return 0, false, err
			}
			sli[i] = int8(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(sli))
		return offset, true, nil
	}

	return offset, false, nil
}
