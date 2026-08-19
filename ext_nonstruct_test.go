package msgpack_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3"
	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

// Role is a named non-struct type, as commonly used for "enums" in Go.
// See shamaton/msgpack#55: ext coders registered for such types were ignored.
type Role uint8

const (
	roleUser  Role = 1
	roleAdmin Role = 2

	roleExtCode = 0x02
)

type roleEncoder struct {
	ext.EncoderCommon
}

var _ ext.Encoder = (*roleEncoder)(nil)

func (e *roleEncoder) Code() int8 { return roleExtCode }

func (e *roleEncoder) Type() reflect.Type { return reflect.TypeOf(Role(0)) }

func (e *roleEncoder) CalcByteSize(reflect.Value) (int, error) {
	return def.Byte1 + def.Byte1 + def.Byte1, nil
}

func (e *roleEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, bytes)
	offset = e.SetByte1Int(int(e.Code()), offset, bytes)
	offset = e.SetByte1Uint64(value.Uint(), offset, bytes)
	return offset
}

type roleDecoder struct {
	ext.DecoderCommon
}

var _ ext.Decoder = (*roleDecoder)(nil)

func (d *roleDecoder) Code() int8 { return roleExtCode }

func (d *roleDecoder) IsType(offset int, data *[]byte) bool {
	code, offset := d.ReadSize1(offset, data)
	if code == def.Fixext1 {
		typ, _ := d.ReadSize1(offset, data)
		return int8(typ) == d.Code()
	}
	return false
}

func (d *roleDecoder) AsValue(offset int, k reflect.Kind, data *[]byte) (interface{}, int, error) {
	code, offset := d.ReadSize1(offset, data)
	if code == def.Fixext1 {
		_, offset = d.ReadSize1(offset, data) // type code
		bs, offset := d.ReadSizeN(offset, def.Byte1, data)
		return Role(bs[0]), offset, nil
	}
	return Role(0), 0, fmt.Errorf("unexpected code %x decoding as %v", code, k)
}

// TestExtCoderForNamedNonStructType covers shamaton/msgpack#55 end to end: an ext
// coder registered for a named non-struct type must be used both when encoding and
// when decoding, at the top level and inside a slice, so the round-trip is lossless.
func TestExtCoderForNamedNonStructType(t *testing.T) {
	if err := msgpack.AddExtCoder(&roleEncoder{}, &roleDecoder{}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := msgpack.RemoveExtCoder(&roleEncoder{}, &roleDecoder{}); err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("top level uses the ext frame", func(t *testing.T) {
		b, err := msgpack.Marshal(roleUser)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{def.Fixext1, roleExtCode, byte(roleUser)}
		if !reflect.DeepEqual(b, want) {
			t.Fatalf("encode mismatch. got % 02x, want % 02x", b, want)
		}

		var got Role
		if err := msgpack.Unmarshal(b, &got); err != nil {
			t.Fatal(err)
		}
		if got != roleUser {
			t.Fatalf("round-trip mismatch. got %d, want %d", got, roleUser)
		}
	})

	t.Run("slice uses ext frames", func(t *testing.T) {
		in := []Role{roleUser, roleAdmin}
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{
			def.FixArray + 2,
			def.Fixext1, roleExtCode, byte(roleUser),
			def.Fixext1, roleExtCode, byte(roleAdmin),
		}
		if !reflect.DeepEqual(b, want) {
			t.Fatalf("encode mismatch. got % 02x, want % 02x", b, want)
		}

		var got []Role
		if err := msgpack.Unmarshal(b, &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, in) {
			t.Fatalf("round-trip mismatch. got %v, want %v", got, in)
		}
	})
}
