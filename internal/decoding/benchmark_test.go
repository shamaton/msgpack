package decoding

import (
	"reflect"
	"testing"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
	"github.com/shamaton/msgpack/v3/internal/encoding"
)

// benchStatus is a named non-struct type, as commonly used for "enums" in Go.
type benchStatus int

type benchStatusEncoder struct {
	ext.EncoderCommon
}

func (e *benchStatusEncoder) Code() int8         { return 44 }
func (e *benchStatusEncoder) Type() reflect.Type { return reflect.TypeOf(benchStatus(0)) }
func (e *benchStatusEncoder) CalcByteSize(reflect.Value) (int, error) {
	return def.Byte1 + def.Byte1 + def.Byte1, nil
}
func (e *benchStatusEncoder) WriteToBytes(value reflect.Value, offset int, data *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, data)
	offset = e.SetByte1Int(int(e.Code()), offset, data)
	return e.SetByte1Int64(value.Int(), offset, data)
}

type benchStatusDecoder struct{}

func (*benchStatusDecoder) Code() int8 { return 44 }
func (*benchStatusDecoder) IsType(offset int, data *[]byte) bool {
	return (*data)[offset] == def.Fixext1 && int8((*data)[offset+1]) == 44
}
func (*benchStatusDecoder) AsValue(offset int, _ reflect.Kind, data *[]byte) (interface{}, int, error) {
	return benchStatus(int8((*data)[offset+2])), offset + 3, nil
}

func makeBenchStatusSlice() []benchStatus {
	values := make([]benchStatus, 1024)
	for i := range values {
		values[i] = benchStatus(i % 128)
	}
	return values
}

func benchmarkDecode(b *testing.B, encoded []byte, target func() any) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		v := target()
		if err := Decode(encoded, v, false); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

// BenchmarkDecodeInt measures the common non-ext path: it must not regress
// once ext dispatch is consulted generically in decode(), not just setStruct.
func BenchmarkDecodeInt(b *testing.B) {
	encoded, err := encoding.Encode(int(42), false)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkDecode(b, encoded, func() any { var v int; return &v })
}

func BenchmarkDecodeIntSlice(b *testing.B) {
	values := make([]int, 1024)
	for i := range values {
		values[i] = i % 128
	}
	encoded, err := encoding.Encode(values, false)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkDecode(b, encoded, func() any { var v []int; return &v })
}

type benchStruct struct {
	ID     int64
	Name   string
	Active bool
	Score  float64
	Tags   []string
}

func BenchmarkDecodeStruct(b *testing.B) {
	value := benchStruct{ID: 42, Name: "benchmark", Active: true, Score: 12.5, Tags: []string{"a", "b", "c", "d"}}
	encoded, err := encoding.Encode(value, false)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkDecode(b, encoded, func() any { var v benchStruct; return &v })
}

func BenchmarkDecodeStatus(b *testing.B) {
	encoder := &benchStatusEncoder{}
	decoder := &benchStatusDecoder{}
	encoding.AddExtEncoder(encoder)
	defer encoding.RemoveExtEncoder(encoder)
	AddExtDecoder(decoder)
	defer RemoveExtDecoder(decoder)

	encoded, err := encoding.Encode(benchStatus(42), false)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkDecode(b, encoded, func() any { var v benchStatus; return &v })
}

func BenchmarkDecodeStatusSlice(b *testing.B) {
	encoder := &benchStatusEncoder{}
	decoder := &benchStatusDecoder{}
	encoding.AddExtEncoder(encoder)
	defer encoding.RemoveExtEncoder(encoder)
	AddExtDecoder(decoder)
	defer RemoveExtDecoder(decoder)

	encoded, err := encoding.Encode(makeBenchStatusSlice(), false)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkDecode(b, encoded, func() any { var v []benchStatus; return &v })
}
