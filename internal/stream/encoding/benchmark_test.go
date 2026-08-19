package encoding

import (
	"bytes"
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
	benchmarkStreamLength int
	benchmarkPointerInt   = 42
	benchmarkIntSlice     = makeBenchmarkInts()
	benchmarkStatusSlice  = makeBenchmarkStatuses()
	benchmarkStructValue  = benchmarkStruct{
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
	typ reflect.Type
}

func (e *benchmarkExtEncoder) Code() int8         { return 42 }
func (e *benchmarkExtEncoder) Type() reflect.Type { return e.typ }
func (e *benchmarkExtEncoder) Write(w ext.StreamWriter, value reflect.Value) error {
	if err := w.WriteByte1Int(def.Fixext1); err != nil {
		return err
	}
	if err := w.WriteByte1Int(int(e.Code())); err != nil {
		return err
	}
	return w.WriteByte1Int64(valueToInt64(value))
}

func valueToInt64(value reflect.Value) int64 {
	if value.Kind() == reflect.Slice {
		return int64(value.Len())
	}
	return value.Int()
}

func benchmarkStreamEncode(b *testing.B, value any) {
	b.Helper()
	var buffer bytes.Buffer
	if err := Encode(&buffer, value, false); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		buffer.Reset()
		if err := Encode(&buffer, value, false); err != nil {
			b.Fatal(err)
		}
		benchmarkStreamLength = buffer.Len()
	}
	b.StopTimer()
}

func benchmarkStreamEncodeExt(b *testing.B, value any, encoder *benchmarkExtEncoder) {
	b.Helper()
	AddExtEncoder(encoder)
	defer RemoveExtEncoder(encoder)
	benchmarkStreamEncode(b, value)
}

func BenchmarkStreamEncodeInt(b *testing.B) {
	benchmarkStreamEncode(b, int(42))
}

func BenchmarkStreamEncodeIntSlice(b *testing.B) {
	benchmarkStreamEncode(b, benchmarkIntSlice)
}

func BenchmarkStreamEncodeStruct(b *testing.B) {
	benchmarkStreamEncode(b, benchmarkStructValue)
}

func BenchmarkStreamEncodeTime(b *testing.B) {
	benchmarkStreamEncode(b, benchmarkTimeValue)
}

func BenchmarkStreamEncodeIndirect(b *testing.B) {
	benchmarkStreamEncode(b, benchmarkIndirectValue)
}

func BenchmarkStreamEncodeStatus(b *testing.B) {
	encoder := &benchmarkExtEncoder{typ: reflect.TypeOf(benchmarkStatus(0))}
	benchmarkStreamEncodeExt(b, benchmarkStatus(42), encoder)
}

func BenchmarkStreamEncodeStatusSlice(b *testing.B) {
	encoder := &benchmarkExtEncoder{typ: reflect.TypeOf(benchmarkStatus(0))}
	benchmarkStreamEncodeExt(b, benchmarkStatusSlice, encoder)
}

func BenchmarkStreamEncodeStatuses(b *testing.B) {
	encoder := &benchmarkExtEncoder{typ: reflect.TypeOf(benchmarkStatuses{})}
	benchmarkStreamEncodeExt(b, benchmarkStatuses(benchmarkStatusSlice), encoder)
}
