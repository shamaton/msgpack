// Package rawext provides the built-in encoder support for ext.Raw, the
// fallback value used when decoding an ext payload that no registered
// Decoder claims. It mirrors internal/time's structure: a non-stream and a
// stream encoder, both writing the ext header (choosing the fixext size or
// ext8/16/32 form from the payload length) followed by the type byte and the
// raw payload bytes.
package rawext

import (
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

// TypeOf is the reflect.Type of ext.Raw, used by callers to recognize the
// built-in raw-ext type without importing reflect at every call site.
var TypeOf = reflect.TypeOf(ext.Raw{})

// maxDataLen is the largest payload length representable by the ext32 wire
// format (a 4-byte length prefix). It is compared as uint64 (rather than
// int) because the untyped constant math.MaxUint32 does not fit in a 32-bit
// int, which would otherwise fail to compile on 32-bit platforms.
const maxDataLen uint64 = math.MaxUint32

// validateLength rejects payloads too large to be represented by any ext
// wire format (even ext32's 4-byte length prefix), so encoding fails with a
// clear error instead of silently truncating the length header and writing
// a malformed, oversized ext frame.
func validateLength(n int) error {
	if uint64(n) > maxDataLen {
		return def.ErrUnsupportedLength
	}
	return nil
}

// headerSize returns the ext format byte and the number of bytes used to
// encode the explicit length (0 for the fixext forms, 1/2/4 for ext8/16/32).
// The caller must have already validated n via validateLength.
func headerSize(n int) (code byte, lengthBytes int) {
	switch n {
	case 1:
		return def.Fixext1, 0
	case 2:
		return def.Fixext2, 0
	case 4:
		return def.Fixext4, 0
	case 8:
		return def.Fixext8, 0
	case 16:
		return def.Fixext16, 0
	}
	switch {
	case n <= math.MaxUint8:
		return def.Ext8, 1
	case n <= math.MaxUint16:
		return def.Ext16, 2
	default:
		return def.Ext32, 4
	}
}

// byteSize returns the total number of bytes required to encode r, including
// the ext header, the type byte, and the payload. It returns an error if r.Data
// is too large to be represented by any ext wire format, or if the resulting
// total would overflow int (relevant on 32-bit platforms, where int is only
// 32 bits wide).
func byteSize(r ext.Raw) (int, error) {
	return totalSize(len(r.Data))
}

// totalSize computes the encoded byte size for a payload of length n,
// without requiring an actual []byte of that length to exist. It is the
// core logic behind byteSize, factored out so boundary/overflow cases can
// be exercised in tests without allocating huge slices.
func totalSize(n int) (int, error) {
	if err := validateLength(n); err != nil {
		return 0, err
	}
	_, lengthBytes := headerSize(n)
	total := uint64(def.Byte1) + uint64(lengthBytes) + uint64(def.Byte1) + uint64(n)
	if total > uint64(math.MaxInt) {
		return 0, def.ErrUnsupportedLength
	}
	return int(total), nil
}
