package msgpack_test

import (
	"bytes"
	"testing"

	"github.com/shamaton/msgpack/v2"
)

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := msgpack.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestDecodeIntoInterfaceHoldingPointer(t *testing.T) {
	five := mustMarshal(t, 5)
	nilData := mustMarshal(t, nil)

	t.Run("nonNilPointer", func(t *testing.T) {
		for _, df := range decodeFuncs {
			var d int
			var v interface{} = &d
			if err := mustNotPanic(t, func() error { return df.fn(five, &v) }); err != nil {
				t.Fatalf("%s: %v", df.name, err)
			}
			if d != 5 || v != &d {
				t.Fatalf("%s: got d=%d v=%#v", df.name, d, v)
			}
		}
	})

	// The stream decoder used to read the byte after the value instead.
	t.Run("nonNilPointer/trailingData", func(t *testing.T) {
		var d int
		var v interface{} = &d
		data := append(append([]byte{}, five...), 0x07)
		if err := msgpack.UnmarshalRead(bytes.NewReader(data), &v); err != nil {
			t.Fatal(err)
		}
		if d != 5 {
			t.Fatalf("got %d, want 5", d)
		}
	})

	t.Run("nonNilPointer/nil", func(t *testing.T) {
		for _, df := range decodeFuncs {
			d := 3
			var v interface{} = &d
			if err := mustNotPanic(t, func() error { return df.fn(nilData, &v) }); err != nil {
				t.Fatalf("%s: %v", df.name, err)
			}
			if d != 3 || v != &d {
				t.Fatalf("%s: got d=%d v=%#v", df.name, d, v)
			}
		}
	})

	t.Run("nilPointer", func(t *testing.T) {
		for _, df := range decodeFuncs {
			var v interface{} = (*int)(nil)
			if err := mustNotPanic(t, func() error { return df.fn(five, &v) }); err != nil {
				t.Fatalf("%s: %v", df.name, err)
			}
			want := interface{}(nil)
			if err := msgpack.Unmarshal(five, &want); err != nil {
				t.Fatal(err)
			}
			if v != want {
				t.Fatalf("%s: got %#v, want %#v", df.name, v, want)
			}
		}
	})
}
