package msgpack_test

import (
	"bytes"
	"fmt"
	"net"
	"reflect"

	"github.com/shamaton/msgpack/v2"
	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
)

func ExampleAddExtCoder() {
	err := msgpack.AddExtCoder(&IPNetEncoder{}, &IPNetDecoder{})
	if err != nil {
		panic(err)
	}

	v1 := net.IPNet{IP: net.IP{127, 0, 0, 1}, Mask: net.IPMask{255, 255, 255, 0}}
	r1, err := msgpack.Marshal(v1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("% 02x\n", r1)
	// Output:
	// c7 0c 32 31 32 37 2e 30 2e 30 2e 31 2f 32 34

	var v2 net.IPNet
	err = msgpack.Unmarshal(r1, &v2)
	if err != nil {
		panic(err)
	}
	fmt.Println(v2)
	// Output:
	// {127.0.0.0 ffffff00}
}

func ExampleAddExtStreamCoder() {
	err := msgpack.AddExtStreamCoder(&IPNetStreamEncoder{}, &IPNetStreamDecoder{})
	if err != nil {
		panic(err)
	}

	v1 := net.IPNet{IP: net.IP{127, 0, 0, 1}, Mask: net.IPMask{255, 255, 255, 0}}

	buf := bytes.Buffer{}
	err = msgpack.MarshalWrite(&buf, v1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("% 02x\n", buf.Bytes())
	// Output:
	// c7 0c 32 31 32 37 2e 30 2e 30 2e 31 2f 32 34

	var v2 net.IPNet
	err = msgpack.UnmarshalRead(&buf, &v2)
	if err != nil {
		panic(err)
	}
	fmt.Println(v2)
	// Output:
	// {127.0.0.0 ffffff00}
}

const ipNetCode = 50

type IPNetDecoder struct {
	ext.DecoderCommon
}

var _ ext.Decoder = (*IPNetDecoder)(nil)

func (td *IPNetDecoder) Code() int8 {
	return ipNetCode
}

func (td *IPNetDecoder) IsType(offset int, d *[]byte) bool {
	code, offset := td.ReadSize1(offset, d)
	if code == def.Ext8 {
		_, offset = td.ReadSize1(offset, d)
		t, _ := td.ReadSize1(offset, d)
		return int8(t) == td.Code()
	}
	return false
}

func (td *IPNetDecoder) AsValue(offset int, k reflect.Kind, d *[]byte) (any, int, error) {
	code, offset := td.ReadSize1(offset, d)

	switch code {
	case def.Ext8:
		// size
		size, offset := td.ReadSize1(offset, d)
		// code
		_, offset = td.ReadSize1(offset, d)
		// value
		data, offset := td.ReadSizeN(offset, int(size), d)

		_, v, err := net.ParseCIDR(string(data))
		if err != nil {
			return net.IPNet{}, 0, fmt.Errorf("failed to parse CIDR: %w", err)
		}
		if v == nil {
			return net.IPNet{}, 0, fmt.Errorf("parsed CIDR is nil")
		}
		return *v, offset, nil
	}
	return net.IPNet{}, 0, fmt.Errorf("should not reach this line!! code %x decoding %v", td.Code(), k)
}

type IPNetStreamDecoder struct{}

var _ ext.StreamDecoder = (*IPNetStreamDecoder)(nil)

func (td *IPNetStreamDecoder) Code() int8 {
	return ipNetCode
}

func (td *IPNetStreamDecoder) IsType(code byte, innerType int8, _ int) bool {
	return code == def.Ext8 && innerType == td.Code()
}

func (td *IPNetStreamDecoder) ToValue(code byte, data []byte, k reflect.Kind) (any, error) {
	if code == def.Ext8 {
		_, v, err := net.ParseCIDR(string(data))
		if err != nil {
			return net.IPNet{}, fmt.Errorf("failed to parse CIDR: %w", err)
		}
		if v == nil {
			return net.IPNet{}, fmt.Errorf("parsed CIDR is nil")
		}
		return *v, nil
	}
	return net.IPNet{}, fmt.Errorf("should not reach this line!! code %x decoding %v", td.Code(), k)
}

type IPNetEncoder struct {
	ext.EncoderCommon
}

var _ ext.Encoder = (*IPNetEncoder)(nil)

func (s *IPNetEncoder) Code() int8 {
	return ipNetCode
}

func (s *IPNetEncoder) Type() reflect.Type {
	return reflect.TypeOf(net.IPNet{})
}

func (s *IPNetEncoder) CalcByteSize(value reflect.Value) (int, error) {
	v := value.Interface().(net.IPNet)
	fmt.Println(def.Byte1 + def.Byte1 + len([]byte(v.String())))
	return def.Byte1 + def.Byte1 + def.Byte1 + len(v.String()), nil
}

func (s *IPNetEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	v := value.Interface().(net.IPNet)
	data := v.String()

	offset = s.SetByte1Int(def.Ext8, offset, bytes)
	offset = s.SetByte1Int(len(data), offset, bytes)
	offset = s.SetByte1Int(int(s.Code()), offset, bytes)
	offset = s.SetBytes([]byte(data), offset, bytes)
	return offset
}

type IPNetStreamEncoder struct{}

var _ ext.StreamEncoder = (*IPNetStreamEncoder)(nil)

func (s *IPNetStreamEncoder) Code() int8 {
	return ipNetCode
}

func (s *IPNetStreamEncoder) Type() reflect.Type {
	return reflect.TypeOf(net.IPNet{})
}

func (s *IPNetStreamEncoder) Write(w ext.StreamWriter, value reflect.Value) error {
	v := value.Interface().(net.IPNet)
	data := v.String()

	if err := w.WriteByte1Int(def.Ext8); err != nil {
		return err
	}
	if err := w.WriteByte1Int(len(data)); err != nil {
		return err
	}
	if err := w.WriteByte1Int(int(s.Code())); err != nil {
		return err
	}
	if err := w.WriteBytes([]byte(data)); err != nil {
		return err
	}

	return nil
}
