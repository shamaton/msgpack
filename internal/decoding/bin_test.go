package decoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func Test_isCodeBin(t *testing.T) {
	d := decoder{}
	for i := 0x00; i <= 0xff; i++ {
		v := byte(i)
		isBin := v == def.Bin8 || v == def.Bin16 || v == def.Bin32
		tu.Equal(t, d.isCodeBin(v), isBin)
	}
}

func Test_asBin(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) ([]byte, int, error) {
		return d.asBin
	}
	testcases := AsXXXTestCases[[]byte]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin8.error.size",
			Data:     []byte{def.Bin8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin8.error.data",
			Data:     []byte{def.Bin8, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin8.ok",
			Data:     []byte{def.Bin8, 1, 'a'},
			Expected: []byte{'a'},
			MethodAs: method,
		},
		{
			Name:     "Bin16.error.size",
			Data:     []byte{def.Bin16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin16.error.data",
			Data:     []byte{def.Bin16, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin16.ok",
			Data:     []byte{def.Bin16, 0, 1, 'b'},
			Expected: []byte{'b'},
			MethodAs: method,
		},
		{
			Name:     "Bin32.error.size",
			Data:     []byte{def.Bin32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin32.error.data",
			Data:     []byte{def.Bin32, 0, 0, 0, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin32.ok",
			Data:     []byte{def.Bin32, 0, 0, 0, 1, 'c'},
			Expected: []byte{'c'},
			MethodAs: method,
		},
		{
			Name:     "Unexpected",
			Data:     []byte{def.Nil},
			Error:    def.ErrCanNotDecode,
			MethodAs: method,
		},
	}
	testcases.Run(t)
}
