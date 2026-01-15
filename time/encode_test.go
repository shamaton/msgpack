package time

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v3/def"
	tu "github.com/shamaton/msgpack/v3/internal/common/testutil"
)

func TestCode(t *testing.T) {
	encoder := Encoder
	code := encoder.Code()
	tu.Equal(t, code, def.TimeStamp)
}

func TestType(t *testing.T) {
	encoder := Encoder
	typ := encoder.Type()
	expected := reflect.TypeOf(time.Time{})
	tu.Equal(t, typ, expected)
}

func TestCalcByteSize(t *testing.T) {
	encoder := Encoder

	tests := []struct {
		name     string
		time     time.Time
		expected int
	}{
		{
			name:     "Fixext4 - epoch",
			time:     time.Unix(0, 0),
			expected: def.Byte1 + def.Byte1 + def.Byte4, // 6
		},
		{
			name:     "Fixext4 - small timestamp",
			time:     time.Unix(1000, 0),
			expected: def.Byte1 + def.Byte1 + def.Byte4, // 6
		},
		{
			name:     "Fixext8 - with nanoseconds",
			time:     time.Unix(1000, 999999999),
			expected: def.Byte1 + def.Byte1 + def.Byte8, // 10
		},
		{
			name:     "Fixext8 - boundary (2^34-1)",
			time:     time.Unix((1<<34)-1, 0),
			expected: def.Byte1 + def.Byte1 + def.Byte8, // 10
		},
		{
			name:     "Ext8 - large timestamp (2^34)",
			time:     time.Unix(1<<34, 0),
			expected: def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8, // 15
		},
		{
			name:     "Ext8 - large timestamp with nanoseconds",
			time:     time.Unix(1<<34, 123456789),
			expected: def.Byte1 + def.Byte1 + def.Byte1 + def.Byte4 + def.Byte8, // 15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.time)
			size, err := encoder.CalcByteSize(value)
			tu.NoError(t, err)
			tu.Equal(t, size, tt.expected)
		})
	}
}

func TestEncodedDataAccuracy(t *testing.T) {
	encoder := Encoder

	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "Fixext4 - simple",
			time: time.Unix(12345, 0),
		},
		{
			name: "Fixext8 - with nanoseconds",
			time: time.Unix(67890, 123456789),
		},
		{
			name: "Ext8 - large",
			time: time.Unix(17179869184, 987654321),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.time)
			size, err := encoder.CalcByteSize(value)
			tu.NoError(t, err)

			bytes := make([]byte, size)
			encoder.WriteToBytes(value, 0, &bytes)

			// Verify we can extract the correct time back
			switch bytes[0] {
			case def.Fixext4:
				data := binary.BigEndian.Uint32(bytes[2:6])
				secs := int64(data)
				tu.Equal(t, secs, tt.time.Unix())

			case def.Fixext8:
				data := binary.BigEndian.Uint64(bytes[2:10])
				nano := int64(data >> 34)
				secs := int64(data & 0x00000003ffffffff)
				tu.Equal(t, secs, tt.time.Unix())
				tu.Equal(t, nano, int64(tt.time.Nanosecond()))

			case def.Ext8:
				nano := binary.BigEndian.Uint32(bytes[3:7])
				secs := binary.BigEndian.Uint64(bytes[7:15])
				tu.Equal(t, int64(secs), tt.time.Unix())
				tu.Equal(t, int64(nano), int64(tt.time.Nanosecond()))
			}
		})
	}
}

