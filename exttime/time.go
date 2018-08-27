package exttime

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/shamaton/msgpack/def"
	"github.com/shamaton/msgpack/ext"
)

type timeSerializer struct {
	ext.Common
}

// todo : typo
func GetExtSerilizer() ext.ExtSeri {
	return new(timeSerializer)
}

func (s *timeSerializer) IsType(value reflect.Value) bool {
	_, ok := value.Interface().(time.Time)
	return ok
}

func (s *timeSerializer) CalcByteSize(value reflect.Value) (int, error) {
	t := value.Interface().(time.Time)
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			return def.Byte1 + def.Byte4, nil
		}
		return def.Byte1 + def.Byte8, nil
	}

	return def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8, nil
}

func (s *timeSerializer) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	t := value.Interface().(time.Time)

	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			offset = s.SetByte1Int(def.Fixext4, offset, bytes)
			offset = s.SetByte1Int(def.TimeStamp, offset, bytes)
			offset = s.SetByte4Uint64(data, offset, bytes)
			return offset
		}

		offset = s.SetByte1Int(def.Fixext8, offset, bytes)
		offset = s.SetByte1Int(def.TimeStamp, offset, bytes)
		offset = s.SetByte8Uint64(data, offset, bytes)
		return offset
	}

	offset = s.SetByte1Int(def.Ext8, offset, bytes)
	offset = s.SetByte1Int(12, offset, bytes)
	offset = s.SetByte1Int(def.TimeStamp, offset, bytes)
	offset = s.SetByte4Int(t.Nanosecond(), offset, bytes)
	offset = s.SetByte8Uint64(secs, offset, bytes)
	return offset
}

type timeDeserializer struct {
	ext.CommonDeseri
}

// todo : typo
func GetExtDeserilizer() ext.ExtDeseri {
	return new(timeDeserializer)
}

func (td *timeDeserializer) IsType(offset int, d *[]byte) bool {
	code, offset := td.ReadSize1(offset, d)

	if code == def.Fixext4 {
		t, _ := td.ReadSize1(offset, d)
		return int8(t) == def.TimeStamp
	} else if code == def.Fixext8 {
		t, _ := td.ReadSize1(offset, d)
		return int8(t) == def.TimeStamp
	} else if code == def.Ext8 {
		l, offset := td.ReadSize1(offset, d)
		t, _ := td.ReadSize1(offset, d)
		return l == 12 && int8(t) == def.TimeStamp
	}
	return false
}

func (td *timeDeserializer) AsValue(offset int, k reflect.Kind, d *[]byte) (interface{}, int, error) {
	code, offset := td.ReadSize1(offset, d)

	// TODO : In timestamp 64 and timestamp 96 formats, nanoseconds must not be larger than 999999999.

	switch code {
	case def.Fixext4:
		_, offset = td.ReadSize1(offset, d)
		bs, offset := td.ReadSize4(offset, d)
		return time.Unix(int64(binary.BigEndian.Uint32(bs)), 0), offset, nil

	case def.Fixext8:
		_, offset = td.ReadSize1(offset, d)
		bs, offset := td.ReadSize8(offset, d)
		data64 := binary.BigEndian.Uint64(bs)
		return time.Unix(int64(data64&0x00000003ffffffff), int64(data64>>34)), offset, nil

	case def.Ext8:
		_, offset = td.ReadSize1(offset, d)
		_, offset = td.ReadSize1(offset, d)
		nanobs, offset := td.ReadSize4(offset, d)
		secbs, offset := td.ReadSize8(offset, d)
		nano := binary.BigEndian.Uint32(nanobs)
		sec := binary.BigEndian.Uint64(secbs)
		return time.Unix(int64(sec), int64(nano)), offset, nil
	}

	// todo : const now
	return time.Now(), 0, fmt.Errorf("should not reach this line!! code %x decoding %v", code, k)
}
