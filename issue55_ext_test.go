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

// --- round-trip fixtures (encode AND decode) -------------------------------
//
// The types and coders above only exercise the encoder side (they assert
// that the ext frame appears in the encoded bytes). The fixtures below prove
// the full round trip: Unmarshal must restore a named non-struct type (and a
// slice of it) from the ext frame that Marshal produced.

const issue55RoundTripCode int8 = 43

type issue55RoundTripStatus int

const (
	issue55RoundTripActive issue55RoundTripStatus = 1
	issue55RoundTripClosed issue55RoundTripStatus = 2
)

type issue55RoundTripHolder struct {
	Status   issue55RoundTripStatus
	Statuses []issue55RoundTripStatus
}

type issue55RoundTripEncoder struct {
	ext.EncoderCommon
}

func (*issue55RoundTripEncoder) Code() int8 { return issue55RoundTripCode }
func (*issue55RoundTripEncoder) Type() reflect.Type {
	return reflect.TypeOf(issue55RoundTripStatus(0))
}
func (*issue55RoundTripEncoder) CalcByteSize(reflect.Value) (int, error) {
	return def.Byte1 + def.Byte1 + def.Byte1, nil
}
func (e *issue55RoundTripEncoder) WriteToBytes(value reflect.Value, offset int, data *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, data)
	offset = e.SetByte1Int(int(e.Code()), offset, data)
	return e.SetByte1Int64(value.Int(), offset, data)
}

type issue55RoundTripDecoder struct{}

func (*issue55RoundTripDecoder) Code() int8 { return issue55RoundTripCode }
func (*issue55RoundTripDecoder) IsType(offset int, data *[]byte) bool {
	if (*data)[offset] != def.Fixext1 {
		return false
	}
	return int8((*data)[offset+1]) == issue55RoundTripCode
}
func (*issue55RoundTripDecoder) AsValue(offset int, _ reflect.Kind, data *[]byte) (any, int, error) {
	return issue55RoundTripStatus(int8((*data)[offset+2])), offset + 3, nil
}

