package time

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	"github.com/shamaton/msgpack/v2/internal/common"
	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func TestStreamCode(t *testing.T) {
	encoder := StreamEncoder
	code := encoder.Code()
	tu.Equal(t, code, def.TimeStamp)
}

func TestStreamType(t *testing.T) {
	encoder := StreamEncoder
	typ := encoder.Type()
	expected := reflect.TypeOf(time.Time{})
	tu.Equal(t, typ, expected)
}

func TestStreamWrite(t *testing.T) {
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

	encoder := StreamEncoder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.time)
			buf := &bytes.Buffer{}
			buffer := common.GetBuffer()
			defer common.PutBuffer(buffer)
			w := ext.CreateStreamWriter(buf, buffer)

			err := encoder.Write(w, value)
			tu.NoError(t, err)

			// Flush buffer to writer
			err = buffer.Flush(buf)
			tu.NoError(t, err)

			b := buf.Bytes()
			tu.Equal(t, len(b), tt.expectedLen)

			// Verify format type
			switch tt.expectedFormat {
			case "fixext4":
				tu.Equal(t, b[0], def.Fixext4)
				tu.Equal(t, int8(b[1]), def.TimeStamp)

			case "fixext8":
				tu.Equal(t, b[0], def.Fixext8)
				tu.Equal(t, int8(b[1]), def.TimeStamp)

			case "ext8":
				tu.Equal(t, b[0], def.Ext8)
				tu.Equal(t, b[1], 12)
				tu.Equal(t, int8(b[2]), def.TimeStamp)

			default:
				t.Errorf("Unknown expected format: %s", tt.expectedFormat)
			}
		})
	}
}

func TestStreamWriteEdgeCases(t *testing.T) {
	encoder := StreamEncoder

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
			buf := &bytes.Buffer{}
			buffer := common.GetBuffer()
			defer common.PutBuffer(buffer)
			w := ext.CreateStreamWriter(buf, buffer)

			err := encoder.Write(w, value)
			tu.NoError(t, err)

			// Flush buffer to writer
			err = buffer.Flush(buf)
			tu.NoError(t, err)

			b := buf.Bytes()
			// Verify basic structure is valid
			if len(b) < 2 {
				t.Error("Byte slice too short")
			}
		})
	}
}

func TestStreamEncodedDataAccuracy(t *testing.T) {
	encoder := StreamEncoder

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
			buf := &bytes.Buffer{}
			buffer := common.GetBuffer()
			defer common.PutBuffer(buffer)
			w := ext.CreateStreamWriter(buf, buffer)

			err := encoder.Write(w, value)
			tu.NoError(t, err)

			// Flush buffer to writer
			err = buffer.Flush(buf)
			tu.NoError(t, err)

			b := buf.Bytes()

			// Verify we can extract the correct time back
			switch b[0] {
			case def.Fixext4:
				data := binary.BigEndian.Uint32(b[2:6])
				secs := int64(data)
				tu.Equal(t, secs, tt.time.Unix())

			case def.Fixext8:
				data := binary.BigEndian.Uint64(b[2:10])
				nano := int64(data >> 34)
				secs := int64(data & 0x00000003ffffffff)
				tu.Equal(t, secs, tt.time.Unix())
				tu.Equal(t, nano, int64(tt.time.Nanosecond()))

			case def.Ext8:
				nano := binary.BigEndian.Uint32(b[3:7])
				secs := binary.BigEndian.Uint64(b[7:15])
				tu.Equal(t, int64(secs), tt.time.Unix())
				tu.Equal(t, int64(nano), int64(tt.time.Nanosecond()))
			}
		})
	}
}

func TestStreamWriteWithVariousNanoseconds(t *testing.T) {
	encoder := StreamEncoder

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
			buf := &bytes.Buffer{}
			buffer := common.GetBuffer()
			defer common.PutBuffer(buffer)
			w := ext.CreateStreamWriter(buf, buffer)

			err := encoder.Write(w, value)
			tu.NoError(t, err)

			// Flush buffer to writer
			err = buffer.Flush(buf)
			tu.NoError(t, err)

			b := buf.Bytes()
			tu.Equal(t, b[0], tt.expectFmt)
		})
	}
}

