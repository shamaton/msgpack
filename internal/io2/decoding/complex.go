package decoding

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func asComplex64(r io.Reader, code byte, k reflect.Kind) (complex64, error) {
	switch code {
	case def.Fixext8:
		t, err := readSize1(r)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext8. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := readSize4(r)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := readSize4(r)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float32frombits(binary.BigEndian.Uint32(rb))
		i := math.Float32frombits(binary.BigEndian.Uint32(ib))
		return complex(r, i), nil

	case def.Fixext16:
		t, err := readSize1(r)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext16. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := readSize8(r)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := readSize8(r)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float64frombits(binary.BigEndian.Uint64(rb))
		i := math.Float64frombits(binary.BigEndian.Uint64(ib))
		return complex64(complex(r, i)), nil

	}

	return complex(0, 0), fmt.Errorf("should not reach this line!! code %x decoding %v", code, k)
}

func asComplex128(r io.Reader, code byte, k reflect.Kind) (complex128, error) {
	switch code {
	case def.Fixext8:
		t, err := readSize1(r)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext8. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := readSize4(r)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := readSize4(r)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float32frombits(binary.BigEndian.Uint32(rb))
		i := math.Float32frombits(binary.BigEndian.Uint32(ib))
		return complex128(complex(r, i)), nil

	case def.Fixext16:
		t, err := readSize1(r)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext16. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := readSize8(r)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := readSize8(r)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float64frombits(binary.BigEndian.Uint64(rb))
		i := math.Float64frombits(binary.BigEndian.Uint64(ib))
		return complex(r, i), nil

	}

	return complex(0, 0), fmt.Errorf("should not reach this line!! code %x decoding %v", code, k)
}
