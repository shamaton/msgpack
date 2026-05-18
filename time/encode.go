package time

import (
	"reflect"
	"time"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

var Encoder = new(timeEncoder)

type timeEncoder struct {
	ext.EncoderCommon
}

var typeOf = reflect.TypeOf(time.Time{})

func (td *timeEncoder) Code() int8 {
	return def.TimeStamp
}

func (s *timeEncoder) Type() reflect.Type {
	return typeOf
}

func (s *timeEncoder) CalcByteSize(value reflect.Value) (int, error) {
	t := value.Interface().(time.Time)
	sec := t.Unix()
	if sec >= 0 {
		secs := uint64(sec) // #nosec G115 -- non-negative Unix seconds are checked before timestamp64 packing.
		if secs>>34 == 0 {
			data := uint64(t.Nanosecond())<<34 | secs // #nosec G115 -- time.Nanosecond is always in [0, 999999999].
			if data&0xffffffff00000000 == 0 {
				return def.Byte1 + def.Byte1 + def.Byte4, nil
			}
			return def.Byte1 + def.Byte1 + def.Byte8, nil
		}
	}

	return def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8, nil
}

func (s *timeEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	t := value.Interface().(time.Time)

	sec := t.Unix()
	if sec >= 0 {
		secs := uint64(sec) // #nosec G115 -- non-negative Unix seconds are checked before timestamp64 packing.
		if secs>>34 == 0 {
			data := uint64(t.Nanosecond())<<34 | secs // #nosec G115 -- time.Nanosecond is always in [0, 999999999].
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
	}

	offset = s.SetByte1Int(def.Ext8, offset, bytes)
	offset = s.SetByte1Int(12, offset, bytes)
	offset = s.SetByte1Int(def.TimeStamp, offset, bytes)
	offset = s.SetByte4Int(t.Nanosecond(), offset, bytes)
	offset = s.SetByte8Int64(sec, offset, bytes)
	return offset
}
