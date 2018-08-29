package deserialize

import (
	"reflect"

	"github.com/shamaton/msgpack/def"
)

func (d *deserializer) asBool(offset int, k reflect.Kind) (bool, int, error) {
	code := d.data[offset]
	offset++

	// todo : use switch
	if code == def.True {
		return true, offset, nil
	} else if code == def.False {
		return false, offset, nil
	}
	return false, 0, d.errorTemplate(code, k)
}
