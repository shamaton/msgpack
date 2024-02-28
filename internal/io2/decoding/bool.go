package decoding

import (
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func asBool(r io.Reader, k reflect.Kind) (bool, error) {
	code, err := readSize1(r)
	if err != nil {
		return false, err
	}
	return asBoolWithCode(r, code, k)
}

func asBoolWithCode(_ io.Reader, code byte, k reflect.Kind) (bool, error) {
	switch code {
	case def.True:
		return true, nil
	case def.False:
		return false, nil
	}
	return false, errorTemplate(code, k)
}
