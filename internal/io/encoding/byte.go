package encoding

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

var typeByte = reflect.TypeOf(byte(0))

func (e *encoder) isByteSlice(rv reflect.Value) bool {
	return rv.Type().Elem() == typeByte
}

func (e *encoder) calcByteSlice(l int) (int, error) {
	if l <= math.MaxUint8 {
		return def.Byte1 + l, nil
	} else if l <= math.MaxUint16 {
		return def.Byte2 + l, nil
	} else if uint(l) <= math.MaxUint32 {
		return def.Byte4 + l, nil
	}
	// not supported error
	return 0, fmt.Errorf("not support this array length : %d", l)
}

func (e *encoder) writeByteSliceLength(l int) error {
	if l <= math.MaxUint8 {
		if err := e.setByte1Int(def.Bin8); err != nil {
			return err
		}
		if err := e.setByte1Int(l); err != nil {
			return err
		}
	} else if l <= math.MaxUint16 {
		if err := e.setByte1Int(def.Bin16); err != nil {
			return err
		}
		if err := e.setByte2Int(l); err != nil {
			return err
		}
	} else if uint(l) <= math.MaxUint32 {
		if err := e.setByte1Int(def.Bin32); err != nil {
			return err
		}
		if err := e.setByte4Int(l); err != nil {
			return err
		}
	}
	return nil
}
