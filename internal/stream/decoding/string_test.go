package decoding

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/internal/common"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func Test_stringByteLength(t *testing.T) {
	testcases := []struct {
		name     string
		code     byte
		length   int
		expected int
		errSkip  bool
	}{
		{
			name:     "FixStr",
			code:     def.FixStr + 1,
			expected: 1,
			errSkip:  true,
		},
		{
			name:     "Str8",
			code:     def.Str8,
			length:   1,
			expected: math.MaxUint8,
		},
		{
			name:     "Str16",
			code:     def.Str16,
			length:   2,
			expected: math.MaxUint16,
		},
		{
			name:     "Str32",
			code:     def.Str32,
			length:   4,
			expected: math.MaxUint32,
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
				_, err := d.stringByteLength(tc.code, reflect.String)
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
				v, err := d.stringByteLength(tc.code, reflect.String)
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

func Test_asString(t *testing.T) {
	t.Run("read error", func(t *testing.T) {
		d := decoder{
			r:   tu.NewErrReader(),
			buf: common.GetBuffer(),
		}
		v, err := d.asString(reflect.String)
		tu.IsError(t, err, tu.ErrReaderErr)
		tu.Equal(t, v, emptyString)
	})
	t.Run("ok", func(t *testing.T) {
		d := decoder{
			r:   bytes.NewReader([]byte{def.FixStr + 1, 'a'}),
			buf: common.GetBuffer(),
		}
		v, err := d.asString(reflect.String)
		tu.NoError(t, err)
		tu.Equal(t, v, "a")
	})
}

func Test_asStringByte(t *testing.T) {
	t.Run("read error", func(t *testing.T) {
		d := decoder{
			r:   tu.NewErrReader(),
			buf: common.GetBuffer(),
		}
		v, err := d.asStringByte(reflect.String)
		tu.IsError(t, err, tu.ErrReaderErr)
		tu.EqualSlice(t, v, emptyBytes)
	})
	t.Run("ok", func(t *testing.T) {
		d := decoder{
			r:   bytes.NewReader([]byte{def.FixStr + 1, 'a'}),
			buf: common.GetBuffer(),
		}
		v, err := d.asStringByte(reflect.String)
		tu.NoError(t, err)
		tu.EqualSlice(t, v, []byte("a"))
	})
}