type testErrWriter struct {
	ErrorBytes []byte
	Count      int
}

func (w *testErrWriter) Write(p []byte) (n int, err error) {
	if bytes.Equal(w.ErrorBytes, p) {
		return 0, errors.New("equal bytes error")
	}
	return len(p), nil
}

func TestStreamWriteErrors(t *testing.T) {
	encoder := StreamEncoder

	ts := def.TimeStamp

	tests := []struct {
		name         string
		timeValue    time.Time
		errorBytes   []byte
		prepareSize  int
		prepareBytes []byte
	}{
		// Fixext4 tests
		{
			name:         "Fixext4 - Error on writing Fixext4 type byte",
			timeValue:    time.Unix(1000, 0),
			errorBytes:   []byte{255},
			prepareSize:  1,
			prepareBytes: []byte{255},
		},
		{
			name:        "Fixext4 - Error on writing TimeStamp type byte",
			timeValue:   time.Unix(1000, 0),
			errorBytes:  []byte{def.Fixext4},
			prepareSize: 1,
		},
		{
			name:        "Fixext4 - Error on writing 4 bytes of data",
			timeValue:   time.Unix(1000, 0),
			errorBytes:  []byte{def.Fixext4, byte(ts)},
			prepareSize: 2,
		},
		// Fixext8 tests
		{
			name:         "Fixext8 - Error on writing Fixext8 type byte",
			timeValue:    time.Unix(1000, 1),
			errorBytes:   []byte{255},
			prepareSize:  1,
			prepareBytes: []byte{255},
		},
		{
			name:        "Fixext8 - Error on writing TimeStamp type byte",
			timeValue:   time.Unix(1000, 1),
			errorBytes:  []byte{def.Fixext8},
			prepareSize: 1,
		},
		{
			name:        "Fixext8 - Error on writing 8 bytes of data",
			timeValue:   time.Unix(1000, 1),
			errorBytes:  []byte{def.Fixext8, byte(ts)},
			prepareSize: 2,
		},
		// Ext8 tests
		{
			name:         "Ext8 - Error on writing Ext8 type byte",
			timeValue:    time.Unix(1<<34, 0),
			errorBytes:   []byte{255},
			prepareSize:  1,
			prepareBytes: []byte{255},
		},
		{
			name:        "Ext8 - Error on writing length byte",
			timeValue:   time.Unix(1<<34, 0),
			errorBytes:  []byte{def.Ext8},
			prepareSize: 1,
		},
		{
			name:        "Ext8 - Error on writing TimeStamp type byte",
			timeValue:   time.Unix(1<<34, 0),
			errorBytes:  []byte{def.Ext8, 12},
			prepareSize: 2,
		},
		{
			name:        "Ext8 - Error on writing 4 bytes of nanoseconds",
			timeValue:   time.Unix(1<<34, 0),
			errorBytes:  []byte{def.Ext8, 12, byte(ts)},
			prepareSize: 3,
		},
		{
			name:        "Ext8 - Error on writing 8 bytes of seconds",
			timeValue:   time.Unix(1<<34, 123456789),
			errorBytes:  []byte{def.Ext8, 12, byte(ts), 0x07, 0x5b, 0xcd, 0x15},
			prepareSize: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &common.Buffer{Data: make([]byte, tt.prepareSize)}
			err := buffer.Write(nil, tt.prepareBytes...)
			tu.NoError(t, err)

			errWriter := &testErrWriter{ErrorBytes: tt.errorBytes}
			w := ext.CreateStreamWriter(errWriter, buffer)
			value := reflect.ValueOf(tt.timeValue)
			err = encoder.Write(w, value)
			tu.Error(t, err)
		})
	}
}
