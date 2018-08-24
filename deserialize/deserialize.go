package deserialize

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/shamaton/msgpack/def"
)

type deserializer struct {
	data    []byte
	asArray bool
}

func Exec(data []byte, holder interface{}, asArray bool) error {
	d := deserializer{data: data, asArray: asArray}

	rv := reflect.ValueOf(holder)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", holder)
	}

	rv = rv.Elem()
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	last, err := d.deserialize(rv, 0)
	if err != nil {
		return err
	}
	if len(data) != last {
		return fmt.Errorf("failed deserialization size=%d, last=%d", len(data), last)
	}
	return err
}

func (d *deserializer) deserialize(rv reflect.Value, offset int) (int, error) {
	k := rv.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, o, err := d.asInt(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetInt(v)
		offset = o

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, o, err := d.asUint(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetUint(v)
		offset = o

	case reflect.Float32:
		v, o, err := d.asFloat32(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetFloat(float64(v))
		offset = o

	case reflect.Float64:
		v, o, err := d.asFloat64(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetFloat(v)
		offset = o

	case reflect.String:
		v, o, err := d.asString(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetString(v)
		offset = o

	case reflect.Bool:
		v, o, err := d.asBool(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetBool(v)
		offset = o

	case reflect.Slice:
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}
		// byte slice
		if d.isCodeBin(d.data[offset]) {
			bs, offset, err := d.asBin(offset, k)
			if err != nil {
				return 0, err
			}
			rv.SetBytes(bs)
			return offset, nil
		}
		// string to bytes
		if d.isCodeString(d.data[offset]) {
			l, offset, err := d.stringByteLength(offset, k)
			if err != nil {
				return 0, err
			}
			bs, offset := d.asStringByte(offset, l, k)
			rv.SetBytes(bs)
			return offset, nil
		}

		// get slice length
		l, o, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		// check fixed type
		fixedOffset, found, err := d.asFixedSlice(rv, o, l)
		if err != nil {
			return 0, err
		}
		if found {
			return fixedOffset, nil
		}

		// create slice dynamically
		e := rv.Type().Elem()
		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)
		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			tmpSlice.Index(i).Set(v)
		}
		rv.Set(tmpSlice)
		offset = o

	case reflect.Array:
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}
		// byte slice
		if d.isCodeBin(d.data[offset]) {
			// todo : length check
			bs, offset, err := d.asBin(offset, k)
			if err != nil {
				return 0, err
			}
			rv.SetBytes(bs)
			return offset, nil
		}
		// string to bytes
		if d.isCodeString(d.data[offset]) {
			l, offset, err := d.stringByteLength(offset, k)
			if err != nil {
				return 0, err
			}
			if l > rv.Len() {
				return 0, fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), l)
			}
			bs, offset := d.asStringByte(offset, l, k)
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return offset, nil
		}

		// get slice length
		l, o, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		if l > rv.Len() {
			return 0, fmt.Errorf("%v len is %d, but msgpack has %d elements", rv.Type(), rv.Len(), l)
		}

		// create array dynamically
		for i := 0; i < l; i++ {
			o, err = d.deserialize(rv.Index(i), o)
			if err != nil {
				return 0, err
			}
		}
		offset = o

	case reflect.Map:
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}

		// get map length
		l, o, err := d.mapLength(offset, k)
		if err != nil {
			return 0, err
		}

		// check fixed type
		fixedOffset, found, err := d.asFixedMap(rv, o, l)
		if err != nil {
			return 0, err
		}
		if found {
			return fixedOffset, nil
		}

		// create dynamically
		key := rv.Type().Key()
		value := rv.Type().Elem()
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		for i := 0; i < l; i++ {
			k := reflect.New(key).Elem()
			v := reflect.New(value).Elem()
			o, err = d.deserialize(k, o)
			if err != nil {
				return 0, err
			}
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			rv.SetMapIndex(k, v)
		}
		offset = o

	case reflect.Struct:
		// todo : ext
		if d.isDateTime(offset) {
			dt, offset, err := d.asDateTime(offset, k)
			if err != nil {
				return 0, err
			}
			rv.Set(reflect.ValueOf(dt))
			return offset, nil
		}

		if d.asArray {
			l, o, err := d.sliceLength(offset, k)
			if err != nil {
				return 0, err
			}
			// find or create reference
			cm, findCache := cachemap2[rv.Type()]
			if !findCache {
				cm = &structCache2{}
				for i := 0; i < rv.NumField(); i++ {
					if ok, _ := d.checkField(rv.Type().Field(i)); ok {
						cm.m = append(cm.m, i)
					}
				}
				cachemap2[rv.Type()] = cm
			}
			// set value
			for i := 0; i < l; i++ {
				if i < len(cm.m) {
					o, err = d.deserialize(rv.Field(cm.m[i]), o)
					if err != nil {
						return 0, err
					}
				} else {
					o = d.jumpByte(o)
				}
			}
			offset = o

		} else {
			l, o, err := d.mapLength(offset, k)
			if err != nil {
				return 0, err
			}
			// find or create reference
			cm, cacheFind := cachemap[rv.Type()]
			if !cacheFind {
				cm = &structCache{m: map[string]int{}}
				for i := 0; i < rv.NumField(); i++ {
					if ok, name := d.checkField(rv.Type().Field(i)); ok {
						cm.m[name] = i
					}
				}
				cachemap[rv.Type()] = cm
			}
			// set value if string correct
			for i := 0; i < l; i++ {
				key, o2, err := d.asString(o, k)
				if err != nil {
					return 0, err
				}
				if _, ok := cm.m[key]; ok {
					o2, err = d.deserialize(rv.Field(cm.m[key]), o2)
					if err != nil {
						return 0, err
					}
				} else {
					o2 = d.jumpByte(o2)
				}
				o = o2
			}
			offset = o
		}

	case reflect.Ptr:
		fmt.Println("nnnnnnnnill")
		// nil
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}
		o, err := d.deserialize(rv.Elem(), offset)
		if err != nil {
			return 0, err
		}
		offset = o

	case reflect.Interface:
		// all type...

	default:
		return 0, d.errorTemplate(d.data[offset], k)
	}
	return offset, nil
}

