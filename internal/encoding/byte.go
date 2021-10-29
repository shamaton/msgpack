package encoding

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"

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

func (e *encoder) writeByteSliceLength(l int, writer Writer) error {
	if l <= math.MaxUint8 {
		err := e.setByte1Int(def.Bin8, writer)
		if err != nil {
			return err
		}

		return e.setByte1Int(l, writer)
	}

	if l <= math.MaxUint16 {
		err := e.setByte1Int(def.Bin16, writer)
		if err != nil {
			return err
		}

		return e.setByte2Int(l, writer)
	}

	if uint(l) <= math.MaxUint32 {
		err := e.setByte1Int(def.Bin32, writer)
		if err != nil {
			return err
		}

		return e.setByte4Int(l, writer)
	}

	return errors.New("slice too large: " + strconv.FormatInt(int64(l), 10) + " elements")
}
