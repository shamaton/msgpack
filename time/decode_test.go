package time

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3/def"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func TestDecodeCode(t *testing.T) {
	decoder := Decoder
	code := decoder.Code()
	tu.Equal(t, code, def.TimeStamp)
}

func TestDecodeIsType(t *testing.T) {
	decoder := Decoder
	ts := def.TimeStamp

	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "Fixext4 with TimeStamp",
			data:     []byte{def.Fixext4, byte(ts), 0, 0, 0, 0},
			expected: true,
		},
		{
			name:     "Fixext4 with wrong type",
			data:     []byte{def.Fixext4, 0x00, 0, 0, 0, 0},
			expected: false,
		},
		{
			name:     "Fixext8 with TimeStamp",
			data:     []byte{def.Fixext8, byte(ts), 0, 0, 0, 0, 0, 0, 0, 0},
			expected: true,
		},
		{
			name:     "Fixext8 with wrong type",
			data:     []byte{def.Fixext8, 0x00, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: false,
		},
		{
			name:     "Ext8 with length 12 and TimeStamp",
			data:     []byte{def.Ext8, 12, byte(ts), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: true,
		},
		{
			name:     "Ext8 with wrong length",
			data:     []byte{def.Ext8, 10, byte(ts), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: false,
		},
		{
			name:     "Ext8 with wrong type",
			data:     []byte{def.Ext8, 12, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: false,
		},
		{
			name:     "Wrong format",
			data:     []byte{def.Nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decoder.IsType(0, &tt.data)
			tu.Equal(t, result, tt.expected)
		})
	}
}

func TestDecodeAsValueFixext4(t *testing.T) {
	decoder := Decoder
	ts := def.TimeStamp

	tests := []struct {
		name     string
		unixTime int64
	}{
		{
			name:     "Unix epoch",
			unixTime: 0,
		},
		{
			name:     "Small timestamp",
			unixTime: 1000,
		},
		{
			name:     "Large timestamp",
			unixTime: 1<<32 - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fixext4 format data
			data := make([]byte, 6)
			data[0] = def.Fixext4
			data[1] = byte(ts)
			binary.BigEndian.PutUint32(data[2:], uint32(tt.unixTime))

			value, offset, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
			tu.NoError(t, err)
			tu.Equal(t, offset, 6)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			tu.Equal(t, timeValue.Unix(), tt.unixTime)
			tu.Equal(t, timeValue.Nanosecond(), 0)
		})
	}
}

func TestDecodeAsValueFixext8(t *testing.T) {
	decoder := Decoder
	ts := def.TimeStamp

	tests := []struct {
		name       string
		unixTime   int64
		nanosecond int64
	}{
		{
			name:       "With small nanoseconds",
			unixTime:   1000,
			nanosecond: 123456789,
		},
		{
			name:       "With max nanoseconds",
			unixTime:   67890,
			nanosecond: 999999999,
		},
		{
			name:       "Boundary at 2^34-1",
			unixTime:   (1 << 34) - 1,
			nanosecond: 999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fixext8 format data
			data := make([]byte, 10)
			data[0] = def.Fixext8
			data[1] = byte(ts)

			// Pack nanoseconds in upper 30 bits and seconds in lower 34 bits
			data64 := uint64(tt.nanosecond)<<34 | uint64(tt.unixTime)
			binary.BigEndian.PutUint64(data[2:], data64)

			value, offset, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
			tu.NoError(t, err)
			tu.Equal(t, offset, 10)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			tu.Equal(t, timeValue.Unix(), tt.unixTime)
			tu.Equal(t, int64(timeValue.Nanosecond()), tt.nanosecond)
		})
	}
}

func TestDecodeAsValueExt8(t *testing.T) {
	decoder := Decoder
	ts := def.TimeStamp

	tests := []struct {
		name       string
		unixTime   int64
		nanosecond int32
	}{
		{
			name:       "Large timestamp at 2^34",
			unixTime:   1 << 34,
			nanosecond: 123456789,
		},
		{
			name:       "Very large timestamp",
			unixTime:   253402300799,
			nanosecond: 999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Ext8 format data
			data := make([]byte, 15)
			data[0] = def.Ext8
			data[1] = 12 // length
			data[2] = byte(ts)
			binary.BigEndian.PutUint32(data[3:], uint32(tt.nanosecond))
			binary.BigEndian.PutUint64(data[7:], uint64(tt.unixTime))

			value, offset, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
			tu.NoError(t, err)
			tu.Equal(t, offset, 15)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			tu.Equal(t, timeValue.Unix(), tt.unixTime)
			tu.Equal(t, int64(timeValue.Nanosecond()), int64(tt.nanosecond))
		})
	}
}

