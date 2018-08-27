package deserialize

import (
	"encoding/binary"
	"reflect"
	"time"

	"github.com/shamaton/msgpack/def"
	"github.com/shamaton/msgpack/ext"
)

var extFuncMaps = map[reflect.Type]ext.ExtDeseri{}
var extFuncs = []ext.ExtDeseri{}

func createCacheFuncs() {
	extFuncs = make([]ext.ExtDeseri, len(extFuncMaps))
	i := 0
	for k := range extFuncMaps {
		extFuncs[i] = extFuncMaps[k]
		i++
	}
}

func SetExtFunc(f ext.ExtDeseri) {
	t := reflect.TypeOf(f)
	_, ok := extFuncMaps[t]
	if !ok {
		extFuncMaps[t] = f
		createCacheFuncs()
	}
}

func UnsetExtFunc(f ext.ExtDeseri) {
	t := reflect.TypeOf(f)
	_, ok := extFuncMaps[t]
	if ok {
		delete(extFuncMaps, t)
		createCacheFuncs()
	}
}

func (d *deserializer) isDateTime(offset int) bool {
	code, offset := d.readSize1(offset)

	if code == def.Fixext4 {
		t, _ := d.readSize1(offset)
		return int8(t) == def.TimeStamp
	} else if code == def.Fixext8 {
		t, _ := d.readSize1(offset)
		return int8(t) == def.TimeStamp
	} else if code == def.Ext8 {
		l, offset := d.readSize1(offset)
		t, _ := d.readSize1(offset)
		return l == 12 && int8(t) == def.TimeStamp
	}
	return false
}

func (d *deserializer) asDateTime(offset int, k reflect.Kind) (time.Time, int, error) {
	code, offset := d.readSize1(offset)

	// TODO : In timestamp 64 and timestamp 96 formats, nanoseconds must not be larger than 999999999.

	switch code {
	case def.Fixext4:
		_, offset = d.readSize1(offset)
		bs, offset := d.readSize4(offset)
		return time.Unix(int64(binary.BigEndian.Uint32(bs)), 0), offset, nil

	case def.Fixext8:
		_, offset = d.readSize1(offset)
		bs, offset := d.readSize8(offset)
		data64 := binary.BigEndian.Uint64(bs)
		return time.Unix(int64(data64&0x00000003ffffffff), int64(data64>>34)), offset, nil

	case def.Ext8:
		_, offset = d.readSize1(offset)
		_, offset = d.readSize1(offset)
		nanobs, offset := d.readSize4(offset)
		secbs, offset := d.readSize8(offset)
		nano := binary.BigEndian.Uint32(nanobs)
		sec := binary.BigEndian.Uint64(secbs)
		return time.Unix(int64(sec), int64(nano)), offset, nil
	}

	return time.Now(), 0, d.errorTemplate(code, k)
}
