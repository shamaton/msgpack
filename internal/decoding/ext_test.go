package decoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
	"github.com/shamaton/msgpack/v2/time"
)

func Test_AddExtDecoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		AddExtDecoder(time.Decoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

func Test_RemoveExtDecoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		RemoveExtDecoder(time.Decoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

type trackingExtDecoder struct {
	isTypeCalls  int
	asValueCalls int
}

func (td *trackingExtDecoder) Code() int8 {
	return 42
}

func (td *trackingExtDecoder) IsType(_ int, _ *[]byte) bool {
	td.isTypeCalls++
	return false
}

func (td *trackingExtDecoder) AsValue(_ int, _ reflect.Kind, _ *[]byte) (interface{}, int, error) {
	td.asValueCalls++
	return nil, 0, ErrTestExtDecoder
}

func TestExtValidationRejectsTruncatedBytesBeforeCustomDecoders(t *testing.T) {
	dec := &trackingExtDecoder{}
	AddExtDecoder(dec)
	defer RemoveExtDecoder(dec)

	testcases := []struct {
		name string
		data []byte
	}{
		{name: "fixext1", data: []byte{def.Fixext1, 42}},
		{name: "fixext2", data: []byte{def.Fixext2, 42, 0}},
		{name: "fixext4", data: []byte{def.Fixext4, 42, 0, 0, 0}},
		{name: "fixext8", data: []byte{def.Fixext8, 42, 0, 0, 0, 0, 0, 0, 0}},
		{name: "fixext16", data: []byte{def.Fixext16, 42, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "ext8", data: []byte{def.Ext8, 1, 42}},
		{name: "ext16", data: []byte{def.Ext16, 0, 1, 42}},
		{name: "ext32", data: []byte{def.Ext32, 0, 0, 0, 1, 42}},
	}

	for _, tc := range testcases {
		t.Run("interface/"+tc.name, func(t *testing.T) {
			dec.isTypeCalls = 0
			dec.asValueCalls = 0

			d := decoder{data: tc.data}
			_, _, err := d.asInterface(0, reflect.Interface)
			tu.IsError(t, err, def.ErrTooShortBytes)
			tu.Equal(t, dec.isTypeCalls, 0)
			tu.Equal(t, dec.asValueCalls, 0)
		})

		t.Run("struct/"+tc.name, func(t *testing.T) {
			dec.isTypeCalls = 0
			dec.asValueCalls = 0

			d := decoder{data: tc.data}
			var v struct{}
			_, err := d.setStruct(reflect.ValueOf(&v).Elem(), 0, reflect.Struct)
			tu.IsError(t, err, def.ErrTooShortBytes)
			tu.Equal(t, dec.isTypeCalls, 0)
			tu.Equal(t, dec.asValueCalls, 0)
		})
	}
}
