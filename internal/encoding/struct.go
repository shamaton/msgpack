package encoding

import (
	"math"
	"reflect"
	"sync"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/internal/common"
)

type structCache struct {
	// common fields
	names  []string
	omits  []bool
	noOmit bool

	// fast path detection
	hasEmbedded bool

	// fast path (hasEmbedded == false): direct field access
	simpleIndexes []int

	// embedded path (hasEmbedded == true): path-based access
	indexes   [][]int   // field path (support for embedded structs)
	omitPaths [][][]int // embedded omitempty parent paths

	common.Common
}

var cachemap = sync.Map{}

type (
	structCalcFunc  func(rv reflect.Value) (int, error)
	structWriteFunc func(rv reflect.Value, offset int) int
)

// getFieldByPath returns the field value by following the path of indices.
// The bool indicates whether the path was reachable (no nil pointer in the path).
func getFieldByPath(rv reflect.Value, path []int) (reflect.Value, bool) {
	for _, idx := range path {
		// Handle pointer indirection if needed
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				// Return invalid value if pointer is nil
				return reflect.Value{}, false
			}
			rv = rv.Elem()
		}
		rv = rv.Field(idx)
	}
	return rv, true
}

func shouldOmitByParent(rv reflect.Value, omitPaths [][]int) bool {
	for _, path := range omitPaths {
		parentValue, ok := getFieldByPath(rv, path)
		if !ok || parentValue.IsZero() {
			return true
		}
	}
	return false
}

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
		c = &structCache{}
		fields := e.CollectFields(t, nil)

		// detect embedded fields
		hasEmbedded := false
		for _, f := range fields {
			if len(f.Path) > 1 || len(f.OmitPaths) > 0 {
				hasEmbedded = true
				break
			}
		}
		c.hasEmbedded = hasEmbedded

		omitCount := 0
		for _, field := range fields {
			c.names = append(c.names, field.Name)
			c.omits = append(c.omits, field.Omit)
			if hasEmbedded {
				c.indexes = append(c.indexes, field.Path)
				c.omitPaths = append(c.omitPaths, field.OmitPaths)
			} else {
				c.simpleIndexes = append(c.simpleIndexes, field.Path[0])
			}
			if field.Omit {
				omitCount++
			}
		}
		c.noOmit = omitCount == 0
		cachemap.Store(t, c)
	} else {
		c = cache.(*structCache)
	}

	// calculate size based on path type
	var numFields int
	if c.hasEmbedded {
		numFields = len(c.indexes)
		for i := 0; i < numFields; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				fieldValue = reflect.Value{}
			}
			size, err := e.calcSize(fieldValue)
			if err != nil {
				return 0, err
			}
			ret += size
		}
	} else {
		numFields = len(c.simpleIndexes)
		for i := 0; i < numFields; i++ {
			size, err := e.calcSize(rv.Field(c.simpleIndexes[i]))
			if err != nil {
				return 0, err
			}
			ret += size
		}
	}

	// format size
	size, err := e.calcLength(numFields)
	if err != nil {
		return 0, err
	}
	ret += size
	return ret, nil
}

func (e *encoder) calcStructMap(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	cache, find := cachemap.Load(t)
	var c *structCache
	if !find {
		c = &structCache{}
		fields := e.CollectFields(t, nil)

		// detect embedded fields
		hasEmbedded := false
		for _, f := range fields {
			if len(f.Path) > 1 || len(f.OmitPaths) > 0 {
				hasEmbedded = true
				break
			}
		}
		c.hasEmbedded = hasEmbedded

		omitCount := 0
		for _, field := range fields {
			c.names = append(c.names, field.Name)
			c.omits = append(c.omits, field.Omit)
			if hasEmbedded {
				c.indexes = append(c.indexes, field.Path)
				c.omitPaths = append(c.omitPaths, field.OmitPaths)
			} else {
				c.simpleIndexes = append(c.simpleIndexes, field.Path[0])
			}
			if field.Omit {
				omitCount++
			}
		}
		c.noOmit = omitCount == 0
		cachemap.Store(t, c)
	} else {
		c = cache.(*structCache)
	}

	l := 0
	if c.hasEmbedded {
		for i := 0; i < len(c.indexes); i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				continue
			}
			size, err := e.calcSizeWithOmitEmpty(fieldValue, c.names[i], c.omits[i])
			if err != nil {
				return 0, err
			}
			ret += size
			if size > 0 {
				l++
			}
		}
	} else {
		for i := 0; i < len(c.simpleIndexes); i++ {
			size, err := e.calcSizeWithOmitEmpty(rv.Field(c.simpleIndexes[i]), c.names[i], c.omits[i])
			if err != nil {
				return 0, err
			}
			ret += size
			if size > 0 {
				l++
			}
		}
	}

	// format size
	size, err := e.calcLength(l)
	if err != nil {
		return 0, err
	}
	ret += size
	return ret, nil
}

