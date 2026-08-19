package encoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
	"github.com/shamaton/msgpack/v3/time"
)

func Test_AddExtEncoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		AddExtEncoder(time.Encoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

func Test_RemoveExtEncoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		RemoveExtEncoder(time.Encoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

// enumUint8 is a named non-struct type, as commonly used for "enums" in Go.
type enumUint8 uint8

const enumExtCode = 0x02

type enumUint8Encoder struct {
	ext.EncoderCommon
}

var _ ext.Encoder = (*enumUint8Encoder)(nil)

func (e *enumUint8Encoder) Code() int8 { return enumExtCode }

func (e *enumUint8Encoder) Type() reflect.Type { return reflect.TypeOf(enumUint8(0)) }

func (e *enumUint8Encoder) CalcByteSize(reflect.Value) (int, error) {
	return def.Byte1 + def.Byte1 + def.Byte1, nil
}

func (e *enumUint8Encoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, bytes)
	offset = e.SetByte1Int(int(e.Code()), offset, bytes)
	offset = e.SetByte1Uint64(value.Uint(), offset, bytes)
	return offset
}

// Test_ExtEncoderForNamedNonStructType covers shamaton/msgpack#55: an ext encoder
// registered for a named non-struct type must be used at the top level and inside
// a slice, instead of the plain int encoding.
func Test_ExtEncoderForNamedNonStructType(t *testing.T) {
	enc := &enumUint8Encoder{}
	AddExtEncoder(enc)
	defer RemoveExtEncoder(enc)

	t.Run("top level", func(t *testing.T) {
		b, err := Encode(enumUint8(1), false)
		tu.NoError(t, err)
		tu.EqualSlice(t, b, []byte{def.Fixext1, enumExtCode, 0x01})
	})

	t.Run("slice", func(t *testing.T) {
		b, err := Encode([]enumUint8{1, 2}, false)
		tu.NoError(t, err)
		tu.EqualSlice(t, b, []byte{
			def.FixArray + 2,
			def.Fixext1, enumExtCode, 0x01,
			def.Fixext1, enumExtCode, 0x02,
		})
	})
}