// TestIssue55DefinedTypeExtRoundTrip proves the fix for shamaton/msgpack#55
// end to end: registering an ext coder for a named non-struct type must
// produce a lossless Marshal -> Unmarshal round trip, at the top level, as a
// struct field, and inside a slice.
func TestIssue55DefinedTypeExtRoundTrip(t *testing.T) {
	encoder := &issue55RoundTripEncoder{}
	decoder := &issue55RoundTripDecoder{}
	if err := msgpack.AddExtCoder(encoder, decoder); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := msgpack.RemoveExtCoder(encoder, decoder); err != nil {
			t.Error(err)
		}
	})

	t.Run("top level", func(t *testing.T) {
		b, err := msgpack.Marshal(issue55RoundTripActive)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte{def.Fixext1, byte(issue55RoundTripCode), byte(issue55RoundTripActive)}
		if !bytes.Equal(b, want) {
			t.Fatalf("encode mismatch. got % 02x, want % 02x", b, want)
		}

		var got issue55RoundTripStatus
		if err := msgpack.Unmarshal(b, &got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if got != issue55RoundTripActive {
			t.Fatalf("round-trip mismatch. got %d, want %d", got, issue55RoundTripActive)
		}
	})

	t.Run("slice", func(t *testing.T) {
		in := []issue55RoundTripStatus{issue55RoundTripActive, issue55RoundTripClosed}
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		var got []issue55RoundTripStatus
		if err := msgpack.Unmarshal(b, &got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if !reflect.DeepEqual(got, in) {
			t.Fatalf("round-trip mismatch. got %v, want %v", got, in)
		}
	})

	t.Run("struct field and slice field", func(t *testing.T) {
		in := issue55RoundTripHolder{
			Status:   issue55RoundTripClosed,
			Statuses: []issue55RoundTripStatus{issue55RoundTripActive, issue55RoundTripClosed},
		}
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		var got issue55RoundTripHolder
		if err := msgpack.Unmarshal(b, &got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if !reflect.DeepEqual(got, in) {
			t.Fatalf("round-trip mismatch. got %+v, want %+v", got, in)
		}
	})
}

// --- stream round-trip fixtures (encode AND decode via io.Writer/io.Reader) -

const issue55StreamRoundTripCode int8 = 45

type issue55StreamRoundTripStatus int

const (
	issue55StreamRoundTripActive issue55StreamRoundTripStatus = 1
	issue55StreamRoundTripClosed issue55StreamRoundTripStatus = 2
)

type issue55StreamRoundTripHolder struct {
	Status   issue55StreamRoundTripStatus
	Statuses []issue55StreamRoundTripStatus
}

type issue55StreamRoundTripEncoder struct{}

func (*issue55StreamRoundTripEncoder) Code() int8 { return issue55StreamRoundTripCode }
func (*issue55StreamRoundTripEncoder) Type() reflect.Type {
	return reflect.TypeOf(issue55StreamRoundTripStatus(0))
}
func (e *issue55StreamRoundTripEncoder) Write(w ext.StreamWriter, value reflect.Value) error {
	if err := w.WriteByte1Int(def.Fixext1); err != nil {
		return err
	}
	if err := w.WriteByte1Int(int(e.Code())); err != nil {
		return err
	}
	return w.WriteByte1Int64(value.Int())
}

type issue55StreamRoundTripDecoder struct{}

func (*issue55StreamRoundTripDecoder) Code() int8 { return issue55StreamRoundTripCode }
func (*issue55StreamRoundTripDecoder) IsType(code byte, innerType int8, _ int) bool {
	return code == def.Fixext1 && innerType == issue55StreamRoundTripCode
}
func (*issue55StreamRoundTripDecoder) ToValue(_ byte, data []byte, _ reflect.Kind) (any, error) {
	return issue55StreamRoundTripStatus(int8(data[0])), nil
}

// TestIssue55DefinedTypeExtStreamRoundTrip proves the fix for
// shamaton/msgpack#55 through the streaming public APIs: registering a
// stream ext coder for a named non-struct type must produce a lossless
// MarshalWrite -> UnmarshalRead round trip, at the top level, as a struct
// field, and inside a slice -- not just through the byte-slice Marshal API.
func TestIssue55DefinedTypeExtStreamRoundTrip(t *testing.T) {
	encoder := &issue55StreamRoundTripEncoder{}
	decoder := &issue55StreamRoundTripDecoder{}
	if err := msgpack.AddExtStreamCoder(encoder, decoder); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := msgpack.RemoveExtStreamCoder(encoder, decoder); err != nil {
			t.Error(err)
		}
	})

	t.Run("top level", func(t *testing.T) {
		var buf bytes.Buffer
		if err := msgpack.MarshalWrite(&buf, issue55StreamRoundTripActive); err != nil {
			t.Fatal(err)
		}
		want := []byte{def.Fixext1, byte(issue55StreamRoundTripCode), byte(issue55StreamRoundTripActive)}
		if !bytes.Equal(buf.Bytes(), want) {
			t.Fatalf("encode mismatch. got % 02x, want % 02x", buf.Bytes(), want)
		}

		var got issue55StreamRoundTripStatus
		if err := msgpack.UnmarshalRead(&buf, &got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if got != issue55StreamRoundTripActive {
			t.Fatalf("round-trip mismatch. got %d, want %d", got, issue55StreamRoundTripActive)
		}
	})

	t.Run("slice", func(t *testing.T) {
		in := []issue55StreamRoundTripStatus{issue55StreamRoundTripActive, issue55StreamRoundTripClosed}
		var buf bytes.Buffer
		if err := msgpack.MarshalWrite(&buf, in); err != nil {
			t.Fatal(err)
		}

		var got []issue55StreamRoundTripStatus
		if err := msgpack.UnmarshalRead(&buf, &got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if !reflect.DeepEqual(got, in) {
			t.Fatalf("round-trip mismatch. got %v, want %v", got, in)
		}
	})

	t.Run("struct field and slice field", func(t *testing.T) {
		in := issue55StreamRoundTripHolder{
			Status:   issue55StreamRoundTripClosed,
			Statuses: []issue55StreamRoundTripStatus{issue55StreamRoundTripActive, issue55StreamRoundTripClosed},
		}
		var buf bytes.Buffer
		if err := msgpack.MarshalWrite(&buf, in); err != nil {
			t.Fatal(err)
		}

		var got issue55StreamRoundTripHolder
		if err := msgpack.UnmarshalRead(&buf, &got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if !reflect.DeepEqual(got, in) {
			t.Fatalf("round-trip mismatch. got %+v, want %+v", got, in)
		}
	})
}
