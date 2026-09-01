package decoding

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/time"
)

var (
	extCoderMap = map[int8]ext.Decoder{time.Decoder.Code(): time.Decoder}
	extCoders   = []ext.Decoder{time.Decoder}
)

// AddExtDecoder adds decoders for extension types.
func AddExtDecoder(f ext.Decoder) {
	// ignore time
	if f.Code() == time.Decoder.Code() {
		return
	}

	_, ok := extCoderMap[f.Code()]
	if !ok {
		extCoderMap[f.Code()] = f
		updateExtCoders()
	}
}

// RemoveExtDecoder removes decoders for extension types.
func RemoveExtDecoder(f ext.Decoder) {
	// ignore time
	if f.Code() == time.Decoder.Code() {
		return
	}

	_, ok := extCoderMap[f.Code()]
	if ok {
		delete(extCoderMap, f.Code())
		updateExtCoders()
	}
}

func updateExtCoders() {
	extCoders = make([]ext.Decoder, len(extCoderMap))
	i := 0
	for k := range extCoderMap {
		extCoders[i] = extCoderMap[k]
		i++
	}
}

// tryExtDecode attempts to decode the value at offset using a registered ext
// decoder whose Go type matches the destination rv, mirroring the dispatch
// setStruct already performed for struct destinations only. extEndOffset both
// recognizes the ext wire codes and validates that the full ext frame (header
// and payload) is present, so truncated input is rejected with
// def.ErrTooShortBytes before any custom decoder callback runs; non-ext bytes
// fall out of extEndOffset's switch immediately, so this costs a single
// bounds-checked byte read for the common non-ext case. Returns ok=false
// (with offset==0) when no registered decoder claims the bytes, so callers
// fall back to the normal kind-based decoding (which will surface a clear
// type-mismatch error).
func (d *decoder) tryExtDecode(rv reflect.Value, offset int, k reflect.Kind) (int, bool, error) {
	if len(extCoders) == 0 {
		return 0, false, nil
	}
	// The only registered decoder is the built-in time.Decoder, which targets
	// a struct destination; skip the check entirely for other kinds. Use
	// rv.Kind() rather than the passed-in k: callers such as the slice-of-
	// struct fast path pass the container's kind, not the element's.
	if len(extCoders) == 1 && rv.Kind() != reflect.Struct {
		return 0, false, nil
	}
	isExt, _, err := d.extEndOffset(offset)
	if err != nil {
		return 0, false, err
	}
	if !isExt {
		return 0, false, nil
	}
	for i := range extCoders {
		if !extCoders[i].IsType(offset, &d.data) {
			continue
		}
		v, o, err := extCoders[i].AsValue(offset, k, &d.data)
		if err != nil {
			return 0, false, err
		}
		if rv.Type() == reflect.TypeOf(v) {
			rv.Set(reflect.ValueOf(v))
			return o, true, nil
		}
	}
	return 0, false, nil
}

func (d *decoder) extEndOffset(offset int) (bool, int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return false, 0, err
	}
	return d.extEndOffsetWithCode(code, offset)
}

func (d *decoder) extEndOffsetWithCode(code byte, offset int) (bool, int, error) {
	switch code {
	case def.Fixext1:
		_, offset, err := d.readSizeN(offset, def.Byte1+def.Byte1)
		return true, offset, err
	case def.Fixext2:
		_, offset, err := d.readSizeN(offset, def.Byte1+def.Byte2)
		return true, offset, err
	case def.Fixext4:
		_, offset, err := d.readSizeN(offset, def.Byte1+def.Byte4)
		return true, offset, err
	case def.Fixext8:
		_, offset, err := d.readSizeN(offset, def.Byte1+def.Byte8)
		return true, offset, err
	case def.Fixext16:
		_, offset, err := d.readSizeN(offset, def.Byte1+def.Byte16)
		return true, offset, err
	case def.Ext8:
		size, offset, err := d.readSize1(offset)
		if err != nil {
			return true, 0, err
		}
		_, offset, err = d.readSizeN(offset, def.Byte1+int(size))
		return true, offset, err
	case def.Ext16:
		sizeBytes, offset, err := d.readSize2(offset)
		if err != nil {
			return true, 0, err
		}
		_, offset, err = d.readSizeN(offset, def.Byte1+int(binary.BigEndian.Uint16(sizeBytes)))
		return true, offset, err
	case def.Ext32:
		sizeBytes, offset, err := d.readSize4(offset)
		if err != nil {
			return true, 0, err
		}
		_, offset, err = d.readSizeN(offset, def.Byte1+int(binary.BigEndian.Uint32(sizeBytes)))
		return true, offset, err
	default:
		return false, 0, nil
	}
}

/*
var zero = time.Unix(0,0)

func (d *decoder) isDateTime(offset int) bool {
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

func (d *decoder) asDateTime(offset int, k reflect.Kind) (time.Time, int, error) {
	code, offset := d.readSize1(offset)

	switch code {
	case def.Fixext4:
		_, offset = d.readSize1(offset)
		bs, offset := d.readSize4(offset)
		return time.Unix(int64(binary.BigEndian.Uint32(bs)), 0), offset, nil

	case def.Fixext8:
		_, offset = d.readSize1(offset)
		bs, offset := d.readSize8(offset)
		data64 := binary.BigEndian.Uint64(bs)
		nano := int64(data64 >> 34)
		if nano > 999999999 {
			return zero, 0, fmt.Errorf("In timestamp 64 formats, nanoseconds must not be larger than 999999999 : %d", nano)
		}
		return time.Unix(int64(data64&0x00000003ffffffff), nano), offset, nil

	case def.Ext8:
		_, offset = d.readSize1(offset)
		_, offset = d.readSize1(offset)
		nanobs, offset := d.readSize4(offset)
		secbs, offset := d.readSize8(offset)
		nano := binary.BigEndian.Uint32(nanobs)
		if nano > 999999999 {
			return zero, 0, fmt.Errorf("In timestamp 96 formats, nanoseconds must not be larger than 999999999 : %d", nano)
		}
		sec := binary.BigEndian.Uint64(secbs)
		return time.Unix(int64(sec), int64(nano)), offset, nil
	}

	return zero, 0, d.errorTemplate(code, k)
}
*/
