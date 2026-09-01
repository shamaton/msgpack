package decoding

import (
	"encoding/binary"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/common/decodingutil"
	"github.com/shamaton/msgpack/v3/time"
)

var (
	extCoderMap = map[int8]ext.StreamDecoder{time.StreamDecoder.Code(): time.StreamDecoder}
	extCoders   = []ext.StreamDecoder{time.StreamDecoder}
)

// tryExtDecode attempts to decode the value using a registered ext decoder
// whose Go type matches the destination rv, mirroring the dispatch setStruct
// already performed for struct destinations only. readIfExtType both
// recognizes the ext wire codes and reads the full frame (header and
// payload) from the reader, so truncated input surfaces as a read error
// before any custom decoder callback runs; non-ext codes make it return with
// data == nil and no error, at the cost of a single already-read code byte.
// Returns ok=false when no registered decoder claims the bytes, so callers
// fall back to the normal kind-based decoding (which will surface a clear
// type-mismatch error).
func (d *decoder) tryExtDecode(code byte, rv reflect.Value, k reflect.Kind) (bool, error) {
	if len(extCoders) == 0 {
		return false, nil
	}
	// The only registered decoder is the built-in time.StreamDecoder, which
	// targets a struct destination; skip the ext read entirely for other
	// kinds. Use rv.Kind() rather than the passed-in k: callers such as the
	// slice-of-struct fast path pass the container's kind, not the element's.
	if len(extCoders) == 1 && rv.Kind() != reflect.Struct {
		return false, nil
	}
	innerType, data, err := d.readIfExtType(code)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil
	}
	for i := range extCoders {
		if !extCoders[i].IsType(code, innerType, len(data)) {
			continue
		}
		v, err := extCoders[i].ToValue(code, data, k)
		if err != nil {
			return false, err
		}
		if rv.Type() == reflect.TypeOf(v) {
			rv.Set(reflect.ValueOf(v))
			return true, nil
		}
	}
	return false, nil
}

// AddExtDecoder adds decoders for extension types.
func AddExtDecoder(f ext.StreamDecoder) {
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
func RemoveExtDecoder(f ext.StreamDecoder) {
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
	extCoders = make([]ext.StreamDecoder, len(extCoderMap))
	i := 0
	for k := range extCoderMap {
		extCoders[i] = extCoderMap[k]
		i++
	}
}

func (d *decoder) readIfExtType(code byte) (innerType int8, data []byte, err error) {
	switch code {
	case def.Fixext1:
		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		v, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), []byte{v}, nil

	case def.Fixext2:
		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSize2()
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil

	case def.Fixext4:
		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSize4()
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil

	case def.Fixext8:
		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSize8()
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil

	case def.Fixext16:
		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSize16()
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil

	case def.Ext8:
		bs, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		size := int(bs)

		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSizeN(size)
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil

	case def.Ext16:
		bs, err := d.readSize2()
		if err != nil {
			return 0, nil, err
		}
		size := int(binary.BigEndian.Uint16(bs))

		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSizeN(size)
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil

	case def.Ext32:
		bs, err := d.readSize4()
		if err != nil {
			return 0, nil, err
		}
		size, err := lengthFromUint32(binary.BigEndian.Uint32(bs))
		if err != nil {
			return 0, nil, err
		}

		typ, err := d.readSize1()
		if err != nil {
			return 0, nil, err
		}
		data, err = d.readSizeN(size)
		if err != nil {
			return 0, nil, err
		}
		return decodingutil.Int8FromByte(typ), data, nil
	}

	return 0, nil, nil
}
