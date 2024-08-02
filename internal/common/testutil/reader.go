package testutil

import "errors"

var ErrReaderErr = errors.New("reader error")

type ErrReader struct{}

func NewErrReader() *ErrReader {
	return &ErrReader{}
}

func (ErrReader) Read(_ []byte) (int, error) {
	return 0, ErrReaderErr
}
