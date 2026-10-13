package rawext

import (
	"reflect"

	"github.com/shamaton/msgpack/v3/ext"
)

// StreamEncoder is the built-in stream ext.StreamEncoder for ext.Raw.
var StreamEncoder = new(streamEncoder)

type streamEncoder struct{}

var _ ext.StreamEncoder = (*streamEncoder)(nil)

// Code returns a placeholder; see encoder.Code for why it is unused.
func (s *streamEncoder) Code() int8 {
	return 0
}

func (s *streamEncoder) Type() reflect.Type {
	return TypeOf
}

func (s *streamEncoder) Write(w ext.StreamWriter, value reflect.Value) error {
	r := value.Interface().(ext.Raw)
	n := len(r.Data)
	if err := validateLength(n); err != nil {
		return err
	}
	code, lengthBytes := headerSize(n)

	if err := w.WriteByte1Int(int(code)); err != nil {
		return err
	}
	switch lengthBytes {
	case 1:
		if err := w.WriteByte1Int(n); err != nil {
			return err
		}
	case 2:
		if err := w.WriteByte2Int(n); err != nil {
			return err
		}
	case 4:
		if err := w.WriteByte4Int(n); err != nil {
			return err
		}
	}
	if err := w.WriteByte1Int64(int64(r.Type)); err != nil {
		return err
	}
	return w.WriteBytes(r.Data)
}
