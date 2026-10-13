package msgpack_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3"
	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

// TestUnmarshalUnregisteredExtFallsBackToRaw verifies that ext payloads with
// no registered Decoder are decoded into ext.Raw instead of failing, per the
// MessagePack spec (unregistered ext types must still be parsable).
func TestUnmarshalUnregisteredExtFallsBackToRaw(t *testing.T) {
	// fixext1: type=5, payload=[0x7b]
	data := []byte{0xd4, 0x05, 0x7b}

	methods := []struct {
		name string
		fn   func([]byte, interface{}) error
	}{
		{name: "Unmarshal", fn: msgpack.Unmarshal},
		{name: "UnmarshalAsArray", fn: msgpack.UnmarshalAsArray},
		{name: "UnmarshalAsMap", fn: msgpack.UnmarshalAsMap},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			var v interface{}
			if err := method.fn(data, &v); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := ext.Raw{Type: 5, Data: []byte{0x7b}}
			got, ok := v.(ext.Raw)
			if !ok {
				t.Fatalf("expected ext.Raw, got %T (%v)", v, v)
			}
			if got.Type != want.Type || !reflect.DeepEqual(got.Data, want.Data) {
				t.Fatalf("expected %+v, got %+v", want, got)
			}
		})
	}
}

// TestMarshalExtRawRoundTrip verifies that ext.Raw values obtained from
// decoding an unregistered ext payload can be marshaled back to the exact
// same bytes when the original encoding already used the canonical
// (shortest) ext wire format for its payload length.
func TestMarshalExtRawRoundTrip(t *testing.T) {
	testcases := []struct {
		name string
		data []byte
	}{
		{name: "fixext1", data: []byte{0xd4, 0x05, 0x7b}},
		{name: "fixext2", data: []byte{0xd5, 0x05, 0x01, 0x02}},
		{name: "fixext4", data: []byte{0xd6, 0x05, 0x01, 0x02, 0x03, 0x04}},
		{name: "fixext8", data: append([]byte{0xd7, 0x05}, make([]byte, 8)...)},
		{name: "fixext16", data: append([]byte{0xd8, 0x05}, make([]byte, 16)...)},
		{name: "ext8", data: append([]byte{0xc7, 20, 0x05}, make([]byte, 20)...)},
		{name: "ext8.empty", data: []byte{0xc7, 0, 0x05}},
		{name: "ext16", data: append([]byte{0xc8, 0x01, 0x2c, 0x05}, make([]byte, 300)...)},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var v interface{}
			if err := msgpack.Unmarshal(tc.data, &v); err != nil {
				t.Fatalf("unexpected decode error: %v", err)
			}
			raw, ok := v.(ext.Raw)
			if !ok {
				t.Fatalf("expected ext.Raw, got %T", v)
			}

			out, err := msgpack.Marshal(raw)
			if err != nil {
				t.Fatalf("unexpected encode error: %v", err)
			}
			if !reflect.DeepEqual(out, tc.data) {
				t.Fatalf("round-trip mismatch:\n got: %#v\nwant: %#v", out, tc.data)
			}

			var v2 interface{}
			if err := msgpack.Unmarshal(out, &v2); err != nil {
				t.Fatalf("unexpected re-decode error: %v", err)
			}
			if !reflect.DeepEqual(v, v2) {
				t.Fatalf("re-decoded value mismatch: %+v vs %+v", v, v2)
			}
		})
	}
}

// TestMarshalExtRawNonCanonicalPreservesValue verifies that re-encoding an
// ext.Raw decoded from a non-canonical wire form (one that does not use the
// shortest possible ext format for its payload length) preserves Type and
// Data, even though the exact wire bytes differ: the encoder always picks
// the shortest matching ext format for the payload length.
func TestMarshalExtRawNonCanonicalPreservesValue(t *testing.T) {
	// ext8 with a 1-byte payload; the canonical form for a 1-byte payload is
	// fixext1, so this round-trip changes the wire bytes but not the value.
	data := []byte{0xc7, 1, 0x05, 0x7b}

	var v interface{}
	if err := msgpack.Unmarshal(data, &v); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	raw, ok := v.(ext.Raw)
	if !ok {
		t.Fatalf("expected ext.Raw, got %T", v)
	}

	out, err := msgpack.Marshal(raw)
	if err != nil {
		t.Fatalf("unexpected encode error: %v", err)
	}
	wantCanonical := []byte{0xd4, 0x05, 0x7b}
	if !reflect.DeepEqual(out, wantCanonical) {
		t.Fatalf("expected canonical fixext1 encoding %#v, got %#v", wantCanonical, out)
	}

	var v2 interface{}
	if err := msgpack.Unmarshal(out, &v2); err != nil {
		t.Fatalf("unexpected re-decode error: %v", err)
	}
	if !reflect.DeepEqual(v, v2) {
		t.Fatalf("expected re-decoded value to equal original: %+v vs %+v", v, v2)
	}
}

