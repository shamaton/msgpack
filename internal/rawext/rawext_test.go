package rawext

import (
	"math"
	"strconv"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

func Test_headerSize(t *testing.T) {
	testcases := []struct {
		n           int
		wantCode    byte
		wantLenSize uint64
	}{
		{n: 1, wantCode: def.Fixext1, wantLenSize: 0},
		{n: 2, wantCode: def.Fixext2, wantLenSize: 0},
		{n: 4, wantCode: def.Fixext4, wantLenSize: 0},
		{n: 8, wantCode: def.Fixext8, wantLenSize: 0},
		{n: 16, wantCode: def.Fixext16, wantLenSize: 0},
		{n: 0, wantCode: def.Ext8, wantLenSize: 1},
		{n: 3, wantCode: def.Ext8, wantLenSize: 1},
		{n: math.MaxUint8, wantCode: def.Ext8, wantLenSize: 1},
		{n: math.MaxUint8 + 1, wantCode: def.Ext16, wantLenSize: 2},
		{n: math.MaxUint16, wantCode: def.Ext16, wantLenSize: 2},
		{n: math.MaxUint16 + 1, wantCode: def.Ext32, wantLenSize: 4},
	}
	for _, tc := range testcases {
		code, lengthBytes := headerSize(tc.n)
		if code != tc.wantCode || lengthBytes != tc.wantLenSize {
			t.Fatalf("headerSize(%d) = (%x, %d), want (%x, %d)", tc.n, code, lengthBytes, tc.wantCode, tc.wantLenSize)
		}
	}
}

func Test_validateLength(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("requires a 64-bit int to represent lengths beyond ext16's range")
	}
	// math.MaxUint32 cannot be passed directly as an int argument: as an
	// untyped constant it would fail to compile on 32-bit platforms, where
	// int is only 32 bits wide. Routing it through a uint64 variable first
	// makes the conversion a runtime operation instead of a constant one, so
	// the test file itself still compiles everywhere (it is simply skipped
	// above on 32-bit platforms).
	var maxU64 uint64 = math.MaxUint32
	maxLen := int(maxU64)

	if err := validateLength(maxLen); err != nil {
		t.Fatalf("unexpected error for max representable length: %v", err)
	}
	err := validateLength(maxLen + 1)
	if err == nil {
		t.Fatal("expected error for length exceeding ext32's 4-byte length prefix, got nil")
	}
	if err != def.ErrUnsupportedLength {
		t.Fatalf("expected %v, got %v", def.ErrUnsupportedLength, err)
	}
}

// Test_validateLength_negative exercises the n < 0 guard directly. n can
// never actually be negative via the normal len(r.Data) call path, but the
// guard exists to make the uint64(n) conversion provably non-negative for
// static analysis (gosec G115); this test ensures that branch is reachable
// rather than being untested dead code.
func Test_validateLength_negative(t *testing.T) {
	err := validateLength(-1)
	if err != def.ErrUnsupportedLength {
		t.Fatalf("expected %v, got %v", def.ErrUnsupportedLength, err)
	}
}

func Test_byteSize_emptyPayload(t *testing.T) {
	r := ext.Raw{Type: 1, Data: make([]byte, 0)}
	if _, err := byteSize(r); err != nil {
		t.Fatalf("unexpected error for empty payload: %v", err)
	}
}

// Test_totalSize_boundsAndOverflow exercises totalSize right at the largest
// n it can accept, and one past it, without allocating a real []byte of
// that length. Which of totalSize's two guards actually produces the error
// at that boundary is platform-dependent:
//   - On 64-bit platforms, math.MaxInt is far larger than maxDataLen (ext32's
//     ~4GiB length-prefix limit), so validateLength's maxDataLen check is the
//     one that rejects n+1.
//   - On 32-bit platforms, math.MaxInt (int is 32 bits wide there) is
//     smaller than maxDataLen, so the header-byte arithmetic would overflow
//     int first; totalSize's explicit `total > uint64(math.MaxInt)` check is
//     the one that rejects n+1 instead.
//
// Computing the boundary as min(maxDataLen, math.MaxInt-headerOverhead)
// means this test exercises whichever guard is actually reachable on the
// platform it runs on, instead of only ever covering the 64-bit case.
func Test_totalSize_boundsAndOverflow(t *testing.T) {
	// Ext32's code byte + 4 length bytes + the ext type byte: the fixed
	// overhead totalSize adds on top of n for any n large enough to need
	// the ext32 form.
	const headerOverhead = uint64(def.Byte1 + 4 + def.Byte1)

	maxValidN := maxDataLen
	if uint64(math.MaxInt)-headerOverhead < maxDataLen {
		maxValidN = uint64(math.MaxInt) - headerOverhead
	}
	n := int(maxValidN)

	if _, err := totalSize(n); err != nil {
		t.Fatalf("unexpected error at boundary n=%d: %v", n, err)
	}
	if _, err := totalSize(n + 1); err != def.ErrUnsupportedLength {
		t.Fatalf("expected %v for n=%d, got %v", def.ErrUnsupportedLength, n+1, err)
	}
}

// Test_totalSize_negative exercises totalSize's own n < 0 guard directly. n
// can never actually be negative via the normal len(r.Data) call path, but
// the guard exists to make the uint64(n) conversion provably non-negative
// for static analysis (gosec G115); this test ensures that branch is
// reachable rather than being untested dead code.
func Test_totalSize_negative(t *testing.T) {
	if _, err := totalSize(-1); err != def.ErrUnsupportedLength {
		t.Fatalf("expected %v, got %v", def.ErrUnsupportedLength, err)
	}
}
