package serialize

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (s *serializer) calcStructArray(rv reflect.Value) (int, error) {
	ret := 0
	num := 0
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.CanSet() {
			size, err := s.calcSize(rv.Field(i))
			if err != nil {
				return 0, err
			}
			ret += size
			num++
		}
	}

	// format size
	if num <= 0x0f {
		// format code only
	} else if num <= math.MaxUint16 {
		ret += def.Byte2
	} else if num <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this array length : %d", num)
	}
	return ret, nil
}

func (s *serializer) calcStructMap(rv reflect.Value) (int, error) {
	ret := 0
	l := rv.NumField()
	num := 0
	for i := 0; i < l; i++ {
		field := rv.Field(i)
		if field.CanSet() {
			// TODO : tag check
			keySize := def.Byte1 + s.calcString(rv.Type().Field(i).Name)
			valueSize, err := s.calcSize(field)
			if err != nil {
				return 0, err
			}
			ret += keySize + valueSize
			num++
		}
	}

	// format size
	if num <= 0x0f {
		// format code only
	} else if num <= math.MaxUint16 {
		ret += def.Byte2
	} else if num <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return ret, nil
}

func (s *serializer) writeStructArray(rv reflect.Value, offset int) (int, error) {

	num := 0
	l := rv.NumField()
	for i := 0; i < l; i++ {
		field := rv.Field(i)
		if field.CanSet() {
			num++
		}
	}
	// write format
	if num <= 0x0f {
		offset = s.setByte1Int(def.FixArray+num, offset)
	} else if num <= math.MaxUint16 {
		offset = s.setByte1Int(def.Array16, offset)
		offset = s.setByte2Int(num, offset)
	} else if num <= math.MaxUint32 {
		offset = s.setByte1Int(def.Array32, offset)
		offset = s.setByte4Int(num, offset)
	}

	for i := 0; i < l; i++ {
		field := rv.Field(i)
		if field.CanSet() {
			o, err := s.create(rv.Field(i), offset)
			if err != nil {
				return 0, err
			}
			offset = o
		}
	}
	return offset, nil
}

func (s *serializer) writeStructMap(rv reflect.Value, offset int) (int, error) {
	l := rv.NumField()
	num := 0
	for i := 0; i < l; i++ {
		field := rv.Field(i)
		if field.CanSet() {
			num++
		}
	}
	// format size
	if num <= 0x0f {
		offset = s.setByte1Int(def.FixMap+num, offset)
	} else if num <= math.MaxUint16 {
		offset = s.setByte1Int(def.Map16, offset)
		offset = s.setByte2Int(num, offset)
	} else if num <= math.MaxUint32 {
		offset = s.setByte1Int(def.Map32, offset)
		offset = s.setByte4Int(num, offset)
	}

	for i := 0; i < l; i++ {
		field := rv.Field(i)
		if field.CanSet() {
			// TODO : tag check
			o := s.writeString(rv.Type().Field(i).Name, offset)
			o, err := s.create(field, o)
			if err != nil {
				return 0, err
			}
			offset = o
		}
	}
	return offset, nil
}