// TestMarshalExtRawNested verifies ext.Raw encodes correctly when nested
// inside slices, maps, and structs, not just as a top-level value.
func TestMarshalExtRawNested(t *testing.T) {
	raw := ext.Raw{Type: 7, Data: []byte{1, 2, 3}}

	t.Run("slice", func(t *testing.T) {
		in := []ext.Raw{raw, raw}
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatalf("unexpected encode error: %v", err)
		}
		var out []ext.Raw
		if err := msgpack.Unmarshal(b, &out); err != nil {
			t.Fatalf("unexpected decode error: %v", err)
		}
		if !reflect.DeepEqual(in, out) {
			t.Fatalf("expected %+v, got %+v", in, out)
		}
	})

	t.Run("map", func(t *testing.T) {
		in := map[string]ext.Raw{"a": raw}
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatalf("unexpected encode error: %v", err)
		}
		var out map[string]ext.Raw
		if err := msgpack.Unmarshal(b, &out); err != nil {
			t.Fatalf("unexpected decode error: %v", err)
		}
		if !reflect.DeepEqual(in, out) {
			t.Fatalf("expected %+v, got %+v", in, out)
		}
	})

	t.Run("struct field", func(t *testing.T) {
		type S struct {
			R ext.Raw
		}
		in := S{R: raw}
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatalf("unexpected encode error: %v", err)
		}
		var out S
		if err := msgpack.Unmarshal(b, &out); err != nil {
			t.Fatalf("unexpected decode error: %v", err)
		}
		if !reflect.DeepEqual(in, out) {
			t.Fatalf("expected %+v, got %+v", in, out)
		}
	})

	t.Run("interface", func(t *testing.T) {
		var in interface{} = raw
		b, err := msgpack.Marshal(in)
		if err != nil {
			t.Fatalf("unexpected encode error: %v", err)
		}
		var out interface{}
		if err := msgpack.Unmarshal(b, &out); err != nil {
			t.Fatalf("unexpected decode error: %v", err)
		}
		if !reflect.DeepEqual(in, out) {
			t.Fatalf("expected %+v, got %+v", in, out)
		}
	})
}

// TestUnmarshalExtRawAsMapKeyRejected verifies that an unregistered ext value
// used as a MessagePack map key is rejected with a clear error instead of
// panicking: ext.Raw embeds a []byte field, so it is not a comparable Go
// type and cannot be used as a map[interface{}]interface{} key.
func TestUnmarshalExtRawAsMapKeyRejected(t *testing.T) {
	// fixmap(1){ fixext1(type=5,[0x7b]): 1 }
	data := []byte{0x81, 0xd4, 0x05, 0x7b, 0x01}

	t.Run("Unmarshal", func(t *testing.T) {
		var v interface{}
		err := msgpack.Unmarshal(data, &v)
		if !errors.Is(err, def.ErrCanNotSetUncomparableAsMapKey) {
			t.Fatalf("expected %v, got %v", def.ErrCanNotSetUncomparableAsMapKey, err)
		}
	})

	t.Run("UnmarshalRead", func(t *testing.T) {
		var v interface{}
		err := msgpack.UnmarshalRead(bytes.NewReader(data), &v)
		if !errors.Is(err, def.ErrCanNotSetUncomparableAsMapKey) {
			t.Fatalf("expected %v, got %v", def.ErrCanNotSetUncomparableAsMapKey, err)
		}
	})
}

