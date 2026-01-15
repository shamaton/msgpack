package time

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3/def"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func TestStreamDecodeCode(t *testing.T) {
	decoder := StreamDecoder
	code := decoder.Code()
	tu.Equal(t, code, def.TimeStamp)
}

func TestStreamDecodeIsType(t *testing.T) {
	decoder := StreamDecoder
	ts := int8(def.TimeStamp)

	tests := []struct {
		name       string
		code       byte
		innerType  int8
		dataLength int
		expected   bool
	}{
		{
			name:       "Fixext4 with TimeStamp",
			code:       def.Fixext4,
			innerType:  ts,
			dataLength: 4,
			expected:   true,
		},
		{
			name:       "Fixext4 with wrong type",
			code:       def.Fixext4,
			innerType:  0x00,
			dataLength: 4,
			expected:   false,
		},
		{
			name:       "Fixext8 with TimeStamp",
			code:       def.Fixext8,
			innerType:  ts,
			dataLength: 8,
			expected:   true,
		},
		{
			name:       "Fixext8 with wrong type",
			code:       def.Fixext8,
			innerType:  0x00,
			dataLength: 8,
			expected:   false,
		},
		{
			name:       "Ext8 with length 12 and TimeStamp",
			code:       def.Ext8,
			innerType:  ts,
			dataLength: 12,
			expected:   true,
		},
		{
			name:       "Ext8 with wrong length",
			code:       def.Ext8,
			innerType:  ts,
			dataLength: 10,
			expected:   false,
		},
		{
			name:       "Ext8 with wrong type",
			code:       def.Ext8,
			innerType:  0x00,
			dataLength: 12,
			expected:   false,
		},
		{
			name:       "Wrong format",
			code:       def.Nil,
			innerType:  ts,
			dataLength: 0,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decoder.IsType(tt.code, tt.innerType, tt.dataLength)
			tu.Equal(t, result, tt.expected)
		})
	}
}

func TestStreamDecodeToValueFixext4(t *testing.T) {
	decoder := StreamDecoder

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
			// Create Fixext4 format data (without header)
			data := make([]byte, 4)
			binary.BigEndian.PutUint32(data, uint32(tt.unixTime))

			value, err := decoder.ToValue(def.Fixext4, data, reflect.TypeOf(time.Time{}).Kind())
			tu.NoError(t, err)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			tu.Equal(t, timeValue.Unix(), tt.unixTime)
			tu.Equal(t, timeValue.Nanosecond(), 0)
		})
	}
}

func TestStreamDecodeToValueFixext8(t *testing.T) {
	decoder := StreamDecoder

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
			// Create Fixext8 format data (without header)
			data := make([]byte, 8)

			// Pack nanoseconds in upper 30 bits and seconds in lower 34 bits
			data64 := uint64(tt.nanosecond)<<34 | uint64(tt.unixTime)
			binary.BigEndian.PutUint64(data, data64)

			value, err := decoder.ToValue(def.Fixext8, data, reflect.TypeOf(time.Time{}).Kind())
			tu.NoError(t, err)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			tu.Equal(t, timeValue.Unix(), tt.unixTime)
			tu.Equal(t, int64(timeValue.Nanosecond()), tt.nanosecond)
		})
	}
}

func TestStreamDecodeToValueExt8(t *testing.T) {
	decoder := StreamDecoder

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
			// Create Ext8 format data (without header)
			data := make([]byte, 12)
			binary.BigEndian.PutUint32(data[:4], uint32(tt.nanosecond))
			binary.BigEndian.PutUint64(data[4:12], uint64(tt.unixTime))

			value, err := decoder.ToValue(def.Ext8, data, reflect.TypeOf(time.Time{}).Kind())
			tu.NoError(t, err)

			timeValue, ok := value.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", value)
			}

			tu.Equal(t, timeValue.Unix(), tt.unixTime)
			tu.Equal(t, int64(timeValue.Nanosecond()), int64(tt.nanosecond))
		})
	}
}

