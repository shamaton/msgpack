package decoding

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/internal/common"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func Test_asUint(t *testing.T) {
	t.Run("read error", func(t *testing.T) {
		d := decoder{
			r:   tu.NewErrReader(),
			buf: common.GetBuffer(),
		}
		v, err := d.asUint(reflect.Uint)
		tu.IsError(t, err, tu.ErrReaderErr)
		tu.Equal(t, v, 0)
	})
	t.Run("ok", func(t *testing.T) {
		d := decoder{
			r:   bytes.NewReader([]byte{1}),
			buf: common.GetBuffer(),
		}
		v, err := d.asUint(reflect.Uint)
		tu.NoError(t, err)
		tu.Equal(t, v, 1)
	})
}

func Test_asUintWithCode(t *testing.T) {
	testcases := []struct {
		name     string
		code     byte
		length   int
		expected uint64
		errSkip  bool
	}{
		{
			name:     "Uint8",
			code:     def.Uint8,
			length:   1,
			expected: math.MaxUint8,
		},
		{
			name:     "Int8",
			code:     def.Int8,
			length:   1,
			expected: math.MaxUint64,
		},
		{
			name:     "Uint16",
			code:     def.Uint16,
			length:   2,
			expected: math.MaxUint16,
		},
		{
			name:     "Int16",
			code:     def.Int16,
			length:   2,
			expected: math.MaxUint64,
		},
		{
			name:     "Uint32",
			code:     def.Uint32,
			length:   4,
			expected: math.MaxUint32,
		},
		{
			name:     "Int32",
			code:     def.Int32,
			length:   4,
			expected: math.MaxUint64,
		},
		{
			name:     "Uint64",
			code:     def.Uint64,
			length:   8,
			expected: math.MaxUint64,
		},
		{
			name:     "Int64",
			code:     def.Int64,
			length:   8,
			expected: math.MaxUint64,
		},
		{
			name:     "Nil",
			code:     def.Nil,
			expected: 0,
			errSkip:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name+"", func(t *testing.T) {
			t.Run("ng", func(t *testing.T) {
				if tc.errSkip {
					t.Log("this testcase is skipped by skip flag")
					return
				}
				d := decoder{
					r:   tu.NewErrReader(),
					buf: common.GetBuffer(),
				}
				defer common.PutBuffer(d.buf)
				_, err := d.asUintWithCode(tc.code, reflect.String)
				tu.IsError(t, err, tu.ErrReaderErr)
			})
			t.Run("ok", func(t *testing.T) {
				data := make([]byte, tc.length)
				for i := range data {
					data[i] = 0xff
				}

				d := decoder{
					r:   bytes.NewReader(data),
					buf: common.GetBuffer(),
				}
				defer common.PutBuffer(d.buf)
				v, err := d.asUintWithCode(tc.code, reflect.String)
				tu.NoError(t, err)
				tu.Equal(t, v, tc.expected)

				p := make([]byte, 1)
				n, err := d.r.Read(p)
				tu.IsError(t, err, io.EOF)
				tu.Equal(t, n, 0)
			})
		})
	}
}