func TestDecodeAsValueErrors(t *testing.T) {
	decoder := Decoder
	ts := def.TimeStamp

	t.Run("Fixext8 with invalid nanoseconds", func(t *testing.T) {
		data := make([]byte, 10)
		data[0] = def.Fixext8
		data[1] = byte(ts)

		// Set nanoseconds > 999999999 (invalid)
		invalidNano := uint64(1000000000)
		data64 := invalidNano<<34 | 1000
		binary.BigEndian.PutUint64(data[2:], data64)

		_, _, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
		tu.ErrorContains(t, err, "in timestamp 64 formats")
	})

	t.Run("Ext8 with invalid nanoseconds", func(t *testing.T) {
		data := make([]byte, 15)
		data[0] = def.Ext8
		data[1] = 12
		data[2] = byte(ts)

		// Set nanoseconds > 999999999 (invalid)
		binary.BigEndian.PutUint32(data[3:], 1000000000)
		binary.BigEndian.PutUint64(data[7:], 1000)

		_, _, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
		tu.ErrorContains(t, err, "in timestamp 96 formats")
	})

	t.Run("Invalid format code", func(t *testing.T) {
		data := []byte{def.Nil}

		_, _, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
		tu.ErrorContains(t, err, "should not reach")
	})
}

func TestDecodeTimezone(t *testing.T) {
	decoder := Decoder
	ts := def.TimeStamp

	tests := []struct {
		name       string
		time       time.Time
		createData func(time.Time) []byte
	}{
		{
			name: "Fixext4",
			time: time.Unix(1000, 0),
			createData: func(testTime time.Time) []byte {
				data := make([]byte, 6)
				data[0] = def.Fixext4
				data[1] = byte(ts)
				binary.BigEndian.PutUint32(data[2:], uint32(testTime.Unix()))
				return data
			},
		},
		{
			name: "Fixext8",
			time: time.Unix(1000, 123456789),
			createData: func(testTime time.Time) []byte {
				data := make([]byte, 10)
				data[0] = def.Fixext8
				data[1] = byte(ts)
				data64 := uint64(testTime.Nanosecond())<<34 | uint64(testTime.Unix())
				binary.BigEndian.PutUint64(data[2:], data64)
				return data
			},
		},
		{
			name: "Ext8",
			time: time.Unix(1<<34, 123456789),
			createData: func(testTime time.Time) []byte {
				data := make([]byte, 15)
				data[0] = def.Ext8
				data[1] = 12
				data[2] = byte(ts)
				binary.BigEndian.PutUint32(data[3:], uint32(testTime.Nanosecond()))
				binary.BigEndian.PutUint64(data[7:], uint64(testTime.Unix()))
				return data
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" - Decode as local (default)", func(t *testing.T) {
			// Set to local timezone
			SetDecodedAsLocal(true)
			defer SetDecodedAsLocal(true) // Reset to default

			data := tt.createData(tt.time)
			value, _, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
			tu.NoError(t, err)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			// Should be in local timezone
			tu.Equal(t, timeValue.Location(), time.Local)
		})

		t.Run(tt.name+" - Decode as UTC", func(t *testing.T) {
			// Set to UTC timezone
			SetDecodedAsLocal(false)
			defer SetDecodedAsLocal(true) // Reset to default

			data := tt.createData(tt.time)
			value, _, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &data)
			tu.NoError(t, err)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			// Should be in UTC timezone
			tu.Equal(t, timeValue.Location(), time.UTC)
		})
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	encoder := Encoder
	decoder := Decoder

	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "Fixext4 format",
			time: time.Unix(1000, 0),
		},
		{
			name: "Fixext8 format",
			time: time.Unix(67890, 123456789),
		},
		{
			name: "Ext8 format",
			time: time.Unix(17179869184, 987654321),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			value := reflect.ValueOf(tt.time)
			size, err := encoder.CalcByteSize(value)
			tu.NoError(t, err)

			bytes := make([]byte, size)
			encoder.WriteToBytes(value, 0, &bytes)

			// Decode
			decodedValue, offset, err := decoder.AsValue(0, reflect.TypeOf(time.Time{}).Kind(), &bytes)
			tu.NoError(t, err)
			tu.Equal(t, offset, size)

			decodedTime, ok := decodedValue.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", decodedValue)
			}

			// Compare Unix time and nanoseconds
			tu.Equal(t, decodedTime.Unix(), tt.time.Unix())
			tu.Equal(t, decodedTime.Nanosecond(), tt.time.Nanosecond())
		})
	}
}
