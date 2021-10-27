package decoding

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) asComplex64(reader *bufio.Reader, k reflect.Kind) (complex64, error) {
	code, err := reader.ReadByte()
	if err != nil {
		return complex(0, 0), err
	}

	switch code {
	case def.Fixext8:
		t, err := d.readSize1(reader)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext8. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := d.readSize4(reader)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := d.readSize4(reader)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float32frombits(binary.BigEndian.Uint32(rb[:]))
		i := math.Float32frombits(binary.BigEndian.Uint32(ib[:]))
		return complex(r, i), nil

	case def.Fixext16:
		t, err := d.readSize1(reader)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext16. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := d.readSize8(reader)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := d.readSize8(reader)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float64frombits(binary.BigEndian.Uint64(rb[:]))
		i := math.Float64frombits(binary.BigEndian.Uint64(ib[:]))
		return complex64(complex(r, i)), nil

	}

	return complex(0, 0), fmt.Errorf("should not reach this line!! code %x decoding %v", code, k)
}

func (d *decoder) asComplex128(reader *bufio.Reader, k reflect.Kind) (complex128, error) {
	code, err := reader.ReadByte()
	if err != nil {
		return complex(0, 0), err
	}

	switch code {
	case def.Fixext8:
		t, err := d.readSize1(reader)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext8. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := d.readSize4(reader)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := d.readSize4(reader)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float32frombits(binary.BigEndian.Uint32(rb[:]))
		i := math.Float32frombits(binary.BigEndian.Uint32(ib[:]))
		return complex128(complex(r, i)), nil

	case def.Fixext16:
		t, err := d.readSize1(reader)
		if err != nil {
			return complex(0, 0), err
		}
		if int8(t) != def.ComplexTypeCode() {
			return complex(0, 0), fmt.Errorf("fixext16. complex type is diffrent %d, %d", t, def.ComplexTypeCode())
		}
		rb, err := d.readSize8(reader)
		if err != nil {
			return complex(0, 0), err
		}
		ib, err := d.readSize8(reader)
		if err != nil {
			return complex(0, 0), err
		}
		r := math.Float64frombits(binary.BigEndian.Uint64(rb[:]))
		i := math.Float64frombits(binary.BigEndian.Uint64(ib[:]))
		return complex(r, i), nil

	}

	return complex(0, 0), fmt.Errorf("should not reach this line!! code %x decoding %v", code, k)
}
