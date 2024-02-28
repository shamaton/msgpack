package decoding

import (
	"encoding/binary"
	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	"github.com/shamaton/msgpack/v2/time"
	"io"
)

var extCoderMap = map[int8]ext.StreamDecoder{time.StreamDecoder.Code(): time.StreamDecoder}
var extCoders = []ext.StreamDecoder{time.StreamDecoder}

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

func readIfExtType(r io.Reader, code byte) (innerType int8, data []byte, err error) {
	switch code {
	case def.Fixext1:
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSizeN(r, def.Byte1)
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Fixext2:
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSize2(r)
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Fixext4:
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSize4(r)
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Fixext8:
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSize8(r)
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Fixext16:
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSizeN(r, def.Byte8+def.Byte8)
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Ext8:
		bs, err := readSizeN(r, def.Byte1)
		if err != nil {
			return 0, nil, err
		}
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSizeN(r, int(binary.BigEndian.Uint32(bs)))
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Ext16:
		bs, err := readSize2(r)
		if err != nil {
			return 0, nil, err
		}
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSizeN(r, int(binary.BigEndian.Uint32(bs)))
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil

	case def.Ext32:
		bs, err := readSize4(r)
		if err != nil {
			return 0, nil, err
		}
		typ, err := readSize1(r)
		if err != nil {
			return 0, nil, err
		}
		data, err = readSizeN(r, int(binary.BigEndian.Uint32(bs)))
		if err != nil {
			return 0, nil, err
		}
		return int8(typ), data, nil
	}

	return 0, nil, nil
}
