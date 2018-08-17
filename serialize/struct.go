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
		if s.isPublic(rv.Type().Field(i).Name) {
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
		if ok, name := s.checkField(rv.Type().Field(i)); ok {
			// TODO : tag check
			keySize := def.Byte1 + s.calcString(name)
			valueSize, err := s.calcSize(rv.Field(i))
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

func (s *serializer) writeStructArray(rv reflect.Value, offset int) int {

	num := 0
	l := rv.NumField()
	for i := 0; i < l; i++ {
		// TODO : ignore check
		if s.isPublic(rv.Type().Field(i).Name) {
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
		if s.isPublic(rv.Type().Field(i).Name) {
			offset = s.create(rv.Field(i), offset)
		}
	}
	return offset
}

func (s *serializer) writeStructMap(rv reflect.Value, offset int) int {
	l := rv.NumField()
	num := 0
	for i := 0; i < l; i++ {
		if ok, _ := s.checkField(rv.Type().Field(i)); ok {
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
		if ok, name := s.checkField(rv.Type().Field(i)); ok {
			// TODO : tag check
			offset = s.writeString(name, offset)
			offset = s.create(rv.Field(i), offset)
		}
	}
	return offset
}

func (s *serializer) checkField(field reflect.StructField) (bool, string) {
	// A to Z
	if s.isPublic(field.Name) {
		if tag := field.Tag.Get("msgpack"); tag == "ignore" {
			return false, ""
		} else if len(tag) > 0 {
			return true, tag
		}
		return true, field.Name
	}
	return false, ""
}

func (s *serializer) isPublic(name string) bool {
	return 0x41 <= name[0] && name[0] <= 0x5a
}
