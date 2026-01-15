package encoding

import (
	"testing"

	"github.com/shamaton/msgpack/v3/def"
)

func Test_asBool(t *testing.T) {
	method := func(e *encoder) func(bool) error {
		return e.writeBool
	}
	testcases := AsXXXTestCases[bool]{
		{
			Name:         "True.error",
			Value:        true,
			BufferSize:   1,
			PreWriteSize: 1,
			Method:       method,
		},
		{
			Name:       "True.ok",
			Value:      true,
			Expected:   []byte{def.True},
			BufferSize: 1,
			Method:     method,
		},
	}
	testcases.Run(t)
}
