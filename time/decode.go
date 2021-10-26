package time

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/shamaton/msgpack/v2/def"
)

var zero = time.Unix(0, 0)

var Decoder = new(timeDecoder)

type timeDecoder struct{}

func (td *timeDecoder) Code() int8 {
	return def.TimeStamp
}

func (td *timeDecoder) AsValue(data []byte, k reflect.Kind) (interface{}, error) {
	switch len(data) {
	case 4:
		return time.Unix(int64(binary.BigEndian.Uint32(data)), 0), nil

	case 8:
		data64 := binary.BigEndian.Uint64(data)
		nano := int64(data64 >> 34)
		if nano > 999999999 {
			return zero, fmt.Errorf("In timestamp 64 formats, nanoseconds must not be larger than 999999999 : %d", nano)
		}
		return time.Unix(int64(data64&0x00000003ffffffff), nano), nil

	case 12:
		nanobs := data[:4]
		secbs := data[4:]
		nano := binary.BigEndian.Uint32(nanobs)
		if nano > 999999999 {
			return zero, fmt.Errorf("In timestamp 96 formats, nanoseconds must not be larger than 999999999 : %d", nano)
		}
		sec := binary.BigEndian.Uint64(secbs)
		return time.Unix(int64(sec), int64(nano)), nil
	}

	return zero, fmt.Errorf("should not reach this line!! data %x decoding %v", data, k)
}
