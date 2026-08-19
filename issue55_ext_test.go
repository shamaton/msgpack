package msgpack_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3"
	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

const issue55PublicCode int8 = 40

type issue55PublicStatus int

type issue55PublicHolder struct {
	Status issue55PublicStatus
}

type issue55PublicEncoder struct {
	ext.EncoderCommon
}

func (*issue55PublicEncoder) Code() int8         { return issue55PublicCode }
func (*issue55PublicEncoder) Type() reflect.Type { return reflect.TypeOf(issue55PublicStatus(0)) }
func (*issue55PublicEncoder) CalcByteSize(reflect.Value) (int, error) {
	return 3, nil
}
func (e *issue55PublicEncoder) WriteToBytes(_ reflect.Value, offset int, data *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, data)
	offset = e.SetByte1Int(int(e.Code()), offset, data)
	return e.SetByte1Uint64(0xaa, offset, data)
}

type issue55PublicDecoder struct{}

func (*issue55PublicDecoder) Code() int8               { return issue55PublicCode }
func (*issue55PublicDecoder) IsType(int, *[]byte) bool { return false }
func (*issue55PublicDecoder) AsValue(int, reflect.Kind, *[]byte) (any, int, error) {
	return nil, 0, nil
}

type issue55PublicStreamEncoder struct{}

func (*issue55PublicStreamEncoder) Code() int8         { return issue55PublicCode }
func (*issue55PublicStreamEncoder) Type() reflect.Type { return reflect.TypeOf(issue55PublicStatus(0)) }
func (e *issue55PublicStreamEncoder) Write(w ext.StreamWriter, _ reflect.Value) error {
	if err := w.WriteByte1Int(def.Fixext1); err != nil {
		return err
	}
	if err := w.WriteByte1Int(int(e.Code())); err != nil {
		return err
	}
	return w.WriteByte1Uint64(0xaa)
}

type issue55PublicStreamDecoder struct{}

func (*issue55PublicStreamDecoder) Code() int8 { return issue55PublicCode }
func (*issue55PublicStreamDecoder) IsType(byte, int8, int) bool {
	return false
}
func (*issue55PublicStreamDecoder) ToValue(byte, []byte, reflect.Kind) (any, error) {
	return nil, nil
}

func TestIssue55DefinedTypeExtThroughPublicAPIs(t *testing.T) {
	encoder := &issue55PublicEncoder{}
	decoder := &issue55PublicDecoder{}
	if err := msgpack.AddExtCoder(encoder, decoder); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := msgpack.RemoveExtCoder(encoder, decoder); err != nil {
			t.Error(err)
		}
	})

	streamEncoder := &issue55PublicStreamEncoder{}
	streamDecoder := &issue55PublicStreamDecoder{}
	if err := msgpack.AddExtStreamCoder(streamEncoder, streamDecoder); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := msgpack.RemoveExtStreamCoder(streamEncoder, streamDecoder); err != nil {
			t.Error(err)
		}
	})

	extValue := []byte{def.Fixext1, byte(issue55PublicCode), 0xaa}
	value := issue55PublicHolder{Status: 1}
	marshalCases := []struct {
		name string
		fn   func(any) ([]byte, error)
	}{
		{name: "map", fn: msgpack.MarshalAsMap},
		{name: "array", fn: msgpack.MarshalAsArray},
	}
	for _, tc := range marshalCases {
		t.Run("marshal "+tc.name, func(t *testing.T) {
			encoded, err := tc.fn(value)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(encoded, extValue) {
				t.Fatalf("defined type extension not found in %x", encoded)
			}
		})
	}

	streamCases := []struct {
		name string
		fn   func(*bytes.Buffer, any) error
	}{
		{name: "map", fn: func(buffer *bytes.Buffer, value any) error {
			return msgpack.MarshalWriteAsMap(buffer, value)
		}},
		{name: "array", fn: func(buffer *bytes.Buffer, value any) error {
			return msgpack.MarshalWriteAsArray(buffer, value)
		}},
	}
	for _, tc := range streamCases {
		t.Run("stream "+tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := tc.fn(&buffer, value); err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(buffer.Bytes(), extValue) {
				t.Fatalf("defined type extension not found in %x", buffer.Bytes())
			}
		})
	}
}
