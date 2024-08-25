package common

import (
	"reflect"
	"testing"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func TestCommon_CheckField(t *testing.T) {
	common := Common{}

	t.Run("tag:-", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:"-"`,
		}
		b, v := common.CheckField(field)
		tu.Equal(t, b, false)
		tu.Equal(t, v, "")
	})
	t.Run("tag:B", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:"B"`,
		}
		b, v := common.CheckField(field)
		tu.Equal(t, b, true)
		tu.Equal(t, v, "B")
	})
	t.Run("name:A", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:""`,
		}
		b, v := common.CheckField(field)
		tu.Equal(t, b, true)
		tu.Equal(t, v, "A")
	})
	t.Run("private", func(t *testing.T) {
		field := reflect.StructField{
			Name: "a",
		}
		b, v := common.CheckField(field)
		tu.Equal(t, b, false)
		tu.Equal(t, v, "")
	})
}
