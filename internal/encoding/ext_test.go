package encoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/ext"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
	"github.com/shamaton/msgpack/v3/time"
)

func Test_AddExtEncoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		AddExtEncoder(time.Encoder)
		coders := currentExtEncoderRegistry.Load().byKind[reflect.Struct]
		tu.Equal(t, len(coders), 1)
		tu.Equal(t, coders[time.Encoder.Type()], ext.Encoder(time.Encoder))
	})
}

func Test_RemoveExtEncoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		RemoveExtEncoder(time.Encoder)
		coders := currentExtEncoderRegistry.Load().byKind[reflect.Struct]
		tu.Equal(t, len(coders), 1)
		tu.Equal(t, coders[time.Encoder.Type()], ext.Encoder(time.Encoder))
	})
}
