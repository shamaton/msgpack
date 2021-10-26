package encoding

import (
	"io"
	"math"

	"github.com/shamaton/msgpack/v2/def"
)

func (e *encoder) calcUint(v uint64) int {
	if v <= math.MaxInt8 {
		// format code only
		return 0
	} else if v <= math.MaxUint8 {
		return def.Byte1
	} else if v <= math.MaxUint16 {
		return def.Byte2
	} else if v <= math.MaxUint32 {
		return def.Byte4
	}
	return def.Byte8
}

func (e *encoder) writeUint(v uint64, writer io.Writer) error {
	if v <= math.MaxInt8 {
		return e.setByte1Uint64(v, writer)
	}

	if v <= math.MaxUint8 {
		err := e.setByte1Int(def.Uint8, writer)
		if err != nil {
			return err
		}

		return e.setByte1Uint64(v, writer)
	}

	if v <= math.MaxUint16 {
		err := e.setByte1Int(def.Uint16, writer)
		if err != nil {
			return err
		}

		return e.setByte2Uint64(v, writer)
	}

	if v <= math.MaxUint32 {
		err := e.setByte1Int(def.Uint32, writer)
		if err != nil {
			return err
		}

		return e.setByte4Uint64(v, writer)
	}

	err := e.setByte1Int(def.Uint64, writer)
	if err != nil {
		return err
	}

	return e.setByte8Uint64(v, writer)
}
