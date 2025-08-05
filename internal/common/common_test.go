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
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, false)
		tu.Equal(t, omit, false)
		tu.Equal(t, v, "")
	})
	t.Run("tag:,omitempty", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:",omitempty"`,
		}
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, true)
		tu.Equal(t, omit, true)
		tu.Equal(t, v, "A")
	})
	t.Run("tag:-,omitempty", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:"-,omitempty"`,
		}
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, false)
		tu.Equal(t, omit, false)
		tu.Equal(t, v, "")
	})
	t.Run("tag:B", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:"B"`,
		}
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, true)
		tu.Equal(t, omit, false)
		tu.Equal(t, v, "B")
	})
	t.Run("tag:B,omitempty", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:"B,omitempty"`,
		}
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, true)
		tu.Equal(t, omit, true)
		tu.Equal(t, v, "B")
	})
	t.Run("name:A", func(t *testing.T) {
		field := reflect.StructField{
			Name: "A",
			Tag:  `msgpack:""`,
		}
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, true)
		tu.Equal(t, omit, false)
		tu.Equal(t, v, "A")
	})
	t.Run("private", func(t *testing.T) {
		field := reflect.StructField{
			Name: "a",
		}
		public, omit, v := common.CheckField(field)
		tu.Equal(t, public, false)
		tu.Equal(t, omit, false)
		tu.Equal(t, v, "")
	})
}
