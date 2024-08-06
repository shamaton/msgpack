package decoding

import (
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_isCodeBin(t *testing.T) {
	d := decoder{}
	for i := 0x00; i <= 0xff; i++ {
		v := byte(i)
		isBin := v == def.Bin8 || v == def.Bin16 || v == def.Bin32
		tu.Equal(t, d.isCodeBin(v), isBin)
	}
}

func Test_asBinWithCode(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) ([]byte, error) {
		return d.asBinWithCode
	}
	testcases := AsXXXTestCases[[]byte]{
		{
			Name:             "Bin8.error.size",
			Code:             def.Bin8,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin8.error.data",
			Code:             def.Bin8,
			Data:             []byte{1},
			ReadCount:        1,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin8.ok",
			Code:             def.Bin8,
			Data:             []byte{1, 'a'},
			Expected:         []byte{'a'},
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin16.error",
			Code:             def.Bin16,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin16.error.data",
			Code:             def.Bin16,
			Data:             []byte{0, 1},
			ReadCount:        1,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin16.ok",
			Code:             def.Bin16,
			Data:             []byte{0, 1, 'b'},
			Expected:         []byte{'b'},
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin32.error",
			Code:             def.Bin32,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin32.error.data",
			Code:             def.Bin32,
			Data:             []byte{0, 0, 0, 1},
			ReadCount:        1,
			Error:            io.EOF,
			MethodAsWithCode: method,
		},
		{
			Name:             "Bin32.ok",
			Code:             def.Bin32,
			Data:             []byte{0, 0, 0, 1, 'c'},
			Expected:         []byte{'c'},
			ReadCount:        2,
			MethodAsWithCode: method,
		},
		{
			Name:             "Unexpected",
			Code:             def.Nil,
			IsTemplateError:  true,
			MethodAsWithCode: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}
