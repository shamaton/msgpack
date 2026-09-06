package decoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

// issue55Status/issue55Statuses mirror the encode-side fixtures in
// defined_type_ext_test.go: named non-struct types, as commonly used for Go
// "enums", that must round-trip through a registered ext decoder.
type issue55DecStatus int
type issue55DecStatuses []issue55DecStatus

type issue55DecHolder struct {
	Status issue55DecStatus
}

// issue55DecDecoder decodes a fixext1 frame whose payload byte is the raw
// int value of issue55DecStatus.
type issue55DecDecoder struct {
	code       int8
	isTypeHits int
}

func (d *issue55DecDecoder) Code() int8 { return d.code }

func (d *issue55DecDecoder) IsType(offset int, data *[]byte) bool {
	d.isTypeHits++
	if (*data)[offset] != def.Fixext1 {
		return false
	}
	return int8((*data)[offset+1]) == d.code
}

func (d *issue55DecDecoder) AsValue(offset int, _ reflect.Kind, data *[]byte) (interface{}, int, error) {
	return issue55DecStatus((*data)[offset+2]), offset + 3, nil
}

func issue55DecBytes(code int8, v issue55DecStatus) []byte {
	return []byte{def.Fixext1, byte(code), byte(v)}
}

func Test_DecodeExtCoderForNamedNonStructType(t *testing.T) {
	dec := &issue55DecDecoder{code: 41}
	AddExtDecoder(dec)
	defer RemoveExtDecoder(dec)

	t.Run("top level", func(t *testing.T) {
		var got issue55DecStatus
		d := decoder{data: issue55DecBytes(41, 7)}
		_, err := d.decode(reflect.ValueOf(&got).Elem(), 0)
		tu.NoError(t, err)
		tu.Equal(t, got, issue55DecStatus(7))
	})

	t.Run("slice element", func(t *testing.T) {
		data := append([]byte{def.FixArray + 2}, issue55DecBytes(41, 1)...)
		data = append(data, issue55DecBytes(41, 2)...)

		var got []issue55DecStatus
		d := decoder{data: data}
		_, err := d.decode(reflect.ValueOf(&got).Elem(), 0)
		tu.NoError(t, err)
		tu.EqualSlice(t, []byte(convertToBytes(got)), []byte(convertToBytes(issue55DecStatuses{1, 2})))
	})

	t.Run("struct field via map", func(t *testing.T) {
		data := append([]byte{def.FixMap + 1, def.FixStr + byte(len("Status"))}, "Status"...)
		data = append(data, issue55DecBytes(41, 9)...)

		var got issue55DecHolder
		d := decoder{data: data, asArray: false}
		_, err := d.decode(reflect.ValueOf(&got).Elem(), 0)
		tu.NoError(t, err)
		tu.Equal(t, got.Status, issue55DecStatus(9))
	})

	t.Run("no match falls back to kind switch", func(t *testing.T) {
		// A registered decoder whose IsType never matches must not prevent an
		// ordinary (non-ext) value of the same underlying kind from decoding.
		var got issue55DecStatus
		d := decoder{data: []byte{0x05}} // fixint 5, not an ext frame
		_, err := d.decode(reflect.ValueOf(&got).Elem(), 0)
		tu.NoError(t, err)
		tu.Equal(t, got, issue55DecStatus(5))
	})
}

func convertToBytes(v any) []byte {
	rv := reflect.ValueOf(v)
	bs := make([]byte, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		bs[i] = byte(rv.Index(i).Int())
	}
	return bs
}

var _ = ext.Decoder((*issue55DecDecoder)(nil))
