package msgpack_test

import (
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

// TestDecodeHugeDeclaredLength covers 32-bit declared lengths with the high
// bit set. On 32-bit platforms they used to become negative ints, causing
// panics or silently accepting malformed input.
func TestDecodeHugeDeclaredLength(t *testing.T) {
	hi := []byte{0xff, 0xff, 0xff, 0xff}
	half := []byte{0x80, 0x00, 0x00, 0x00}
	quarter := []byte{0x40, 0x00, 0x00, 0x01}
	type S struct{ A int }
	// fixmap(2){ "A": 1, "B": <unknown field> }
	structPrefix := []byte{0x82, 0xa1, 'A', 0x01, 0xa1, 'B'}

	tests := []struct {
		name string
		data []byte
		newV func() interface{}
	}{
		{"array32/iface", append([]byte{def.Array32}, hi...), func() interface{} { return new(interface{}) }},
		{"array32/slice", append([]byte{def.Array32}, hi...), func() interface{} { return new([]int) }},
		{"array32/half", append([]byte{def.Array32}, half...), func() interface{} { return new([]int) }},
		{"map32/iface", append([]byte{def.Map32}, hi...), func() interface{} { return new(interface{}) }},
		{"map32/map", append([]byte{def.Map32}, hi...), func() interface{} { return new(map[int]int) }},
		{"map32/quarter/iface", append([]byte{def.Map32}, quarter...), func() interface{} { return new(interface{}) }},
		{"map32/quarter/map", append([]byte{def.Map32}, quarter...), func() interface{} { return new(map[int]int) }},
		{"str32/string", append([]byte{def.Str32}, hi...), func() interface{} { return new(string) }},
		{"str32/iface", append([]byte{def.Str32}, hi...), func() interface{} { return new(interface{}) }},
		{"str32/half", append([]byte{def.Str32}, half...), func() interface{} { return new(string) }},
		{"bin32/bytes", append([]byte{def.Bin32}, hi...), func() interface{} { return new([]byte) }},
		{"bin32/iface", append([]byte{def.Bin32}, hi...), func() interface{} { return new(interface{}) }},
		{"bin32/half", append([]byte{def.Bin32}, half...), func() interface{} { return new([]byte) }},
		{"skip/str32", append(append([]byte{}, structPrefix...), append([]byte{def.Str32}, hi...)...), func() interface{} { return new(S) }},
		{"skip/str32/half", append(append([]byte{}, structPrefix...), append([]byte{def.Str32}, half...)...), func() interface{} { return new(S) }},
		{"skip/bin32", append(append([]byte{}, structPrefix...), append([]byte{def.Bin32}, hi...)...), func() interface{} { return new(S) }},
		{"skip/array32", append(append([]byte{}, structPrefix...), append([]byte{def.Array32}, hi...)...), func() interface{} { return new(S) }},
		{"skip/map32", append(append([]byte{}, structPrefix...), append([]byte{def.Map32}, hi...)...), func() interface{} { return new(S) }},
		{"skip/ext32", append(append([]byte{}, structPrefix...), append([]byte{def.Ext32}, hi...)...), func() interface{} { return new(S) }},
		{"skip/map32/half", append(append([]byte{}, structPrefix...), append([]byte{def.Map32}, half...)...), func() interface{} { return new(S) }},
	}
	for _, df := range decodeFuncs {
		for _, tc := range tests {
			t.Run(df.name+"/"+tc.name, func(t *testing.T) {
				err := mustNotPanic(t, func() error { return df.fn(tc.data, tc.newV()) })
				if err == nil {
					t.Fatal("expected error")
				}
			})
		}
	}
}
