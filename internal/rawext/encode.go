package rawext

import (
	"reflect"

	"github.com/shamaton/msgpack/v3/ext"
)

// Encoder is the built-in non-stream ext.Encoder for ext.Raw.
var Encoder = new(encoder)

type encoder struct {
	ext.EncoderCommon
}

var _ ext.Encoder = (*encoder)(nil)

// Code returns a placeholder; ext.Raw is matched by its Go type, not by a
// fixed ext type code (the payload's own type byte varies per value), so
// Code is not used for dispatch.
func (e *encoder) Code() int8 {
	return 0
}

func (e *encoder) Type() reflect.Type {
	return TypeOf
}

func (e *encoder) CalcByteSize(value reflect.Value) (int, error) {
	r := value.Interface().(ext.Raw)
	return byteSize(r)
}

func (e *encoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	r := value.Interface().(ext.Raw)
	n := len(r.Data)
	code, lengthBytes := headerSize(n)

	offset = e.SetByte1Int(int(code), offset, bytes)
	switch lengthBytes {
	case 1:
		offset = e.SetByte1Int(n, offset, bytes)
	case 2:
		offset = e.SetByte2Int(n, offset, bytes)
	case 4:
		offset = e.SetByte4Int(n, offset, bytes)
	}
	offset = e.SetByte1Int64(int64(r.Type), offset, bytes)
	offset = e.SetBytes(r.Data, offset, bytes)
	return offset
}