type structCache struct {
	m map[string]int
}

type structCache2 struct {
	m []int
}

var cachemap = map[reflect.Type]*structCache{}
var cachemap2 = map[reflect.Type]*structCache2{}

// todo : change method name
func (d *deserializer) jumpByte(offset int) int {
	code, offset := d.readSize1(offset)
	switch {
	case code == def.True, code == def.False, code == def.Nil:
		// do nothing

	case d.isPositiveFixNum(code) || d.isNegativeFixNum(code):
		// do nothing
	case code == def.Uint8, code == def.Int8:
		offset += def.Byte1
	case code == def.Uint16, code == def.Int16:
		offset += def.Byte2
	case code == def.Uint32, code == def.Int32, code == def.Float32:
		offset += def.Byte4
	case code == def.Uint64, code == def.Int64, code == def.Float64:
		offset += def.Byte8

	case d.isFixString(code):
		offset += int(code - def.FixStr)
	case code == def.Str8, code == def.Bin8:
		b, offset := d.readSize1(offset)
		offset += int(b)
	case code == def.Str16, code == def.Bin16:
		bs, offset := d.readSize2(offset)
		offset += int(binary.BigEndian.Uint16(bs))
	case code == def.Str32, code == def.Bin32:
		bs, offset := d.readSize4(offset)
		offset += int(binary.BigEndian.Uint32(bs))

	case d.isFixSlice(code):
		l := int(code - def.FixStr)
		for i := 0; i < l; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Array16:
		bs, offset := d.readSize2(offset)
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Array32:
		bs, offset := d.readSize4(offset)
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			offset += d.jumpByte(offset)
		}

	case d.isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Map16:
		bs, offset := d.readSize2(offset)
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			offset += d.jumpByte(offset)
		}
	case code == def.Map32:
		bs, offset := d.readSize4(offset)
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			offset += d.jumpByte(offset)
		}

	case code == def.Fixext1:
		offset += def.Byte1 + def.Byte1
	case code == def.Fixext2:
		offset += def.Byte1 + def.Byte2
	case code == def.Fixext4:
		offset += def.Byte1 + def.Byte4
	case code == def.Fixext8:
		offset += def.Byte1 + def.Byte8
	case code == def.Fixext16:
		offset += def.Byte1 + def.Byte16

	case code == def.Ext8:
		b, offset := d.readSize1(offset)
		offset += def.Byte1 + int(b)
	case code == def.Ext16:
		bs, offset := d.readSize2(offset)
		offset += def.Byte1 + int(binary.BigEndian.Uint16(bs))
	case code == def.Ext32:
		bs, offset := d.readSize4(offset)
		offset += def.Byte1 + int(binary.BigEndian.Uint32(bs))

	}
	return offset
}

