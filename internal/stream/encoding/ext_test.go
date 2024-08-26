package encoding

import (
	"testing"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
	"github.com/shamaton/msgpack/v2/time"
)

func Test_AddExtEncoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		AddExtEncoder(time.StreamEncoder)
		tu.Equal(t, len(extCoders), 1)
	})
}

func Test_RemoveExtEncoder(t *testing.T) {
	t.Run("ignore", func(t *testing.T) {
		RemoveExtEncoder(time.StreamEncoder)
		tu.Equal(t, len(extCoders), 1)
	})
}
