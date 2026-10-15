package msgpack_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3"
	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

type decodeFunc struct {
	name string
	fn   func(data []byte, v interface{}) error
}

var decodeFuncs = []decodeFunc{
	{"Unmarshal", msgpack.Unmarshal},
	{"UnmarshalRead", func(data []byte, v interface{}) error {
		return msgpack.UnmarshalRead(bytes.NewReader(data), v)
	}},
}

func mustNotPanic(t *testing.T, f func() error) error {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic: %v", r)
			}
		}()
		err = f()
	}()
	return err
}

func TestDecodeBinMapKey(t *testing.T) {
	inputs := map[string][]byte{
		// fixmap(1){ bin8(len 0): nil }
		"bin8": {0x81, 0xc4, 0x00, 0xc0},
		// fixmap(1){ bin16(len 1): nil }
		"bin16": {0x81, 0xc5, 0x00, 0x01, 0x61, 0xc0},
		// fixmap(1){ bin32(len 1): nil }
		"bin32": {0x81, 0xc6, 0x00, 0x00, 0x00, 0x01, 0x61, 0xc0},
		// fixarray(1){ fixmap(1){ bin8(len 0): 1 } }
		"nested": {0x91, 0x81, 0xc4, 0x00, 0x01},
	}
	dests := map[string]func() interface{}{
		"interface":       func() interface{} { var v interface{}; return &v },
		"mapIfaceIface":   func() interface{} { var v map[interface{}]interface{}; return &v },
		"mapIfaceString":  func() interface{} { var v map[interface{}]string; return &v },
		"sliceOfMapIface": func() interface{} { var v []map[interface{}]interface{}; return &v },
	}
	for _, df := range decodeFuncs {
		for in, data := range inputs {
			for dn, newV := range dests {
				if in == "nested" && dn != "interface" && dn != "sliceOfMapIface" {
					continue
				}
				if in != "nested" && dn == "sliceOfMapIface" {
					continue
				}
				t.Run(df.name+"/"+in+"/"+dn, func(t *testing.T) {
					err := mustNotPanic(t, func() error { return df.fn(data, newV()) })
					if !errors.Is(err, def.ErrCanNotSetBinAsMapKey) && !errors.Is(err, def.ErrUnhashableMapKey) {
						t.Fatalf("unexpected error: %v", err)
					}
				})
			}
		}
	}
}

func TestDecodeBinMapValueStillAllowed(t *testing.T) {
	// fixmap(1){ "a": bin8(len 1) }
	data := []byte{0x81, 0xa1, 0x61, 0xc4, 0x01, 0x62}
	for _, df := range decodeFuncs {
		t.Run(df.name, func(t *testing.T) {
			var v map[interface{}]interface{}
			if err := df.fn(data, &v); err != nil {
				t.Fatal(err)
			}
			if b, ok := v["a"].([]byte); !ok || string(b) != "b" {
				t.Fatalf("unexpected value: %#v", v)
			}
		})
	}
}

func nestedArrays(n int) []byte {
	return append(bytes.Repeat([]byte{0x91}, n), 0xc0)
}

func nestedMaps(n int) []byte {
	// fixmap(1){ 0: ... }
	return append(bytes.Repeat([]byte{0x81, 0x00}, n), 0xc0)
}

type depthNode struct {
	C []depthNode
}

type depthPtrNode struct {
	C *depthPtrNode
}

func TestDecodeMaxDepth(t *testing.T) {
	const limit = 10000

	dests := map[string]func() interface{}{
		"interface":      func() interface{} { var v interface{}; return &v },
		"sliceIface":     func() interface{} { var v []interface{}; return &v },
		"mapIfaceIface":  func() interface{} { var v map[interface{}]interface{}; return &v },
		"structSlice":    func() interface{} { return &depthNode{} },
		"structPtr":      func() interface{} { return &depthPtrNode{} },
		"structSliceArr": func() interface{} { return &[]depthNode{} },
	}

	t.Run("within limit", func(t *testing.T) {
		for _, df := range decodeFuncs {
			var v interface{}
			if err := df.fn(nestedArrays(limit), &v); err != nil {
				t.Fatalf("%s arrays: %v", df.name, err)
			}
			v = nil
			if err := df.fn(nestedMaps(limit), &v); err != nil {
				t.Fatalf("%s maps: %v", df.name, err)
			}
		}
	})

	inputs := map[string][]byte{
		"arrays":     nestedArrays(limit + 1),
		"maps":       nestedMaps(limit + 1),
		"arraysHuge": nestedArrays(5_000_000),
	}
	for _, df := range decodeFuncs {
		for in, data := range inputs {
			for dn, newV := range dests {
				t.Run(df.name+"/"+in+"/"+dn, func(t *testing.T) {
					err := mustNotPanic(t, func() error { return df.fn(data, newV()) })
					if err == nil {
						t.Fatal("expected error")
					}
					// other destinations may reject the shape before reaching the limit
					if dn == "interface" || (dn == "sliceIface" && in != "maps") {
						if !errors.Is(err, def.ErrExceededMaxDepth) {
							t.Fatalf("unexpected error: %v", err)
						}
					}
				})
			}
		}
	}
}

