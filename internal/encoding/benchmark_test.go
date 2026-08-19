package encoding

import (
	"reflect"
	"testing"
	stdtime "time"

	"github.com/shamaton/msgpack/v3/def"
	"github.com/shamaton/msgpack/v3/ext"
)

type benchmarkStatus int
type benchmarkStatuses []benchmarkStatus

type benchmarkStruct struct {
	ID     int64
	Name   string
	Active bool
	Score  float64
	Tags   []string
	Data   []byte
}

type benchmarkIndirect struct {
	Pointer   *int
	Interface any
	Values    []any
	Nested    [][]int
}

var (
	benchmarkBytesResult []byte
	benchmarkPointerInt  = 42
	benchmarkIntSlice    = makeBenchmarkInts()
	benchmarkStatusSlice = makeBenchmarkStatuses()
	benchmarkStructValue = benchmarkStruct{
		ID: 42, Name: "benchmark", Active: true, Score: 12.5,
		Tags: []string{"a", "b", "c", "d"}, Data: make([]byte, 64),
	}
	benchmarkTimeValue     = stdtime.Unix(1_700_000_000, 123_456_789)
	benchmarkIndirectValue = benchmarkIndirect{
		Pointer:   &benchmarkPointerInt,
		Interface: int(42),
		Values:    []any{int(1), "two", []int{3, 4}},
		Nested:    [][]int{{1, 2}, {3, 4}},
	}
)

func makeBenchmarkInts() []int {
	values := make([]int, 1024)
	for i := range values {
		values[i] = i % 128
	}
	return values
}

func makeBenchmarkStatuses() []benchmarkStatus {
	values := make([]benchmarkStatus, 1024)
	for i := range values {
		values[i] = benchmarkStatus(i % 128)
	}
	return values
}

type benchmarkExtEncoder struct {
	ext.EncoderCommon
	typ reflect.Type
}

func (e *benchmarkExtEncoder) Code() int8         { return 42 }
func (e *benchmarkExtEncoder) Type() reflect.Type { return e.typ }
func (e *benchmarkExtEncoder) CalcByteSize(reflect.Value) (int, error) {
	return def.Byte1 + def.Byte1 + def.Byte1, nil
}
func (e *benchmarkExtEncoder) WriteToBytes(value reflect.Value, offset int, data *[]byte) int {
	offset = e.SetByte1Int(def.Fixext1, offset, data)
	offset = e.SetByte1Int(int(e.Code()), offset, data)
	return e.SetByte1Int64(valueToInt64(value), offset, data)
}

func valueToInt64(value reflect.Value) int64 {
	if value.Kind() == reflect.Slice {
		return int64(value.Len())
	}
	return value.Int()
}

func benchmarkEncode(b *testing.B, value any) {
	b.Helper()
	b.ReportAllocs()
	if _, err := Encode(value, false); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		encoded, err := Encode(value, false)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytesResult = encoded
	}
	b.StopTimer()
}

func benchmarkEncodeExt(b *testing.B, value any, encoder *benchmarkExtEncoder) {
	b.Helper()
	AddExtEncoder(encoder)
	defer RemoveExtEncoder(encoder)
	benchmarkEncode(b, value)
}

func BenchmarkEncodeInt(b *testing.B) {
	benchmarkEncode(b, int(42))
}

func BenchmarkEncodeIntSlice(b *testing.B) {
	benchmarkEncode(b, benchmarkIntSlice)
}

func BenchmarkEncodeStruct(b *testing.B) {
	benchmarkEncode(b, benchmarkStructValue)
}

func BenchmarkEncodeTime(b *testing.B) {
	benchmarkEncode(b, benchmarkTimeValue)
}

func BenchmarkEncodeIndirect(b *testing.B) {
	benchmarkEncode(b, benchmarkIndirectValue)
}

func BenchmarkEncodeStatus(b *testing.B) {
	encoder := &benchmarkExtEncoder{typ: reflect.TypeOf(benchmarkStatus(0))}
	benchmarkEncodeExt(b, benchmarkStatus(42), encoder)
}

func BenchmarkEncodeStatusSlice(b *testing.B) {
	encoder := &benchmarkExtEncoder{typ: reflect.TypeOf(benchmarkStatus(0))}
	benchmarkEncodeExt(b, benchmarkStatusSlice, encoder)
}

func BenchmarkEncodeStatuses(b *testing.B) {
	encoder := &benchmarkExtEncoder{typ: reflect.TypeOf(benchmarkStatuses{})}
	benchmarkEncodeExt(b, benchmarkStatuses(benchmarkStatusSlice), encoder)
}
