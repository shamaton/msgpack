package decoding

import (
	"testing"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
	"github.com/shamaton/msgpack/v2/time"
)

func Test_AddExtDecoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		AddExtDecoder(time.Decoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

func Test_RemoveExtDecoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		RemoveExtDecoder(time.Decoder)
		tu.Equal(t, len(extCoders), 1)
	})
}