func TestWriteToBytes(t *testing.T) {
	tests := []struct {
		name           string
		time           time.Time
		expectedLen    int
		expectedFormat string
	}{
		{
			name:           "Fixext4 format (32-bit timestamp, no nanoseconds)",
			time:           time.Unix(0, 0),
			expectedLen:    6, // 1 (Fixext4) + 1 (TimeStamp) + 4 (data)
			expectedFormat: "fixext4",
		},
		{
			name:           "Fixext4 format (small timestamp, no nanoseconds)",
			time:           time.Unix(1000, 0),
			expectedLen:    6,
			expectedFormat: "fixext4",
		},
		{
			name:           "Fixext8 format (needs 64-bit but secs fits in 34 bits)",
			time:           time.Unix(17179869183, 999999999), // (1 << 34) - 1
			expectedLen:    10,                                // 1 (Fixext8) + 1 (TimeStamp) + 8 (data)
			expectedFormat: "fixext8",
		},
		{
			name:           "Ext8 format (secs >= 2^34)",
			time:           time.Unix(17179869184, 123456789), // 1 << 34
			expectedLen:    15,                                // 1 (Ext8) + 1 (12) + 1 (TimeStamp) + 4 (nsec) + 8 (secs)
			expectedFormat: "ext8",
		},
		{
			name:           "Ext8 format (large timestamp)",
			time:           time.Unix(253402300799, 999999999), // 9999-12-31 23:59:59
			expectedLen:    15,
			expectedFormat: "ext8",
		},
	}

	encoder := Encoder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.time)

			// Calculate expected byte size
			size, err := encoder.CalcByteSize(value)
			tu.NoError(t, err)
			tu.Equal(t, size, tt.expectedLen)

			// Create byte slice
			bytes := make([]byte, size)

			// Write to bytes
			offset := encoder.WriteToBytes(value, 0, &bytes)
			tu.Equal(t, offset, tt.expectedLen)

			// Verify format type
			switch tt.expectedFormat {
			case "fixext4":
				tu.Equal(t, bytes[0], def.Fixext4)
				tu.Equal(t, int8(bytes[1]), def.TimeStamp)

			case "fixext8":
				tu.Equal(t, bytes[0], def.Fixext8)
				tu.Equal(t, int8(bytes[1]), def.TimeStamp)

			case "ext8":
				tu.Equal(t, bytes[0], def.Ext8)
				tu.Equal(t, bytes[1], 12)
				tu.Equal(t, int8(bytes[2]), def.TimeStamp)

			default:
				t.Errorf("Unknown expected format: %s", tt.expectedFormat)
			}
		})
	}
}

func TestWriteToBytesOffset(t *testing.T) {
	encoder := Encoder
	testTime := time.Unix(1000000000, 123456789)
	value := reflect.ValueOf(testTime)

	size, err := encoder.CalcByteSize(value)
	tu.NoError(t, err)

	// Test with non-zero offset
	offset := 10
	bytes := make([]byte, offset+size)

	newOffset := encoder.WriteToBytes(value, offset, &bytes)
	tu.Equal(t, newOffset, offset+size)

	// Verify the data is written at correct position
	if bytes[offset] != def.Fixext4 && bytes[offset] != def.Fixext8 && bytes[offset] != def.Ext8 {
		t.Errorf("Data not written at correct offset")
	}
}

func TestWriteToBytesEdgeCases(t *testing.T) {
	encoder := Encoder

	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "Unix epoch",
			time: time.Unix(0, 0),
		},
		{
			name: "Maximum nanoseconds",
			time: time.Unix(1000, 999999999),
		},
		{
			name: "Boundary at 2^34 - 1",
			time: time.Unix((1<<34)-1, 0),
		},
		{
			name: "Boundary at 2^34",
			time: time.Unix(1<<34, 0),
		},
		{
			name: "Negative Unix timestamp",
			time: time.Unix(-1, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.time)

			size, err := encoder.CalcByteSize(value)
			tu.NoError(t, err)

			bytes := make([]byte, size)
			offset := encoder.WriteToBytes(value, 0, &bytes)
			tu.Equal(t, offset, size)

			// Verify basic structure is valid
			if len(bytes) < 2 {
				t.Error("Byte slice too short")
			}
		})
	}
}

func TestWriteToBytesWithVariousNanoseconds(t *testing.T) {
	encoder := Encoder

	tests := []struct {
		name        string
		time        time.Time
		expectFmt   byte
		description string
	}{
		{
			name:        "Zero nanoseconds",
			time:        time.Unix(1000, 0),
			expectFmt:   def.Fixext4,
			description: "Should use Fixext4 when nanoseconds is 0",
		},
		{
			name:        "Small nanoseconds (1)",
			time:        time.Unix(1000, 1),
			expectFmt:   def.Fixext8,
			description: "Should use Fixext8 with nanoseconds",
		},
		{
			name:        "Mid nanoseconds",
			time:        time.Unix(1000, 500000000),
			expectFmt:   def.Fixext8,
			description: "Should use Fixext8 with mid-range nanoseconds",
		},
		{
			name:        "Max nanoseconds",
			time:        time.Unix(1000, 999999999),
			expectFmt:   def.Fixext8,
			description: "Should use Fixext8 with max nanoseconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.time)
			size, err := encoder.CalcByteSize(value)
			tu.NoError(t, err)

			bytes := make([]byte, size)
			encoder.WriteToBytes(value, 0, &bytes)

			tu.Equal(t, bytes[0], tt.expectFmt)
		})
	}
}
