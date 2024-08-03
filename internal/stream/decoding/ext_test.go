package decoding

import (
	"io"
	"testing"

	"github.com/shamaton/msgpack/v2/def"

	"github.com/shamaton/msgpack/v2/internal/common"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
	"github.com/shamaton/msgpack/v2/time"
)

func Test_AddExtDecoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		AddExtDecoder(time.StreamDecoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

func Test_RemoveExtDecoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		RemoveExtDecoder(time.StreamDecoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

func Test_readIfExtType(t *testing.T) {
	testcases := []struct {
		name  string
		code  byte
		r     []byte
		typ   int8
		data  []byte
		err   error
		count int
	}{
		{
			name:  "Fixext1.error.type",
			code:  def.Fixext1,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Fixext1.error.data",
			code:  def.Fixext1,
			r:     []byte{1},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Fixext1.ok",
			code:  def.Fixext1,
			r:     []byte{1, 2},
			typ:   1,
			data:  []byte{2},
			count: 2,
		},

		{
			name:  "Fixext2.error.type",
			code:  def.Fixext2,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Fixext2.error.data",
			code:  def.Fixext2,
			r:     []byte{2},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Fixext2.ok",
			code:  def.Fixext2,
			r:     []byte{2, 3, 4},
			typ:   2,
			data:  []byte{3, 4},
			count: 2,
		},

		{
			name:  "Fixext4.error.type",
			code:  def.Fixext4,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Fixext4.error.data",
			code:  def.Fixext4,
			r:     []byte{4},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Fixext4.ok",
			code:  def.Fixext4,
			r:     []byte{4, 5, 6, 7, 8},
			typ:   4,
			data:  []byte{5, 6, 7, 8},
			count: 2,
		},

		{
			name:  "Fixext8.error.type",
			code:  def.Fixext8,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Fixext8.error.data",
			code:  def.Fixext8,
			r:     []byte{8},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Fixext8.ok",
			code:  def.Fixext8,
			r:     []byte{8, 1, 2, 3, 4, 5, 6, 7, 8},
			typ:   8,
			data:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
			count: 2,
		},

		{
			name:  "Fixext16.error.type",
			code:  def.Fixext16,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Fixext16.error.data",
			code:  def.Fixext16,
			r:     []byte{16},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Fixext16.ok",
			code:  def.Fixext16,
			r:     []byte{16, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			typ:   16,
			data:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			count: 2,
		},

		{
			name:  "Ext8.error.size",
			code:  def.Ext8,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Ext8.error.type",
			code:  def.Ext8,
			r:     []byte{1},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Ext8.error.data",
			code:  def.Ext8,
			r:     []byte{1, 18},
			err:   io.EOF,
			count: 2,
		},
		{
			name:  "Ext8.ok",
			code:  def.Ext8,
			r:     []byte{1, 18, 2},
			typ:   18,
			data:  []byte{2},
			count: 3,
		},

		{
			name:  "Ext16.error.size",
			code:  def.Ext16,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Ext16.error.type",
			code:  def.Ext16,
			r:     []byte{0, 1},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Ext16.error.data",
			code:  def.Ext16,
			r:     []byte{0, 1, 24},
			err:   io.EOF,
			count: 2,
		},
		{
			name:  "Ext16.ok",
			code:  def.Ext16,
			r:     []byte{0, 1, 24, 3},
			typ:   24,
			data:  []byte{3},
			count: 3,
		},

		{
			name:  "Ext32.error.size",
			code:  def.Ext32,
			r:     []byte{},
			err:   io.EOF,
			count: 0,
		},
		{
			name:  "Ext32.error.type",
			code:  def.Ext32,
			r:     []byte{0, 0, 0, 1},
			err:   io.EOF,
			count: 1,
		},
		{
			name:  "Ext32.error.data",
			code:  def.Ext32,
			r:     []byte{0, 0, 0, 1, 32},
			err:   io.EOF,
			count: 2,
		},
		{
			name:  "Ext32.ok",
			code:  def.Ext32,
			r:     []byte{0, 0, 0, 1, 32, 4},
			typ:   32,
			data:  []byte{4},
			count: 3,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := tu.NewTestReader(tc.r)
			d := decoder{
				r:   r,
				buf: common.GetBuffer(),
			}
			defer common.PutBuffer(d.buf)
			typ, data, err := d.readIfExtType(tc.code)
			tu.IsError(t, err, tc.err)
			tu.Equal(t, tc.typ, typ)
			tu.EqualSlice(t, data, tc.data)
			tu.Equal(t, r.Count(), tc.count)
		})
	}
}