// todo same method...
func (d *deserializer) checkField(field reflect.StructField) (bool, string) {
	// A to Z
	if d.isPublic(field.Name) {
		if tag := field.Tag.Get("msgpack"); tag == "ignore" {
			return false, ""
		} else if len(tag) > 0 {
			return true, tag
		}
		return true, field.Name
	}
	return false, ""
}

// todo same method...
func (d *deserializer) isPublic(name string) bool {
	return 0x41 <= name[0] && name[0] <= 0x5a
}

func (d *deserializer) asFixedMap(rv reflect.Value, offset int, l int) (int, bool, error) {
	t := rv.Type()

	keyKind := rv.Type().Key().Kind()
	valueKind := rv.Type().Elem().Kind()
	switch t {
	case typeMapStringInt:
		m := make(map[string]int, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = int(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil
	}

	return offset, false, nil
}

func (d *deserializer) asFixedSlice(rv reflect.Value, offset int, l int) (int, bool, error) {
	t := rv.Type()
	k := t.Elem().Kind()
	switch t {
	case typeIntSlice:
		sli := make([]int, l)
		for i := range sli {
			v, o, err := d.asInt(offset, k)
			if err != nil {
				return 0, false, err
			}
			sli[i] = int(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(sli))
		return offset, true, nil

	case typeInt8Slice:
		sli := make([]int8, l)
		for i := range sli {
			v, o, err := d.asInt(offset, k)
			if err != nil {
				return 0, false, err
			}
			sli[i] = int8(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(sli))
		return offset, true, nil
	}

	return offset, false, nil
}

func (d *deserializer) isPositiveFixNum(v byte) bool {
	if def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax {
		return true
	}
	return false
}

func (d *deserializer) isNegativeFixNum(v byte) bool {
	if def.NegativeFixintMin <= int8(v) && int8(v) <= def.NegativeFixintMax {
		return true
	}
	return false
}

func (d *deserializer) isFixString(v byte) bool {
	return def.FixStr <= v && v <= def.FixStr+0x1f
}

func (d *deserializer) isCodeBin(v byte) bool {
	switch v {
	case def.Bin8, def.Bin16, def.Bin32:
		return true
	}
	return false
}

func (d *deserializer) isCodeNil(v byte) bool {
	return def.Nil == v
}

func (d *deserializer) isDateTime(offset int) bool {
	code, offset := d.readSize1(offset)

	if code == def.Fixext4 {
		t, _ := d.readSize1(offset)
		return int8(t) == def.TimeStamp
	} else if code == def.Fixext8 {
		t, _ := d.readSize1(offset)
		return int8(t) == def.TimeStamp
	} else if code == def.Ext8 {
		l, offset := d.readSize1(offset)
		t, _ := d.readSize1(offset)
		return l == 12 && int8(t) == def.TimeStamp
	}
	return false
}

func (d *deserializer) asDateTime(offset int, k reflect.Kind) (time.Time, int, error) {
	code, offset := d.readSize1(offset)

	// TODO : In timestamp 64 and timestamp 96 formats, nanoseconds must not be larger than 999999999.

	switch code {
	case def.Fixext4:
		_, offset = d.readSize1(offset)
		bs, offset := d.readSize4(offset)
		return time.Unix(int64(binary.BigEndian.Uint32(bs)), 0), offset, nil

	case def.Fixext8:
		_, offset = d.readSize1(offset)
		bs, offset := d.readSize8(offset)
		data64 := binary.BigEndian.Uint64(bs)
		return time.Unix(int64(data64&0x00000003ffffffff), int64(data64>>34)), offset, nil

	case def.Ext8:
		_, offset = d.readSize1(offset)
		_, offset = d.readSize1(offset)
		nanobs, offset := d.readSize4(offset)
		secbs, offset := d.readSize8(offset)
		nano := binary.BigEndian.Uint32(nanobs)
		sec := binary.BigEndian.Uint64(secbs)
		return time.Unix(int64(sec), int64(nano)), offset, nil
	}

	return time.Now(), 0, d.errorTemplate(code, k)
}

func (d *deserializer) errorTemplate(code byte, k reflect.Kind) error {
	return fmt.Errorf("msgpack : invalid code %x decoding %v", code, k)
}
