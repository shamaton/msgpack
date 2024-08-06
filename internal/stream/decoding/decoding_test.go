package decoding

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/internal/common"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

type AsXXXTestCase[T any] struct {
	Name             string
	Code             byte
	Data             []byte
	ReadCount        int
	Expected         T
	Error            error
	IsTemplateError  bool
	MethodAs         func(d *decoder) func(reflect.Kind) (T, error)
	MethodAsWithCode func(d *decoder) func(byte, reflect.Kind) (T, error)
	MethodAsCustom   func(d *decoder) (T, error)
}

type AsXXXTestCases[T any] []AsXXXTestCase[T]

func (tc *AsXXXTestCase[T]) Run(t *testing.T) {
	const kind = reflect.String
	t.Helper()

	if tc.MethodAs == nil && tc.MethodAsWithCode == nil && tc.MethodAsCustom == nil {
		t.Fatal("must set either method or methodAsWithCode or MethodAsCustom")
	}

	methodAs := func(d *decoder) (T, error) {
		if tc.MethodAs != nil {
			return tc.MethodAs(d)(kind)
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
		r := tu.NewTestReader(tc.Data)
		d := decoder{
			r:   r,
			buf: common.GetBuffer(),
		}
		defer common.PutBuffer(d.buf)

		v, err := methodAs(&d)
		tu.Equal(t, r.Count(), tc.ReadCount)

		if tc.Error != nil {
			tu.IsError(t, err, tc.Error)
			return
		}
		if tc.IsTemplateError {
			tu.ErrorContains(t, err, fmt.Sprintf("msgpack : invalid code %x", tc.Code))
			return
		}
		tu.NoError(t, err)
		tu.Equal(t, v, tc.Expected)

		p := make([]byte, 1)
		n, err := d.r.Read(p)
		tu.IsError(t, err, io.EOF)
		tu.Equal(t, n, 0)
	})
}
