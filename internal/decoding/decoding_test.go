package decoding

import (
	"reflect"
	"testing"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

type AsXXXTestCase[T any] struct {
	Name             string
	Code             byte
	Data             []byte
	Expected         T
	Error            error
	MethodAs         func(d *decoder) func(int, reflect.Kind) (T, int, error)
	MethodAsWithCode func(d *decoder) func(byte, reflect.Kind) (T, int, error)
	MethodAsCustom   func(d *decoder) (T, int, error)
}

type AsXXXTestCases[T any] []AsXXXTestCase[T]

func (tcs AsXXXTestCases[T]) Run(t *testing.T) {
	for _, tc := range tcs {
		tc.Run(t)
	}
}

func (tc *AsXXXTestCase[T]) Run(t *testing.T) {
	const kind = reflect.String
	t.Helper()

	if tc.MethodAs == nil && tc.MethodAsWithCode == nil && tc.MethodAsCustom == nil {
		t.Fatal("must set either method or methodAsWithCode or MethodAsCustom")
	}

	methodAs := func(d *decoder) (T, int, error) {
		if tc.MethodAs != nil {
			return tc.MethodAs(d)(0, kind)
		}
		if tc.MethodAsWithCode != nil {
			return tc.MethodAsWithCode(d)(tc.Code, kind)
		}
		if tc.MethodAsCustom != nil {
			return tc.MethodAsCustom(d)
		}
		panic("unreachable")
	}

	t.Run(tc.Name, func(t *testing.T) {
		d := decoder{
			data: tc.Data,
		}

		v, offset, err := methodAs(&d)
		if tc.Error != nil {
			tu.IsError(t, err, tc.Error)
			return
		}
		tu.NoError(t, err)
		tu.Equal(t, v, tc.Expected)
		tu.Equal(t, offset, len(tc.Data))
	})
}
