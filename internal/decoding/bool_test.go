package decoding

import (
	"github.com/shamaton/msgpack/v2/def"
	"reflect"
	"testing"
)

func Test_asBool(t *testing.T) {
	method := func(d *decoder) func(int, reflect.Kind) (bool, int, error) {
		return d.asBool
	}
	testcases := AsXXXTestCases[bool]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bool.false",
			Data:     []byte{def.False},
			Expected: false,
			MethodAs: method,
		},
		{
			Name:     "Bool.true",
			Data:     []byte{def.True},
			Expected: true,
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