func TestStreamDecodeToValueErrors(t *testing.T) {
	decoder := StreamDecoder

	t.Run("Fixext8 with invalid nanoseconds", func(t *testing.T) {
		data := make([]byte, 8)

		// Set nanoseconds > 999999999 (invalid)
		invalidNano := uint64(1000000000)
		data64 := invalidNano<<34 | 1000
		binary.BigEndian.PutUint64(data, data64)

		_, err := decoder.ToValue(def.Fixext8, data, reflect.TypeOf(time.Time{}).Kind())
		tu.ErrorContains(t, err, "in timestamp 64 formats")
	})

	t.Run("Ext8 with invalid nanoseconds", func(t *testing.T) {
		data := make([]byte, 12)

		// Set nanoseconds > 999999999 (invalid)
		binary.BigEndian.PutUint32(data[:4], 1000000000)
		binary.BigEndian.PutUint64(data[4:12], 1000)

		_, err := decoder.ToValue(def.Ext8, data, reflect.TypeOf(time.Time{}).Kind())
		tu.ErrorContains(t, err, "in timestamp 96 formats")
	})

	t.Run("Invalid format code", func(t *testing.T) {
		data := []byte{0}

		_, err := decoder.ToValue(def.Nil, data, reflect.TypeOf(time.Time{}).Kind())
		tu.ErrorContains(t, err, "should not reach")
	})
}

func TestStreamDecodeTimezone(t *testing.T) {
	decoder := StreamDecoder

	tests := []struct {
		name       string
		time       time.Time
		format     byte
		createData func(time.Time) []byte
	}{
		{
			name:   "Fixext4",
			time:   time.Unix(1000, 0),
			format: def.Fixext4,
			createData: func(testTime time.Time) []byte {
				data := make([]byte, 4)
				binary.BigEndian.PutUint32(data, uint32(testTime.Unix()))
				return data
			},
		},
		{
			name:   "Fixext8",
			time:   time.Unix(1000, 123456789),
			format: def.Fixext8,
			createData: func(testTime time.Time) []byte {
				data := make([]byte, 8)
				data64 := uint64(testTime.Nanosecond())<<34 | uint64(testTime.Unix())
				binary.BigEndian.PutUint64(data, data64)
				return data
			},
		},
		{
			name:   "Ext8",
			time:   time.Unix(1<<34, 123456789),
			format: def.Ext8,
			createData: func(testTime time.Time) []byte {
				data := make([]byte, 12)
				binary.BigEndian.PutUint32(data[:4], uint32(testTime.Nanosecond()))
				binary.BigEndian.PutUint64(data[4:12], uint64(testTime.Unix()))
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
			value, err := decoder.ToValue(tt.format, data, reflect.TypeOf(time.Time{}).Kind())
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
			value, err := decoder.ToValue(tt.format, data, reflect.TypeOf(time.Time{}).Kind())
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

func TestStreamDecodeRoundTrip(t *testing.T) {
	decoder := StreamDecoder

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
			// Use encode_stream to encode
			var encodedData []byte
			secs := uint64(tt.time.Unix())
			if secs>>34 == 0 {
				data := uint64(tt.time.Nanosecond())<<34 | secs
				if data&0xffffffff00000000 == 0 {
					// Fixext4
					encodedData = make([]byte, 4)
					binary.BigEndian.PutUint32(encodedData, uint32(data))

					// Decode
					decodedValue, err := decoder.ToValue(def.Fixext4, encodedData, reflect.TypeOf(time.Time{}).Kind())
					tu.NoError(t, err)

					decodedTime, ok := decodedValue.(time.Time)
					if !ok {
						t.Fatalf("Expected time.Time, got %T", decodedValue)
					}

					tu.Equal(t, decodedTime.Unix(), tt.time.Unix())
					tu.Equal(t, decodedTime.Nanosecond(), tt.time.Nanosecond())
					return
				}

				// Fixext8
				encodedData = make([]byte, 8)
				binary.BigEndian.PutUint64(encodedData, data)

				// Decode
				decodedValue, err := decoder.ToValue(def.Fixext8, encodedData, reflect.TypeOf(time.Time{}).Kind())
				tu.NoError(t, err)

				decodedTime, ok := decodedValue.(time.Time)
				if !ok {
					t.Fatalf("Expected time.Time, got %T", decodedValue)
				}

				tu.Equal(t, decodedTime.Unix(), tt.time.Unix())
				tu.Equal(t, decodedTime.Nanosecond(), tt.time.Nanosecond())
				return
			}

			// Ext8
			encodedData = make([]byte, 12)
			binary.BigEndian.PutUint32(encodedData[:4], uint32(tt.time.Nanosecond()))
			binary.BigEndian.PutUint64(encodedData[4:12], secs)

			// Decode
			decodedValue, err := decoder.ToValue(def.Ext8, encodedData, reflect.TypeOf(time.Time{}).Kind())
			tu.NoError(t, err)

			decodedTime, ok := decodedValue.(time.Time)
			if !ok {
				t.Fatalf("Expected time.Time, got %T", decodedValue)
			}

			tu.Equal(t, decodedTime.Unix(), tt.time.Unix())
			tu.Equal(t, decodedTime.Nanosecond(), tt.time.Nanosecond())
		})
	}
}
