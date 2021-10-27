package encoding

import (
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/internal/common"
)

type structCache struct {
	indexes []int
	names   []string
	common.Common
}

var cachemap = sync.Map{}

type structCalcFunc func(rv reflect.Value) (int, error)
type structWriteFunc func(rv reflect.Value, writer Writer) error

func (e *encoder) getStructCalc(typ reflect.Type) structCalcFunc {

	for j := range extCoders {
		if extCoders[j].Type() == typ {
			return extCoders[j].CalcByteSize
		}
	}
	if e.asArray {
		return e.calcStructArray
	}
	return e.calcStructMap

}

func (e *encoder) calcStruct(rv reflect.Value) (int, error) {

	//if isTime, tm := e.isDateTime(rv); isTime {
	//	size := e.calcTime(tm)
	//	return size, nil
	//}

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].CalcByteSize(rv)
		}
	}

	if e.asArray {
		return e.calcStructArray(rv)
	}
	return e.calcStructMap(rv)
}

func (e *encoder) calcStructArray(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	cache, find := cachemap.Load(t)
	var c *structCache
	if !find {
		var err error
		c, ret, err = e.cacheStructArray(rv, ret, t)
		if err != nil {
			return 0, err
		}
	} else {
		c = cache.(*structCache)
		for i := 0; i < len(c.indexes); i++ {
			size, err := e.calcSize(rv.Field(c.indexes[i]))
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
	} else if uint(l) <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return ret, nil
}



func (e *encoder) cacheStructArray(rv reflect.Value, ret int, t reflect.Type) (*structCache, int, error) {
	c := &structCache{}
	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		if ok, name := e.CheckField(field); ok {
			size, err := e.calcSize(rv.Field(i))
			if err != nil {
				return nil, 0, err
			}
			ret += size
			c.indexes = append(c.indexes, i)
			c.names = append(c.names, name)
		}
	}

	cachemap.Store(t, c)
	return c, ret, nil
}

func (e *encoder) calcStructMap(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	cache, find := cachemap.Load(t)
	var c *structCache
	if !find {
		var err error
		c, ret, err = e.cacheStructMap(rv, ret, t)
		if err != nil {
			return 0, err
		}
	} else {
		c = cache.(*structCache)
		for i := 0; i < len(c.indexes); i++ {
			keySize := def.Byte1 + e.calcString(c.names[i])
			valueSize, err := e.calcSize(rv.Field(c.indexes[i]))
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
	} else if uint(l) <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		// not supported error
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return ret, nil
}

func (e *encoder) cacheStructMap(rv reflect.Value, ret int, t reflect.Type) (*structCache, int, error) {
	c := &structCache{}
	for i := 0; i < rv.NumField(); i++ {
		if ok, name := e.CheckField(rv.Type().Field(i)); ok {
			keySize := def.Byte1 + e.calcString(name)
			valueSize, err := e.calcSize(rv.Field(i))
			if err != nil {
				return nil, 0, err
			}
			ret += keySize + valueSize
			c.indexes = append(c.indexes, i)
			c.names = append(c.names, name)
		}
	}

	cachemap.Store(t, c)
	return c, ret, nil
}

func (e *encoder) getStructWriter(typ reflect.Type) structWriteFunc {
	for i := range extCoders {
		if extCoders[i].Type() == typ {
			return func(rv reflect.Value, writer Writer) error {
				return extCoders[i].WriteToBytes(rv, writer)
			}
		}
	}

	if e.asArray {
		return e.writeStructArray
	}
	return e.writeStructMap
}

func (e *encoder) writeStruct(rv reflect.Value, writer Writer) error {
	/*
		if isTime, tm := e.isDateTime(rv); isTime {
			return e.writeTime(tm, offset)
		}
	*/

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].WriteToBytes(rv, writer)
		}
	}

	if e.asArray {
		return e.writeStructArray(rv, writer)
	}
	return e.writeStructMap(rv, writer)
}

func (e *encoder) writeStructArray(rv reflect.Value, writer Writer) error {
	cache, find := cachemap.Load(rv.Type())
	if !find {
		var err error
		cache, _, err = e.cacheStructArray(rv, 0, rv.Type())
		if err != nil {
			return err
		}
	}
	c := cache.(*structCache)

	// write format
	num := len(c.indexes)
	if num <= 0x0f {
		err := e.setByte1Int(def.FixArray+num, writer)
		if err != nil {
			return err
		}
	} else if num <= math.MaxUint16 {
		err := e.setByte1Int(def.Array16, writer)
		if err != nil {
			return err
		}
		err = e.setByte2Int(num, writer)
		if err != nil {
			return err
		}
	} else if uint(num) <= math.MaxUint32 {
		err := e.setByte1Int(def.Array32, writer)
		if err != nil {
			return err
		}
		err = e.setByte4Int(num, writer)
		if err != nil {
			return err
		}
	}

	for i := 0; i < num; i++ {
		err := e.create(rv.Field(c.indexes[i]), writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) writeStructMap(rv reflect.Value, writer Writer) error {
	cache, find := cachemap.Load(rv.Type())
	if !find {
		var err error
		cache, _, err = e.cacheStructMap(rv, 0, rv.Type())
		if err != nil {
			return err
		}
	}
	c := cache.(*structCache)

	// format size
	num := len(c.indexes)
	if num <= 0x0f {
		err := e.setByte1Int(def.FixMap+num, writer)
		if err != nil {
			return err
		}
	} else if num <= math.MaxUint16 {
		err := e.setByte1Int(def.Map16, writer)
		if err != nil {
			return err
		}

		err = e.setByte2Int(num, writer)
		if err != nil {
			return err
		}
	} else if uint(num) <= math.MaxUint32 {
		err := e.setByte1Int(def.Map32, writer)
		if err != nil {
			return err
		}

		err = e.setByte4Int(num, writer)
		if err != nil {
			return err
		}
	}

	for i := 0; i < num; i++ {
		err := e.writeString(c.names[i], writer)
		if err != nil {
			return err
		}

		err = e.create(rv.Field(c.indexes[i]), writer)
		if err != nil {
			return err
		}
	}

	return nil
}
