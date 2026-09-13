package time

import (
	"reflect"
	"time"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

var StreamEncoder = new(timeStreamEncoder)

type timeStreamEncoder struct{}

var _ ext.StreamEncoder = (*timeStreamEncoder)(nil)

func (timeStreamEncoder) Code() int8 {
	return def.TimeStamp
}

func (timeStreamEncoder) Type() reflect.Type {
	return typeOf
}

func (e timeStreamEncoder) Write(w ext.StreamWriter, value reflect.Value) error {
	t := value.Interface().(time.Time)

	sec := t.Unix()
	if sec >= 0 {
		secs := uint64(sec) // #nosec G115 -- non-negative Unix seconds are checked before timestamp64 packing.
		if secs>>34 == 0 {
			data := uint64(t.Nanosecond())<<34 | secs // #nosec G115 -- time.Nanosecond is always in [0, 999999999].
			if data&0xffffffff00000000 == 0 {
				if err := w.WriteByte1Int(def.Fixext4); err != nil {
					return err
				}
				if err := w.WriteByte1Int(def.TimeStamp); err != nil {
					return err
				}
				if err := w.WriteByte4Uint64(data); err != nil {
					return err
				}
				return nil
			}

			if err := w.WriteByte1Int(def.Fixext8); err != nil {
				return err
			}
			if err := w.WriteByte1Int(def.TimeStamp); err != nil {
				return err
			}
			if err := w.WriteByte8Uint64(data); err != nil {
				return err
			}
			return nil
		}
	}

	if err := w.WriteByte1Int(def.Ext8); err != nil {
		return err
	}
	if err := w.WriteByte1Int(12); err != nil {
		return err
	}
	if err := w.WriteByte1Int(def.TimeStamp); err != nil {
		return err
	}
	if err := w.WriteByte4Int(t.Nanosecond()); err != nil {
		return err
	}
	if err := w.WriteByte8Int64(sec); err != nil {
		return err
	}
	return nil
}
