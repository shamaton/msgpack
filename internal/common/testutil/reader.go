package testutil

import (
	"errors"
	"io"
)

var ErrReaderErr = errors.New("reader error")

type ErrReader struct{}

func NewErrReader() *ErrReader {
	return &ErrReader{}
}

func (ErrReader) Read(_ []byte) (int, error) {
	return 0, ErrReaderErr
}

type TestReader struct {
	s     []byte
	i     int64 // current reading index
	count int
}

func NewTestReader(b []byte) *TestReader {
	return &TestReader{s: b, i: 0, count: 0}
}

func (r *TestReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	r.count++
	return
}

func (r *TestReader) Count() int {
	return r.count
}
