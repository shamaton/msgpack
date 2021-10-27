package decoding

import (
	"bufio"
	"encoding/binary"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) asFloat32(reader *bufio.Reader, k reflect.Kind) (float32, error) {
	code, err := peekCode(reader)
	if err != nil {
		return 0, err
	}

	switch {
	case code == def.Float32:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize4(reader)
		if err != nil {
			return 0, err
		}
		v := math.Float32frombits(binary.BigEndian.Uint32(bs[:]))
		return v, nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, err := d.asUint(reader, k)
		if err != nil {
			break
		}
		return float32(v), nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, err := d.asInt(reader, k)
		if err != nil {
			break
		}
		return float32(v), nil

	case code == def.Nil:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		return 0, nil
	}

	return 0, d.errorTemplate(code, k)
}

func (d *decoder) asFloat64(reader *bufio.Reader, k reflect.Kind) (float64, error) {
	code, err := peekCode(reader)
	if err != nil {
		return 0, err
	}

	switch {
	case code == def.Float64:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize8(reader)
		if err != nil {
			return 0, err
		}
		v := math.Float64frombits(binary.BigEndian.Uint64(bs[:]))
		return v, nil

	case code == def.Float32:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		bs, err := d.readSize4(reader)
		if err != nil {
			return 0, err
		}
		v := math.Float32frombits(binary.BigEndian.Uint32(bs[:]))
		return float64(v), nil

	case d.isPositiveFixNum(code), code == def.Uint8, code == def.Uint16, code == def.Uint32, code == def.Uint64:
		v, err := d.asUint(reader, k)
		if err != nil {
			break
		}
		return float64(v), nil

	case d.isNegativeFixNum(code), code == def.Int8, code == def.Int16, code == def.Int32, code == def.Int64:
		v, err := d.asInt(reader, k)
		if err != nil {
			break
		}
		return float64(v), nil

	case code == def.Nil:
		err = skipOne(reader)
		if err != nil {
			return 0, err
		}

		return 0, nil
	}

	return 0, d.errorTemplate(code, k)
}
