package decoding

import (
	"math"
	"testing"
)

func Test_lengthFromUint32(t *testing.T) {
	tests := []struct {
		u       uint32
		wantErr bool
	}{
		{0, false},
		{1, false},
		{math.MaxInt32, false},
		{uint32(math.MaxInt32) + 1, math.MaxInt == math.MaxInt32},  // error only on 32-bit
		{0xffffffff, math.MaxInt == math.MaxInt32},                   // error only on 32-bit
	}
	for _, tc := range tests {
		got, err := lengthFromUint32(tc.u)
		if (err != nil) != tc.wantErr {
			t.Errorf("lengthFromUint32(%d): err=%v, wantErr=%v", tc.u, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != int(tc.u) {
			t.Errorf("lengthFromUint32(%d) = %d, want %d", tc.u, got, tc.u)
		}
	}
}
