package decoding

import (
	"bufio"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

func (d *decoder) asBool(reader *bufio.Reader, k reflect.Kind) (bool, error) {
	code, err := reader.ReadByte()
	if err != nil {
		return false, err
	}

	switch code {
	case def.True:
		return true, nil
	case def.False:
		return false, nil
	}

	return false, d.errorTemplate(code, k)
}
