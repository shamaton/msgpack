package encoding

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v2/internal/common"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

const dummyByte = 0xc1

type TestWriter struct {
	WrittenBytes []byte
}

var _ io.Writer = (*TestWriter)(nil)

var ErrTestWriter = errors.New("expected written error")

func (w *TestWriter) Write(p []byte) (n int, err error) {
	w.WrittenBytes = append(w.WrittenBytes, p...)
	if bytes.Contains(p, []byte{dummyByte}) {
		return 0, ErrTestWriter
	}
	return len(p), nil
}

func NewTestWriter() *TestWriter {
	return &TestWriter{}
}

type AsXXXTestCase[T any] struct {
	Name            string
	Value           T
	Expected        []byte
	Contains        []byte
	BufferSize      int
	PreWriteSize    int
	Method          func(*encoder) func(T) error
	MethodForFixed  func(*encoder) func(reflect.Value) (bool, error)
	MethodForStruct func(*encoder) func(reflect.Value) error
}

type AsXXXTestCases[T any] []AsXXXTestCase[T]

func (tcs AsXXXTestCases[T]) Run(t *testing.T) {
	for _, tc := range tcs {
		tc.Run(t)
	}
}

func (tc *AsXXXTestCase[T]) Run(t *testing.T) {
	t.Helper()

	if tc.Method == nil && tc.MethodForFixed == nil && tc.MethodForStruct == nil {
		t.Fatal("must set either Method or MethodForFixed or MethodForStruct")
	}

	method := func(e *encoder) error {
		if tc.Method != nil {
			return tc.Method(e)(tc.Value)
		}
		if tc.MethodForFixed != nil {
			_, err := tc.MethodForFixed(e)(reflect.ValueOf(tc.Value))
			return err
		}
		if tc.MethodForStruct != nil {
			return tc.MethodForStruct(e)(reflect.ValueOf(tc.Value))
		}
		panic("unreachable")
	}

	t.Run(tc.Name, func(t *testing.T) {
		w := NewTestWriter()
		e := encoder{
			w:      w,
			buf:    common.GetBuffer(),
			Common: common.Common{},
		}

		if tc.BufferSize < tc.PreWriteSize {
			t.Fatal("buffer size must be greater than pre write size")
		}

		e.buf.Data = make([]byte, tc.BufferSize)
		if tc.PreWriteSize > 0 {
			for i := 0; i < tc.PreWriteSize; i++ {
				_ = e.buf.Write(e.w, dummyByte)
			}
		}

		err := method(&e)
		_ = e.buf.Flush(w)
		common.PutBuffer(e.buf)

		if tc.PreWriteSize > 0 {
			tu.IsError(t, err, ErrTestWriter)
			if !bytes.Contains(w.WrittenBytes, tc.Contains) {
				t.Fatalf("[% 02x] does not contain in [% 02x]", tc.Contains, w.WrittenBytes)
			}
			return
		}

		tu.NoError(t, err)
		tu.EqualSlice(t, w.WrittenBytes, tc.Expected)
	})
}