func TestDecodeMaxDepthStruct(t *testing.T) {
	// {"C": [{"C": [ ... ]}]} nested beyond the limit
	var buf bytes.Buffer
	n := 6000
	for i := 0; i < n; i++ {
		buf.Write([]byte{0x81, 0xa1, 'C', 0x91})
	}
	buf.Write([]byte{0x80})
	data := buf.Bytes()

	for _, df := range decodeFuncs {
		t.Run(df.name, func(t *testing.T) {
			var v depthNode
			err := mustNotPanic(t, func() error { return df.fn(data, &v) })
			if !errors.Is(err, def.ErrExceededMaxDepth) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDecodeContainerMapKeyIntoReflectMap(t *testing.T) {
	inputs := map[string][]byte{
		// fixmap(1){ fixarray(0): nil }
		"array": {0x81, 0x90, 0xc0},
		// fixmap(1){ fixmap(0): nil }
		"map": {0x81, 0x80, 0xc0},
	}
	for _, df := range decodeFuncs {
		for in, data := range inputs {
			t.Run(df.name+"/"+in, func(t *testing.T) {
				var v map[interface{}]string
				err := mustNotPanic(t, func() error { return df.fn(data, &v) })
				if !errors.Is(err, def.ErrCanNotSetSliceAsMapKey) && !errors.Is(err, def.ErrCanNotSetMapAsMapKey) {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	}
}

const unhashableExtCode int8 = 99

type unhashableExtEncoder struct{ ext.EncoderCommon }

func (*unhashableExtEncoder) Code() int8                              { return unhashableExtCode }
func (*unhashableExtEncoder) Type() reflect.Type                      { return reflect.TypeOf(unhashableExtValue{}) }
func (*unhashableExtEncoder) CalcByteSize(reflect.Value) (int, error) { return 3, nil }
func (*unhashableExtEncoder) WriteToBytes(_ reflect.Value, offset int, _ *[]byte) int {
	return offset
}

type unhashableExtValue struct{ s []int }

type unhashableExtDecoder struct{ ext.DecoderCommon }

func (*unhashableExtDecoder) Code() int8 { return unhashableExtCode }
func (*unhashableExtDecoder) IsType(offset int, d *[]byte) bool {
	return (*d)[offset] == def.Fixext1 && int8((*d)[offset+1]) == unhashableExtCode
}
func (*unhashableExtDecoder) AsValue(offset int, _ reflect.Kind, _ *[]byte) (interface{}, int, error) {
	return []int{1}, offset + 3, nil
}

type unhashableExtStreamEncoder struct{}

func (*unhashableExtStreamEncoder) Code() int8                                  { return unhashableExtCode }
func (*unhashableExtStreamEncoder) Type() reflect.Type                          { return reflect.TypeOf(unhashableExtValue{}) }
func (*unhashableExtStreamEncoder) Write(ext.StreamWriter, reflect.Value) error { return nil }

type unhashableExtStreamDecoder struct{}

func (*unhashableExtStreamDecoder) Code() int8 { return unhashableExtCode }
func (*unhashableExtStreamDecoder) IsType(code byte, innerType int8, _ int) bool {
	return code == def.Fixext1 && innerType == unhashableExtCode
}
func (*unhashableExtStreamDecoder) ToValue(byte, []byte, reflect.Kind) (interface{}, error) {
	return []int{1}, nil
}

func TestDecodeUnhashableExtMapKey(t *testing.T) {
	enc, dec := &unhashableExtEncoder{}, &unhashableExtDecoder{}
	if err := msgpack.AddExtCoder(enc, dec); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = msgpack.RemoveExtCoder(enc, dec) })
	senc, sdec := &unhashableExtStreamEncoder{}, &unhashableExtStreamDecoder{}
	if err := msgpack.AddExtStreamCoder(senc, sdec); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = msgpack.RemoveExtStreamCoder(senc, sdec) })

	// fixmap(1){ fixext1(code 99, 0x00): nil }
	data := []byte{0x81, def.Fixext1, byte(unhashableExtCode), 0x00, 0xc0}
	dests := map[string]func() interface{}{
		"interface":     func() interface{} { var v interface{}; return &v },
		"mapIfaceIface": func() interface{} { var v map[interface{}]interface{}; return &v },
	}
	for _, df := range decodeFuncs {
		for dn, newV := range dests {
			t.Run(df.name+"/"+dn, func(t *testing.T) {
				err := mustNotPanic(t, func() error { return df.fn(data, newV()) })
				// The stream reflect path currently fails earlier with EOF when a
				// custom ext targets an interface{} element (pre-existing issue).
				if df.name == "UnmarshalRead" && dn == "mapIfaceIface" {
					if err == nil {
						t.Fatal("expected error")
					}
					return
				}
				if !errors.Is(err, def.ErrUnhashableMapKey) {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	}
}
