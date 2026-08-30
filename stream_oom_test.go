package msgpack_test

import (
	"bytes"
	"io"
	"math"
	"runtime"
	"testing"

	"github.com/shamaton/msgpack/v3"
)

// TestStreamDecoderBoundedAllocation ensures that declared Array/Map/Bin/Str/Ext
// lengths never translate into allocations proportional to the declaration.
// See GHSA-hfpf-8g2f-p4mx: a 5-byte message used to over-commit up to ~68 GB.
func TestStreamDecoderBoundedAllocation(t *testing.T) {
	const maxAlloc = 1 << 20 // 1 MiB

	tests := []struct {
		name string
		data []byte
		v    interface{}
	}{
		{"Array32 into interface", []byte{0xdd, 0xff, 0xff, 0xff, 0xff}, new(interface{})},
		{"Map32 into interface", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(interface{})},
		{"Bin32 into interface", []byte{0xc6, 0x0b, 0xeb, 0xc2, 0x00}, new(interface{})},
		{"Str32 into interface", []byte{0xdb, 0xff, 0xff, 0xff, 0xff}, new(interface{})},
		{"Ext32 into interface", []byte{0xc9, 0xff, 0xff, 0xff, 0xff, 0x00}, new(interface{})},
		{"Array32 into slice", []byte{0xdd, 0xff, 0xff, 0xff, 0xff}, new([]int)},
		{"Str32 into string", []byte{0xdb, 0xff, 0xff, 0xff, 0xff}, new(string)},
		{"Bin32 into bytes", []byte{0xc6, 0xff, 0xff, 0xff, 0xff}, new([]byte)},
		{"Map32 into map", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(map[string]int)},

		// typed map without a fixed fast path: exercises the generic
		// reflect.MakeMapWithSize allocation in decoding.go
		{"Map32 into typed map (generic)", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(map[int]int)},
		{"Map32 into map[string]interface{}", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(map[string]interface{})},
		// structs are decoded from a map header by default
		{"Map32 into struct", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(streamOomTarget)},
		// fixed-size arrays reject declared counts larger than the array
		{"Array32 into fixed-size array", []byte{0xdd, 0xff, 0xff, 0xff, 0xff}, new([16]int)},
		// element types without a fast path go through reflect.MakeSlice
		{"Array32 into nested slice", []byte{0xdd, 0xff, 0xff, 0xff, 0xff}, new([][]int)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)

			err := msgpack.UnmarshalRead(bytes.NewReader(tc.data), tc.v)

			runtime.ReadMemStats(&after)
			if err == nil {
				t.Fatal("expected error for truncated payload")
			}
			if allocated := after.TotalAlloc - before.TotalAlloc; allocated > maxAlloc {
				t.Fatalf("allocated %d bytes from a %d-byte message (limit %d)",
					allocated, len(tc.data), maxAlloc)
			}
		})
	}
}

// streamOomTarget is decoded from a map header (default struct format).
type streamOomTarget struct {
	A int
	B string
}

func TestUnmarshalReadLargeValidPayload(t *testing.T) {
	want := make([]int, 10000)
	for i := range want {
		want[i] = i
	}
	data, err := msgpack.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	var got []int
	if err := msgpack.UnmarshalRead(bytes.NewReader(data), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestUnmarshalReadDeclaredLengthMatchesData(t *testing.T) {
	payload := bytes.Repeat([]byte{0xab}, 200000)
	data := append([]byte{0xc6, 0x00, 0x03, 0x0d, 0x40}, payload...) // bin32

	var got []byte
	if err := msgpack.UnmarshalRead(bytes.NewReader(data), &got); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("decoded %d bytes, want %d", len(got), len(payload))
	}
}

func TestUnmarshalReadShortReader(t *testing.T) {
	data := append([]byte{0xc6, 0x00, 0x00, 0x01, 0x00}, make([]byte, math.MaxUint16+1)...)
	var got []byte
	if err := msgpack.UnmarshalRead(bytes.NewReader(data[:10]), &got); err == nil {
		t.Fatal("expected error when reader ends before declared length")
	}
}

// oneByteReader wraps a reader and returns at most 1 byte per Read call,
// exercising the short-read completion loop in readSize*.
type oneByteReader struct{ r io.Reader }

func (o oneByteReader) Read(p []byte) (int, error) {
	if len(p) > 1 {
		p = p[:1]
	}
	return o.r.Read(p)
}

// TestUnmarshalReadOneByteReader verifies that a valid Bin32 payload
// decodes correctly when the reader delivers data one byte at a time.
// This is the regression test for the short-read correctness issue
// identified in the review: io.Reader.Read may return fewer bytes than
// requested without an error.
func TestUnmarshalReadOneByteReader(t *testing.T) {
	// Build a valid Bin32 payload (256 bytes of 0xab)
	payload := bytes.Repeat([]byte{0xab}, 256)
	data := append([]byte{0xc6, 0x00, 0x00, 0x01, 0x00}, payload...) // bin32, len=256

	var got []byte
	if err := msgpack.UnmarshalRead(oneByteReader{bytes.NewReader(data)}, &got); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("decoded %d bytes, want %d", len(got), len(payload))
	}
}

// TestUnmarshalReadUnrepresentable32BitLengths verifies that 32-bit lengths
// which are not representable as a non-negative int on 32-bit platforms
// (>= 0x80000000) always return an error and never panic or silently succeed.
// On 64-bit hosts these are valid but oversized; on 386 they would become
// negative after the uint32->int conversion and cause wrong results.
func TestUnmarshalReadUnrepresentable32BitLengths(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		v    interface{}
	}{
		{"Array32 0xffffffff -> interface", []byte{0xdd, 0xff, 0xff, 0xff, 0xff}, new(interface{})},
		{"Map32 0xffffffff -> map[string]int", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(map[string]int)},
		{"Map32 0xffffffff -> struct", []byte{0xdf, 0xff, 0xff, 0xff, 0xff}, new(streamOomTarget)},
		{"Str32 0xffffffff -> string", []byte{0xdb, 0xff, 0xff, 0xff, 0xff}, new(string)},
		{"Bin32 0xffffffff -> []byte", []byte{0xc6, 0xff, 0xff, 0xff, 0xff}, new([]byte)},
		{"Ext32 0xffffffff -> interface", []byte{0xc9, 0xff, 0xff, 0xff, 0xff, 0x00}, new(interface{})},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := msgpack.UnmarshalRead(bytes.NewReader(tc.data), tc.v)
			if err == nil {
				t.Fatal("expected error for unrepresentable or oversized length")
			}
		})
	}
}
