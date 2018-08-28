package serialize

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/def"
)

type structCache struct {
	indexes []int
	names   []string
}

var cachemap = map[reflect.Type]*structCache{}

func (s *serializer) calcStructArray(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	c, find := cachemap[t]
	if !find {
		c = &structCache{}
		for i := 0; i < rv.NumField(); i++ {
			field := t.Field(i)
			if ok, name := s.checkField(field); ok {
				size, err := s.calcSize(rv.Field(i))
				if err != nil {
					return 0, err
				}
				ret += size
				c.indexes = append(c.indexes, i)
				c.names = append(c.names, name)
			}
		}
		cachemap[t] = c
	} else {
		for i := 0; i < len(c.indexes); i++ {
			size, err := s.calcSize(rv.Field(c.indexes[i]))
			if err != nil {
				return 0, err
			}
			ret += size
		}
	}

	// format size
	l := len(c.indexes)
	if l <= 0x0f {
		// format code only
	} else if l <= math.MaxUint16 {
		ret += def.Byte2
	} else if l <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return ret, nil
}

func (s *serializer) calcStructMap(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	c, find := cachemap[t]
	if !find {
		c = &structCache{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, name := s.checkField(rv.Type().Field(i)); ok {
				keySize := def.Byte1 + s.calcString(name)
				valueSize, err := s.calcSize(rv.Field(i))
				if err != nil {
					return 0, err
				}
				ret += keySize + valueSize
				c.indexes = append(c.indexes, i)
				c.names = append(c.names, name)
			}
		}
		cachemap[t] = c
	} else {
		for i := 0; i < len(c.indexes); i++ {
			keySize := def.Byte1 + s.calcString(c.names[i])
			valueSize, err := s.calcSize(rv.Field(c.indexes[i]))
			if err != nil {
				return 0, err
			}
			ret += keySize + valueSize
		}
	}

	// format size
	l := len(c.indexes)
	if l <= 0x0f {
		// format code only
	} else if l <= math.MaxUint16 {
		ret += def.Byte2
	} else if l <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return ret, nil
}

func (s *serializer) writeStructArray(rv reflect.Value, offset int) int {

	c := cachemap[rv.Type()]

	// write format
	num := len(c.indexes)
	if num <= 0x0f {
		offset = s.setByte1Int(def.FixArray+num, offset)
	} else if num <= math.MaxUint16 {
		offset = s.setByte1Int(def.Array16, offset)
		offset = s.setByte2Int(num, offset)
	} else if num <= math.MaxUint32 {
		offset = s.setByte1Int(def.Array32, offset)
		offset = s.setByte4Int(num, offset)
	}

	for i := 0; i < num; i++ {
		offset = s.create(rv.Field(c.indexes[i]), offset)
	}
	return offset
}

func (s *serializer) writeStructMap(rv reflect.Value, offset int) int {

	c := cachemap[rv.Type()]

	// format size
	num := len(c.indexes)
	if num <= 0x0f {
		offset = s.setByte1Int(def.FixMap+num, offset)
	} else if num <= math.MaxUint16 {
		offset = s.setByte1Int(def.Map16, offset)
		offset = s.setByte2Int(num, offset)
	} else if num <= math.MaxUint32 {
		offset = s.setByte1Int(def.Map32, offset)
		offset = s.setByte4Int(num, offset)
	}

	for i := 0; i < num; i++ {
		offset = s.writeString(c.names[i], offset)
		offset = s.create(rv.Field(c.indexes[i]), offset)
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
