package encoding

import (
	"reflect"

	"github.com/shamaton/msgpack/ext"
	"github.com/shamaton/msgpack/time"
)

var extCoderMap = map[reflect.Type]ext.Encoder{reflect.TypeOf(time.Encoder): time.Encoder}
var extCoders = []ext.Encoder{time.Encoder}

func AddExtEncoder(f ext.Encoder) {
	t := reflect.TypeOf(f)
	// ignore time
	if t == reflect.TypeOf(time.Encoder) {
		return
	}

	_, ok := extCoderMap[t]
	if !ok {
		extCoderMap[t] = f
		updateExtCoders()
	}
}

func RemoveExtEncoder(f ext.Encoder) {
	t := reflect.TypeOf(f)
	// ignore time
	if t == reflect.TypeOf(time.Encoder) {
		return
	}

	_, ok := extCoderMap[t]
	if ok {
		delete(extCoderMap, t)
		updateExtCoders()
	}
}

func updateExtCoders() {
	extCoders = make([]ext.Encoder, len(extCoderMap))
	i := 0
	for k := range extCoderMap {
		extCoders[i] = extCoderMap[k]
		i++
	}
}

/*
func (e *encoder) isDateTime(value reflect.Value) (bool, time.Time) {
	i := value.Interface()
	switch t := i.(type) {
	case time.Time:
		return true, t
	}
	return false, now
}

func (e *encoder) calcTime(t time.Time) int {
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			return def.Byte1 + def.Byte4
		}
		return def.Byte1 + def.Byte8
	}

	return def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8
}

func (e *encoder) writeTime(t time.Time, offset int) int {
	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			offset = e.setByte1Int(def.Fixext4, offset)
			offset = e.setByte1Int(def.TimeStamp, offset)
			offset = e.setByte4Uint64(data, offset)
			return offset
		}

		offset = e.setByte1Int(def.Fixext8, offset)
		offset = e.setByte1Int(def.TimeStamp, offset)
		offset = e.setByte8Uint64(data, offset)
		return offset
	}

	offset = e.setByte1Int(def.Ext8, offset)
	offset = e.setByte1Int(12, offset)
	offset = e.setByte1Int(def.TimeStamp, offset)
	offset = e.setByte4Int(t.Nanosecond(), offset)
	offset = e.setByte8Uint64(secs, offset)
	return offset
}
*/
