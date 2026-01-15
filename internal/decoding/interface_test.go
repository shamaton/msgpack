package decoding

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

func Test_asInterfaceWithCode(t *testing.T) {
	dec := testExt2Decoder{}
	AddExtDecoder(&dec)
	defer RemoveExtDecoder(&dec)

	method := func(d *decoder) func(int, reflect.Kind) (any, int, error) {
		return d.asInterface
	}
	testcases := AsXXXTestCases[any]{
		{
			Name:     "error.code",
			Data:     []byte{},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint8.error",
			Data:     []byte{def.Uint8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint16.error",
			Data:     []byte{def.Uint16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint32.error",
			Data:     []byte{def.Uint32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Uint64.error",
			Data:     []byte{def.Uint64},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int8.error",
			Data:     []byte{def.Int8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int16.error",
			Data:     []byte{def.Int16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int32.error",
			Data:     []byte{def.Int32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Int64.error",
			Data:     []byte{def.Int64},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Float32.error",
			Data:     []byte{def.Float32},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Float64.error",
			Data:     []byte{def.Float64},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Str.error",
			Data:     []byte{def.Str8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Bin.error",
			Data:     []byte{def.Bin8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array.error.length",
			Data:     []byte{def.Array16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Array.error.required",
			Data:     []byte{def.Array16, 0, 1},
			Error:    def.ErrLackDataLengthToSlice,
			MethodAs: method,
		},
		{
			Name:     "Array.error.set",
			Data:     []byte{def.Array16, 0, 1, def.Int8},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map.error.length",
			Data:     []byte{def.Map16},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map.error.required",
			Data:     []byte{def.Map16, 0, 1},
			Error:    def.ErrLackDataLengthToMap,
			MethodAs: method,
		},
		{
			Name:     "Map.error.set.can.slice",
			Data:     []byte{def.Map16, 0, 1, def.Array16, 0},
			Error:    def.ErrCanNotSetSliceAsMapKey,
			MethodAs: method,
		},
		{
			Name:     "Map.error.set.can.map",
			Data:     []byte{def.Map16, 0, 1, def.Map16, 0},
			Error:    def.ErrCanNotSetMapAsMapKey,
			MethodAs: method,
		},
		{
			Name:     "Map.error.set.key",
			Data:     []byte{def.Map16, 0, 1, def.Str8, 1},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "Map.error.set.value",
			Data:     []byte{def.Map16, 0, 1, def.FixStr + 1, 'a'},
			Error:    def.ErrTooShortBytes,
			MethodAs: method,
		},
		{
			Name:     "ExtCoder.error",
			Data:     []byte{def.Fixext1, 3},
			Error:    ErrTestExtDecoder,
			MethodAs: method,
		},
	}

	for _, tc := range testcases {
		tc.Run(t)
	}
}

// TODO: to testutil
type testExt2Decoder struct {
	ext.DecoderCommon
}

var _ ext.Decoder = (*testExt2Decoder)(nil)

func (td *testExt2Decoder) Code() int8 {
	return 3
}

func (td *testExt2Decoder) IsType(o int, d *[]byte) bool {
	// todo : lack of error handling
	code, _ := td.ReadSize1(o, d)
	if code == def.Fixext1 {
		extCode, _ := td.ReadSize1(o+1, d)
		return int8(extCode) == td.Code()
	}
	return false
}

var ErrTestExtDecoder = fmt.Errorf("testExtDecoder")

func (td *testExt2Decoder) AsValue(_ int, _ reflect.Kind, _ *[]byte) (any, int, error) {
	return nil, 0, ErrTestExtDecoder
}
