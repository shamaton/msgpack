package encoding

import (
	"math"
	"testing"

	"github.com/shamaton/msgpack/v2/def"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_calcByteSlice(t *testing.T) {
	testcases := []struct {
		name   string
		value  int
		result int
		error  error
	}{
		{
			name:   "u8",
			value:  math.MaxUint8,
			result: def.Byte1 + math.MaxUint8,
		},
		{
			name:   "u16",
			value:  math.MaxUint16,
			result: def.Byte2 + math.MaxUint16,
		},
		{
			name:   "u32",
			value:  math.MaxUint32,
			result: def.Byte4 + math.MaxUint32,
		},
		{
			name:  "u32over",
			value: math.MaxUint32 + 1,
			error: def.ErrUnsupportedType,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := encoder{}
			result, err := e.calcByteSlice(tc.value)
			tu.IsError(t, err, tc.error)
			tu.Equal(t, result, tc.result)
		})
	}
}
