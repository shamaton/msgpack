package msgpack

import (
	"fmt"
	"io"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	"github.com/shamaton/msgpack/v2/internal/decoding"
	"github.com/shamaton/msgpack/v2/internal/encoding"
	d2 "github.com/shamaton/msgpack/v2/internal/io/decoding"
	e2 "github.com/shamaton/msgpack/v2/internal/io/encoding"
	d3 "github.com/shamaton/msgpack/v2/internal/io2/decoding"
)

// StructAsArray is encoding option.
// If this option sets true, default encoding sets to array-format.
var StructAsArray = false

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v interface{}) ([]byte, error) {
	return encoding.Encode(v, StructAsArray)
}

func MarshalWrite(w io.Writer, v interface{}) error {
	return e2.Encode(w, v, StructAsArray)
}

// Unmarshal analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Unmarshal(data []byte, v interface{}) error {
	return decoding.Decode(data, v, StructAsArray)
}

func UnmarshalRead(r io.Reader, v interface{}) error {
	return d2.Decode(r, v, StructAsArray)
}

func UnmarshalRead2(r io.Reader, v interface{}) error {
	return d3.Decode(r, v, StructAsArray)
}

// AddExtCoder adds encoders for extension types.
func AddExtCoder(e ext.Encoder, d ext.Decoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	encoding.AddExtEncoder(e)
	decoding.AddExtDecoder(d)
	return nil
}

// AddExtStreamCoder adds encoders for extension types.
func AddExtStreamCoder(e ext.StreamEncoder, d ext.StreamDecoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	e2.AddExtEncoder(e)
	d2.AddExtDecoder(d)
	return nil
}

// RemoveExtCoder removes encoders for extension types.
func RemoveExtCoder(e ext.Encoder, d ext.Decoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	encoding.RemoveExtEncoder(e)
	decoding.RemoveExtDecoder(d)
	return nil
}

// RemoveExtStreamCoder removes encoders for extension types.
func RemoveExtStreamCoder(e ext.StreamEncoder, d ext.StreamDecoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	e2.RemoveExtEncoder(e)
	d2.RemoveExtDecoder(d)
	return nil
}

// SetComplexTypeCode sets def.complexTypeCode
func SetComplexTypeCode(code int8) {
	def.SetComplexTypeCode(code)
}
