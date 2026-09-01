package decoding

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/common"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

// issue55Status/issue55Statuses mirror the non-stream decode fixtures: named
// non-struct types, as commonly used for Go "enums", that must round-trip
// through a registered stream ext decoder.
type issue55StreamDecStatus int

type issue55StreamDecHolder struct {
	Status issue55StreamDecStatus
}

type issue55StreamDecDecoder struct {
	code int8
}

func (d *issue55StreamDecDecoder) Code() int8 { return d.code }

func (d *issue55StreamDecDecoder) IsType(code byte, innerType int8, _ int) bool {
	return code == def.Fixext1 && innerType == d.code
}

func (d *issue55StreamDecDecoder) ToValue(_ byte, data []byte, _ reflect.Kind) (interface{}, error) {
	return issue55StreamDecStatus(int8(data[0])), nil
}

func newTestDecoder(data []byte) decoder {
	return decoder{r: bytes.NewReader(data), buf: common.GetBuffer()}
}

func issue55StreamDecBytes(code int8, v issue55StreamDecStatus) []byte {
	return []byte{def.Fixext1, byte(code), byte(v)}
}

func Test_StreamDecodeExtCoderForNamedNonStructType(t *testing.T) {
	dec := &issue55StreamDecDecoder{code: 42}
	AddExtDecoder(dec)
	defer RemoveExtDecoder(dec)

	t.Run("top level", func(t *testing.T) {
		var got issue55StreamDecStatus
		d := newTestDecoder(issue55StreamDecBytes(42, 7))
		defer common.PutBuffer(d.buf)
		err := d.decode(reflect.ValueOf(&got).Elem())
		tu.NoError(t, err)
		tu.Equal(t, got, issue55StreamDecStatus(7))
	})

	t.Run("slice element", func(t *testing.T) {
		data := append([]byte{def.FixArray + 2}, issue55StreamDecBytes(42, 1)...)
		data = append(data, issue55StreamDecBytes(42, 2)...)

		var got []issue55StreamDecStatus
		d := newTestDecoder(data)
		defer common.PutBuffer(d.buf)
		err := d.decode(reflect.ValueOf(&got).Elem())
		tu.NoError(t, err)
		if !reflect.DeepEqual(got, []issue55StreamDecStatus{1, 2}) {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("struct field via map", func(t *testing.T) {
		data := append([]byte{def.FixMap + 1, def.FixStr + byte(len("Status"))}, "Status"...)
		data = append(data, issue55StreamDecBytes(42, 9)...)

		var got issue55StreamDecHolder
		d := newTestDecoder(data)
		d.asArray = false
		defer common.PutBuffer(d.buf)
		err := d.decode(reflect.ValueOf(&got).Elem())
		tu.NoError(t, err)
		tu.Equal(t, got.Status, issue55StreamDecStatus(9))
	})

	t.Run("no match falls back to kind switch", func(t *testing.T) {
		// A registered decoder whose IsType never matches must not prevent an
		// ordinary (non-ext) value of the same underlying kind from decoding.
		var got issue55StreamDecStatus
		d := newTestDecoder([]byte{0x05}) // fixint 5, not an ext frame
		defer common.PutBuffer(d.buf)
		err := d.decode(reflect.ValueOf(&got).Elem())
		tu.NoError(t, err)
		tu.Equal(t, got, issue55StreamDecStatus(5))
	})
}

var _ = ext.StreamDecoder((*issue55StreamDecDecoder)(nil))