func (e *encoder) calcSizeWithOmitEmpty(rv reflect.Value, name string, omit bool) (int, error) {
	keySize := 0
	valueSize := 0
	if !omit || !rv.IsZero() {
		keySize = e.calcString(name)
		vSize, err := e.calcSize(rv)
		if err != nil {
			return 0, err
		}
		valueSize = vSize
	}
	return keySize + valueSize, nil
}

func (e *encoder) getStructWriter(typ reflect.Type) structWriteFunc {
	for i := range extCoders {
		if extCoders[i].Type() == typ {
			return func(rv reflect.Value, offset int) int {
				return extCoders[i].WriteToBytes(rv, offset, &e.d)
			}
		}
	}

	if e.asArray {
		return e.writeStructArray
	}
	return e.writeStructMap
}

func (e *encoder) writeStruct(rv reflect.Value, offset int) int {
	/*
		if isTime, tm := e.isDateTime(rv); isTime {
			return e.writeTime(tm, offset)
		}
	*/

	for i := range extCoders {
		if extCoders[i].Type() == rv.Type() {
			return extCoders[i].WriteToBytes(rv, offset, &e.d)
		}
	}

	if e.asArray {
		return e.writeStructArray(rv, offset)
	}
	return e.writeStructMap(rv, offset)
}

func (e *encoder) writeStructArray(rv reflect.Value, offset int) int {
	cache, _ := cachemap.Load(rv.Type())
	c := cache.(*structCache)

	// write format
	var num int
	if c.hasEmbedded {
		num = len(c.indexes)
	} else {
		num = len(c.simpleIndexes)
	}

	if num <= 0x0f {
		offset = e.setByte1Int(def.FixArray+num, offset)
	} else if num <= math.MaxUint16 {
		offset = e.setByte1Int(def.Array16, offset)
		offset = e.setByte2Int(num, offset)
	} else if uint(num) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Array32, offset)
		offset = e.setByte4Int(num, offset)
	}

	if c.hasEmbedded {
		for i := 0; i < num; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				fieldValue = reflect.Value{}
			}
			offset = e.create(fieldValue, offset)
		}
	} else {
		for i := 0; i < num; i++ {
			offset = e.create(rv.Field(c.simpleIndexes[i]), offset)
		}
	}
	return offset
}

func (e *encoder) writeStructMap(rv reflect.Value, offset int) int {
	cache, _ := cachemap.Load(rv.Type())
	c := cache.(*structCache)

	// format size
	l := 0
	if c.hasEmbedded {
		num := len(c.indexes)
		for i := 0; i < num; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				continue
			}
			if c.noOmit || !c.omits[i] || !fieldValue.IsZero() {
				l++
			}
		}
	} else {
		num := len(c.simpleIndexes)
		for i := 0; i < num; i++ {
			if c.noOmit || !c.omits[i] || !rv.Field(c.simpleIndexes[i]).IsZero() {
				l++
			}
		}
	}

	if l <= 0x0f {
		offset = e.setByte1Int(def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(l, offset)
	}

	if c.hasEmbedded {
		num := len(c.indexes)
		for i := 0; i < num; i++ {
			fieldValue, ok := getFieldByPath(rv, c.indexes[i])
			if shouldOmitByParent(rv, c.omitPaths[i]) || !ok {
				continue
			}
			if c.noOmit || !c.omits[i] || !fieldValue.IsZero() {
				offset = e.writeString(c.names[i], offset)
				offset = e.create(fieldValue, offset)
			}
		}
	} else {
		num := len(c.simpleIndexes)
		for i := 0; i < num; i++ {
			fieldValue := rv.Field(c.simpleIndexes[i])
			if c.noOmit || !c.omits[i] || !fieldValue.IsZero() {
				offset = e.writeString(c.names[i], offset)
				offset = e.create(fieldValue, offset)
			}
		}
	}
	return offset
}
