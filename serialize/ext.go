package serialize

import (
	"reflect"
	"time"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) isDateTime(value reflect.Value) (bool, time.Time) {
	i := value.Interface()
	switch t := i.(type) {
	case time.Time:
		return true, t
	}
	return false, now
}

func (s *serializer) calcTime(t time.Time) int {
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			return def.Byte1 + def.Byte4
		}
		return def.Byte1 + def.Byte8
	}

	return def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8
}

func (s *serializer) writeTime(t time.Time, offset int) (int, error) {
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			offset = s.setByte1Int(def.Fixext4, offset)
			offset = s.setByte1Int(def.TimeStamp, offset)
			offset = s.setByte4Uint64(data, offset)
			return offset, nil
		}

		offset = s.setByte1Int(def.Fixext8, offset)
		offset = s.setByte1Int(def.TimeStamp, offset)
		offset = s.setByte8Uint64(data, offset)
		return offset, nil
	}

	offset = s.setByte1Int(def.Ext8, offset)
	offset = s.setByte1Int(12, offset)
	offset = s.setByte1Int(def.TimeStamp, offset)
	offset = s.setByte4Int(t.Nanosecond(), offset)
	offset = s.setByte8Uint64(secs, offset)
	return offset, nil
}
