package msgpack

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	"github.com/shamaton/msgpack/v2/internal/decoding"
	"github.com/shamaton/msgpack/v2/internal/encoding"
)

// StructAsArray is encoding option.
// If this option sets true, default encoding sets to array-format.
var StructAsArray = false

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v interface{}) ([]byte, error) {
	return encoding.EncodeBytes(v, StructAsArray)
}

// Unmarshal analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Unmarshal(data []byte, v interface{}) error {
	return decoding.DecodeBytes(data, v, StructAsArray)
}

type Encoder struct {
	writer io.Writer

	StructAsArray bool
}

// NewEncoder will create an Encoder that will encode values into a stream
func NewEncoder(output io.Writer) Encoder {
	return Encoder{writer: output}
}

// Encode a single value into the output stream
func (e Encoder) Encode(v interface{}) error {
	return encoding.Encode(v, e.writer, e.StructAsArray)
}

func (e Encoder) WithStructAsArray(structAsArray bool) Encoder {
	e.StructAsArray = structAsArray
	return e
}

type Decoder struct {
	reader *bufio.Reader

	StructAsArray bool
	InternStrings bool
}

// NewDecoder will create an Decoder that will decode values from a stream one at a time
func NewDecoder(input io.Reader) Decoder {
	// if input is already a bufio.Reader, just type assert it
	bufReader, ok := input.(*bufio.Reader)
	if !ok {
		// otherwise, wrap the input in a bufio reader
		bufReader = bufio.NewReader(input)
	}

	return Decoder{reader: bufReader}
}

func (d Decoder) Decode(v interface{}) error {
	return decoding.Decode(d.reader, v, d.StructAsArray, d.InternStrings)
}

func (d Decoder) WithStructAsArray(structAsArray bool) Decoder {
	d.StructAsArray = structAsArray
	return d
}

func (d Decoder) WithInternStrings(internStrings bool) Decoder {
	d.InternStrings = internStrings
	return d
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

// RemoveExtCoder removes encoders for extension types.
func RemoveExtCoder(e ext.Encoder, d ext.Decoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	encoding.RemoveExtEncoder(e)
	decoding.RemoveExtDecoder(d)
	return nil
}

// SetComplexTypeCode sets def.complexTypeCode
func SetComplexTypeCode(code int8) {
	def.SetComplexTypeCode(code)
}
