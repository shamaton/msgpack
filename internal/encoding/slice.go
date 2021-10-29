package encoding

import (
	"errors"
	"math"
	"reflect"
	"strconv"

	"github.com/shamaton/msgpack/v2/def"
)

func (e *encoder) calcFixedSlice(rv reflect.Value) (int, bool) {
	size := 0

	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			size += def.Byte1 + e.calcInt(int64(v))
		}
		return size, true

	case []uint:
		for _, v := range sli {
			size += def.Byte1 + e.calcUint(uint64(v))
		}
		return size, true

	case []string:
		for _, v := range sli {
			size += def.Byte1 + e.calcString(v)
		}
		return size, true

	case []float32:
		for _, v := range sli {
			size += def.Byte1 + e.calcFloat32(float64(v))
		}
		return size, true

	case []float64:
		for _, v := range sli {
			size += def.Byte1 + e.calcFloat64(v)
		}
		return size, true

	case []bool:
		size += def.Byte1 * len(sli)
		return size, true

	case []int8:
		for _, v := range sli {
			size += def.Byte1 + e.calcInt(int64(v))
		}
		return size, true

	case []int16:
		for _, v := range sli {
			size += def.Byte1 + e.calcInt(int64(v))
		}
		return size, true

	case []int32:
		for _, v := range sli {
			size += def.Byte1 + e.calcInt(int64(v))
		}
		return size, true

	case []int64:
		for _, v := range sli {
			size += def.Byte1 + e.calcInt(v)
		}
		return size, true

	case []uint8:
		for _, v := range sli {
			size += def.Byte1 + e.calcUint(uint64(v))
		}
		return size, true

	case []uint16:
		for _, v := range sli {
			size += def.Byte1 + e.calcUint(uint64(v))
		}
		return size, true

	case []uint32:
		for _, v := range sli {
			size += def.Byte1 + e.calcUint(uint64(v))
		}
		return size, true

	case []uint64:
		for _, v := range sli {
			size += def.Byte1 + e.calcUint(v)
		}
		return size, true
	}

	return size, false
}

func (e *encoder) writeSliceLength(l int, writer Writer) (err error) {
	// format size
	if l <= 0x0f {
		return e.setByte1Int(def.FixArray+l, writer)
	} else if l <= math.MaxUint16 {
		err = e.setByte1Int(def.Array16, writer)
		if err != nil {
			return err
		}

		return e.setByte2Int(l, writer)
	} else if uint(l) <= math.MaxUint32 {
		err = e.setByte1Int(def.Array32, writer)
		if err != nil {
			return err
		}

		return e.setByte4Int(l, writer)
	}

	return errors.New("slice too large: " + strconv.FormatInt(int64(l), 10) + " elements")
}

func (e *encoder) writeFixedSlice(rv reflect.Value, writer Writer) (bool, error) {

	switch sli := rv.Interface().(type) {
	case []int:
		for _, v := range sli {
			err := e.writeInt(int64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []uint:
		for _, v := range sli {
			err := e.writeUint(uint64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []string:
		for _, v := range sli {
			err := e.writeString(v, writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []float32:
		for _, v := range sli {
			err := e.writeFloat32(float64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []float64:
		for _, v := range sli {
			err := e.writeFloat64(float64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []bool:
		for _, v := range sli {
			err := e.writeBool(v, writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []int8:
		for _, v := range sli {
			err := e.writeInt(int64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []int16:
		for _, v := range sli {
			err := e.writeInt(int64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []int32:
		for _, v := range sli {
			err := e.writeInt(int64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []int64:
		for _, v := range sli {
			err := e.writeInt(v, writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []uint8:
		for _, v := range sli {
			err := e.writeUint(uint64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []uint16:
		for _, v := range sli {
			err := e.writeUint(uint64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []uint32:
		for _, v := range sli {
			err := e.writeUint(uint64(v), writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil

	case []uint64:
		for _, v := range sli {
			err := e.writeUint(v, writer)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	}

	return false, nil
}
