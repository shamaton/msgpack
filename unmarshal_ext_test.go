package msgpack_test

import (
	"errors"
	"testing"

	"github.com/shamaton/msgpack/v2"
	"github.com/shamaton/msgpack/v2/def"
)

func TestUnmarshalTruncatedTimestampExtReturnsTooShort(t *testing.T) {
	ts := def.TimeStamp

	testcases := []struct {
		name string
		data []byte
	}{
		{name: "fixext4", data: []byte{def.Fixext4, byte(ts)}},
		{name: "fixext8", data: []byte{def.Fixext8, byte(ts)}},
	}

	methods := []struct {
		name string
		fn   func([]byte, interface{}) error
	}{
		{name: "Unmarshal", fn: msgpack.Unmarshal},
		{name: "UnmarshalAsArray", fn: msgpack.UnmarshalAsArray},
		{name: "UnmarshalAsMap", fn: msgpack.UnmarshalAsMap},
	}

	for _, method := range methods {
		for _, tc := range testcases {
			t.Run(method.name+"/"+tc.name, func(t *testing.T) {
				var v interface{}
				err := method.fn(tc.data, &v)
				if !errors.Is(err, def.ErrTooShortBytes) {
					t.Fatalf("expected %v, got %v", def.ErrTooShortBytes, err)
				}
			})
		}
	}
}
