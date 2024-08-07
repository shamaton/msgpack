package decoding

import (
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
)

func Test_asBool(t *testing.T) {
	method := func(d *decoder) func(reflect.Kind) (bool, error) {
		return d.asBool
	}
	testcases := AsXXXTestCases[bool]{
		{
			Name:      "error",
			ReadCount: 0,
			Error:     io.EOF,
			MethodAs:  method,
		},
		{
			Name:      "ok",
			Data:      []byte{def.True},
			Expected:  true,
			ReadCount: 1,
			MethodAs:  method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

func Test_asBoolWithCode(t *testing.T) {
	method := func(d *decoder) func(byte, reflect.Kind) (bool, error) {
		return d.asBoolWithCode
	}
	testcases := AsXXXTestCases[bool]{
		{
			Name:             "True",
			Code:             def.True,
			Expected:         true,
			MethodAsWithCode: method,
		},
		{
			Name:             "False",
			Code:             def.False,
			Expected:         false,
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
