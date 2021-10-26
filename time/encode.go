package time

import (
	"io"
	"reflect"
	"time"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
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

func (s *timeEncoder) WriteToBytes(value reflect.Value, writer io.Writer) error {
	t := value.Interface().(time.Time)

	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			err := s.SetByte1Int(def.Fixext4, writer)
			if err != nil {
				return err
			}

			err = s.SetByte1Int(def.TimeStamp, writer)
			if err != nil {
				return err
			}

			err = s.SetByte4Uint64(data, writer)
			if err != nil {
				return err
			}

			return nil
		}

		err := s.SetByte1Int(def.Fixext8, writer)
		if err != nil {
			return err
		}

		err = s.SetByte1Int(def.TimeStamp, writer)
		if err != nil {
			return err
		}

		err = s.SetByte8Uint64(data, writer)
		if err != nil {
			return err
		}

		return nil
	}

	err := s.SetByte1Int(def.Ext8, writer)
	if err != nil {
		return err
	}

	err = s.SetByte1Int(12, writer)
	if err != nil {
		return err
	}

	err = s.SetByte1Int(def.TimeStamp, writer)
	if err != nil {
		return err
	}

	err = s.SetByte4Int(t.Nanosecond(), writer)
	if err != nil {
		return err
	}

	err = s.SetByte8Uint64(secs, writer)
	if err != nil {
		return err
	}

	return nil
}