// TestUnmarshalExtRawAsTypedMapKeyRejected verifies the same rejection as
// TestUnmarshalExtRawAsMapKeyRejected, but through the typed-destination
// decode path (decoding directly into map[interface{}]interface{} rather
// than via the interface{} fallback), which constructs the map using
// reflect.Value.SetMapIndex separately from the interface{} decode path.
func TestUnmarshalExtRawAsTypedMapKeyRejected(t *testing.T) {
	// fixmap(1){ fixext1(type=5,[0x7b]): 1 }
	data := []byte{0x81, 0xd4, 0x05, 0x7b, 0x01}

	t.Run("Unmarshal", func(t *testing.T) {
		var v map[interface{}]interface{}
		err := msgpack.Unmarshal(data, &v)
		if !errors.Is(err, def.ErrCanNotSetUncomparableAsMapKey) {
			t.Fatalf("expected %v, got %v", def.ErrCanNotSetUncomparableAsMapKey, err)
		}
	})

	t.Run("UnmarshalRead", func(t *testing.T) {
		var v map[interface{}]interface{}
		err := msgpack.UnmarshalRead(bytes.NewReader(data), &v)
		if !errors.Is(err, def.ErrCanNotSetUncomparableAsMapKey) {
			t.Fatalf("expected %v, got %v", def.ErrCanNotSetUncomparableAsMapKey, err)
		}
	})
}

// TestUnmarshalReadExtRawConsecutiveValuesDoNotAlias verifies that ext.Raw
// values decoded via the stream decoder (UnmarshalRead) each own independent
// backing arrays. The stream decoder's readSize2/4/8/16/N helpers return
// slices aliasing reusable internal buffers, so without a defensive copy,
// decoding several ext values that use those code paths (fixext2 and larger,
// unlike fixext1 which already allocates a fresh 1-byte slice) would corrupt
// earlier results once a later value is read into the same backing array.
func TestUnmarshalReadExtRawConsecutiveValuesDoNotAlias(t *testing.T) {
	// array of 3 fixext2 values, each with a distinct type/payload.
	data := []byte{
		0x93,
		0xd5, 0x01, 0x11, 0x11,
		0xd5, 0x02, 0x22, 0x22,
		0xd5, 0x03, 0x33, 0x33,
	}

	var v []interface{}
	if err := msgpack.UnmarshalRead(bytes.NewReader(data), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []interface{}{
		ext.Raw{Type: 1, Data: []byte{0x11, 0x11}},
		ext.Raw{Type: 2, Data: []byte{0x22, 0x22}},
		ext.Raw{Type: 3, Data: []byte{0x33, 0x33}},
	}
	if !reflect.DeepEqual(v, want) {
		t.Fatalf("expected %+v, got %+v", want, v)
	}
}

// TestUnmarshalReadExtRawIntoTypedDestinationDoesNotAlias verifies the same
// buffer-ownership guarantee as above, but through the typed-destination
// decode path (tryExtDecode), which constructs ext.Raw separately from the
// interface{} fallback path.
func TestUnmarshalReadExtRawIntoTypedDestinationDoesNotAlias(t *testing.T) {
	data := []byte{
		0x92,
		0xd5, 0x01, 0x11, 0x11,
		0xd5, 0x02, 0x22, 0x22,
	}

	var v []ext.Raw
	if err := msgpack.UnmarshalRead(bytes.NewReader(data), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []ext.Raw{
		{Type: 1, Data: []byte{0x11, 0x11}},
		{Type: 2, Data: []byte{0x22, 0x22}},
	}
	if !reflect.DeepEqual(v, want) {
		t.Fatalf("expected %+v, got %+v", want, v)
	}
}

// TestUnmarshalExtRawMalformedExt32LengthReturnsCleanError verifies that an
// ext32 header declaring a huge payload length (far larger than the actual
// remaining input) is rejected with def.ErrTooShortBytes instead of
// panicking. ext32's 4-byte length prefix can declare up to math.MaxUint32,
// which overflows/goes negative when converted directly to a 32-bit `int`
// (relevant on 32-bit platforms); the declared length must be validated
// against the remaining buffer size before any such conversion.
func TestUnmarshalExtRawMalformedExt32LengthReturnsCleanError(t *testing.T) {
	data := []byte{
		def.Ext32, 0xff, 0xff, 0xff, 0xff, // declared length = math.MaxUint32
		0x05, 0x7b, // type byte + 1 payload byte actually present
	}

	t.Run("non-stream", func(t *testing.T) {
		var v ext.Raw
		err := msgpack.Unmarshal(data, &v)
		if !errors.Is(err, def.ErrTooShortBytes) {
			t.Fatalf("expected ErrTooShortBytes, got %v", err)
		}
	})

	t.Run("stream", func(t *testing.T) {
		var v ext.Raw
		err := msgpack.UnmarshalRead(bytes.NewReader(data), &v)
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
	})
}
