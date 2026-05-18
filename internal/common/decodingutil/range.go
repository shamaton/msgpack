package decodingutil

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
)

func RangeError(value interface{}, k reflect.Kind) error {
	return fmt.Errorf("%w %v decoding as %v", def.ErrValueOutOfRange, value, k)
}

func Int64FromUint64(v uint64, k reflect.Kind) (int64, error) {
	if v > math.MaxInt64 {
		return 0, RangeError(v, k)
	}
	return int64(v), nil
}

func Int64FromFloat64(v float64, k reflect.Kind) (int64, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) || v < -1<<63 || v >= 1<<63 {
		return 0, RangeError(v, k)
	}
	return int64(v), nil
}

func Int64FromInt8Byte(v byte) int64 {
	if v <= math.MaxInt8 {
		return int64(v)
	}
	return int64(v) - (math.MaxUint8 + 1)
}

func Int64FromInt16Bits(v uint16) int64 {
	if v <= math.MaxInt16 {
		return int64(v)
	}
	return int64(v) - (math.MaxUint16 + 1)
}

func Int64FromInt32Bits(v uint32) int64 {
	if v <= math.MaxInt32 {
		return int64(v)
	}
	return int64(v) - (math.MaxUint32 + 1)
}

func Int64FromInt64Bits(v uint64) int64 {
	if v <= math.MaxInt64 {
		return int64(v)
	}
	if v == 1<<63 {
		return math.MinInt64
	}
	return -int64(^v + 1) // #nosec G115 -- value is checked to fit before converting the two's-complement magnitude.
}

func Uint64FromInt64(v int64, k reflect.Kind) (uint64, error) {
	if v < 0 {
		return 0, RangeError(v, k)
	}
	return uint64(v), nil
}

func Int8FromByte(v byte) int8 {
	return int8(v) // #nosec G115 -- MessagePack ext type codes are signed one-byte values.
}

func Int8FromInt64(v int64, k reflect.Kind) (int8, error) {
	if v < math.MinInt8 || v > math.MaxInt8 {
		return 0, RangeError(v, k)
	}
	return int8(v), nil
}

func IntFromInt64(v int64, k reflect.Kind) (int, error) {
	if def.IsIntSize32 && (v < math.MinInt32 || v > math.MaxInt32) {
		return 0, RangeError(v, k)
	}
	return int(v), nil
}

func Int16FromInt64(v int64, k reflect.Kind) (int16, error) {
	if v < math.MinInt16 || v > math.MaxInt16 {
		return 0, RangeError(v, k)
	}
	return int16(v), nil
}

func Int32FromInt64(v int64, k reflect.Kind) (int32, error) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, RangeError(v, k)
	}
	return int32(v), nil
}

func Uint8FromUint64(v uint64, k reflect.Kind) (uint8, error) {
	if v > math.MaxUint8 {
		return 0, RangeError(v, k)
	}
	return uint8(v), nil
}

func UintFromUint64(v uint64, k reflect.Kind) (uint, error) {
	if def.IsIntSize32 && v > math.MaxUint32 {
		return 0, RangeError(v, k)
	}
	return uint(v), nil
}

func IntValueForKind(v int64, k reflect.Kind) (int64, error) {
	switch k {
	case reflect.Int:
		vv, err := IntFromInt64(v, k)
		if err != nil {
			return 0, err
		}
		return int64(vv), nil
	case reflect.Int8:
		vv, err := Int8FromInt64(v, k)
		if err != nil {
			return 0, err
		}
		return int64(vv), nil
	case reflect.Int16:
		vv, err := Int16FromInt64(v, k)
		if err != nil {
			return 0, err
		}
		return int64(vv), nil
	case reflect.Int32:
		vv, err := Int32FromInt64(v, k)
		if err != nil {
			return 0, err
		}
		return int64(vv), nil
	case reflect.Int64:
		return v, nil
	}
	return 0, RangeError(v, k)
}

func Uint16FromUint64(v uint64, k reflect.Kind) (uint16, error) {
	if v > math.MaxUint16 {
		return 0, RangeError(v, k)
	}
	return uint16(v), nil
}

func Uint32FromUint64(v uint64, k reflect.Kind) (uint32, error) {
	if v > math.MaxUint32 {
		return 0, RangeError(v, k)
	}
	return uint32(v), nil
}

func UintValueForKind(v uint64, k reflect.Kind) (uint64, error) {
	switch k {
	case reflect.Uint:
		vv, err := UintFromUint64(v, k)
		if err != nil {
			return 0, err
		}
		return uint64(vv), nil
	case reflect.Uint8:
		vv, err := Uint8FromUint64(v, k)
		if err != nil {
			return 0, err
		}
		return uint64(vv), nil
	case reflect.Uint16:
		vv, err := Uint16FromUint64(v, k)
		if err != nil {
			return 0, err
		}
		return uint64(vv), nil
	case reflect.Uint32:
		vv, err := Uint32FromUint64(v, k)
		if err != nil {
			return 0, err
		}
		return uint64(vv), nil
	case reflect.Uint64:
		return v, nil
	}
	return 0, RangeError(v, k)
}
