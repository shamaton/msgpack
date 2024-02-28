package decoding

import (
	"encoding/binary"
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

var (
	typeIntSlice   = reflect.TypeOf([]int{})
	typeInt8Slice  = reflect.TypeOf([]int8{})
	typeInt16Slice = reflect.TypeOf([]int16{})
	typeInt32Slice = reflect.TypeOf([]int32{})
	typeInt64Slice = reflect.TypeOf([]int64{})

	typeUintSlice   = reflect.TypeOf([]uint{})
	typeUint8Slice  = reflect.TypeOf([]uint8{})
	typeUint16Slice = reflect.TypeOf([]uint16{})
	typeUint32Slice = reflect.TypeOf([]uint32{})
	typeUint64Slice = reflect.TypeOf([]uint64{})

	typeFloat32Slice = reflect.TypeOf([]float32{})
	typeFloat64Slice = reflect.TypeOf([]float64{})

	typeStringSlice = reflect.TypeOf([]string{})

	typeBoolSlice = reflect.TypeOf([]bool{})
)

func isFixSlice(v byte) bool {
	return def.FixArray <= v && v <= def.FixArray+0x0f
}

func sliceLength(r io.Reader, code byte, k reflect.Kind) (int, error) {
	switch {
	case isFixSlice(code):
		return int(code - def.FixArray), nil
	case code == def.Array16:
		bs, err := readSize2(r)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(bs)), nil
	case code == def.Array32:
		bs, err := readSize4(r)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(bs)), nil
	}
	return 0, errorTemplate(code, k)
}

// todo : maybe enable to delete
//func    hasRequiredLeastSliceSize(offset, length int) error {
//	// minimum check (byte length)
//	if len(d.data[offset:]) < length {
//		return errors.New("data length lacks to create map")
//	}
//	return nil
//}

func asFixedSlice(r io.Reader, rv reflect.Value, l int) (bool, error) {
	t := rv.Type()
	k := t.Elem().Kind()

	switch t {
	case typeIntSlice:
		sli := make([]int, l)
		for i := range sli {
			v, err := asInt(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = int(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeUintSlice:
		sli := make([]uint, l)
		for i := range sli {
			v, err := asUint(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = uint(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeStringSlice:
		sli := make([]string, l)
		for i := range sli {
			v, err := asString(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = v
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeBoolSlice:
		sli := make([]bool, l)
		for i := range sli {
			v, err := asBool(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = v
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeFloat32Slice:
		sli := make([]float32, l)
		for i := range sli {
			v, err := asFloat32(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = v
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeFloat64Slice:
		sli := make([]float64, l)
		for i := range sli {
			v, err := asFloat64(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = v
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeInt8Slice:
		sli := make([]int8, l)
		for i := range sli {
			v, err := asInt(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = int8(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeInt16Slice:
		sli := make([]int16, l)
		for i := range sli {
			v, err := asInt(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = int16(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeInt32Slice:
		sli := make([]int32, l)
		for i := range sli {
			v, err := asInt(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = int32(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeInt64Slice:
		sli := make([]int64, l)
		for i := range sli {
			v, err := asInt(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = v
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeUint8Slice:
		sli := make([]uint8, l)
		for i := range sli {
			v, err := asUint(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = uint8(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeUint16Slice:
		sli := make([]uint16, l)
		for i := range sli {
			v, err := asUint(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = uint16(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeUint32Slice:
		sli := make([]uint32, l)
		for i := range sli {
			v, err := asUint(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = uint32(v)
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil

	case typeUint64Slice:
		sli := make([]uint64, l)
		for i := range sli {
			v, err := asUint(r, k)
			if err != nil {
				return false, err
			}
			sli[i] = v
		}
		rv.Set(reflect.ValueOf(sli))
		return true, nil
	}

	return false, nil
}
