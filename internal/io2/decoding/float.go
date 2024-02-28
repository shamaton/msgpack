package decoding

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func asFloat32(r io.Reader, k reflect.Kind) (float32, error) {
	code, err := readSize1(r)
	if err != nil {
		return 0, err
	}
	return asFloat32WithCode(r, code, k)
}

func asFloat32WithCode(r io.Reader, code byte, k reflect.Kind) (float32, error) {
	switch {
	case code == def.Float32:
		bs, err := readSize4(r)
		if err != nil {
			return 0, err
		}
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return v, nil

	case isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			break
		}
		return float32(v), nil

	case isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			break
		}
		return float32(v), nil

	case code == def.Nil:
		return 0, nil
	}
	return 0, errorTemplate(code, k)
}

func asFloat64(r io.Reader, k reflect.Kind) (float64, error) {
	code, err := readSize1(r)
	if err != nil {
		return 0, err
	}
	return asFloat64WithCode(r, code, k)
}

func asFloat64WithCode(r io.Reader, code byte, k reflect.Kind) (float64, error) {
	switch {
	case code == def.Float64:
		bs, err := readSize8(r)
		if err != nil {
			return 0, err
		}
		v := math.Float64frombits(binary.BigEndian.Uint64(bs))
		return v, nil

	case code == def.Float32:
		bs, err := readSize4(r)
		if err != nil {
			return 0, err
		}
		v := math.Float32frombits(binary.BigEndian.Uint32(bs))
		return float64(v), nil

	case isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, err := asUintWithCode(r, code, k)
		if err != nil {
			break
		}
		return float64(v), nil

	case isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, err := asIntWithCode(r, code, k)
		if err != nil {
			break
		}
		return float64(v), nil

	case code == def.Nil:
		return 0, nil
	}
	return 0, errorTemplate(code, k)
}
