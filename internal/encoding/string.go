package encoding

import (
	"io"
	"math"
	"unsafe"

	"github.com/shamaton/msgpack/v2/def"
)

func (e *encoder) calcString(v string) int {
	// NOTE : unsafe
	strBytes := *(*[]byte)(unsafe.Pointer(&v))
	l := len(strBytes)
	if l < 32 {
		return l
	} else if l <= math.MaxUint8 {
		return def.Byte1 + l
	} else if l <= math.MaxUint16 {
		return def.Byte2 + l
	}
	return def.Byte4 + l
	// NOTE : length over uint32
}

func (e *encoder) writeString(str string, writer io.Writer) (err error) {
	l := len(str)
	if l < 32 {
		err = e.setByte1Int(def.FixStr+l, writer)
		if err != nil {
			return err
		}
	} else if l <= math.MaxUint8 {
		err = e.setByte1Int(def.Str8, writer)
		if err != nil {
			return err
		}

		err = e.setByte1Int(l, writer)
		if err != nil {
			return err
		}
	} else if l <= math.MaxUint16 {
		err = e.setByte1Int(def.Str16, writer)
		if err != nil {
			return err
		}

		err = e.setByte2Int(l, writer)
		if err != nil {
			return err
		}
	} else {
		err = e.setByte1Int(def.Str32, writer)
		if err != nil {
			return err
		}

		err = e.setByte4Int(l, writer)
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(writer, str)
	return err
}
