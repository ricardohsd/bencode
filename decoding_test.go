package bencode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	testCases := []struct {
		input          string
		byteSize       int
		expectedNumber int64
		err            error
	}{
		{
			input:          "ie",
			byteSize:       0,
			expectedNumber: 0,
			err:            fmt.Errorf("empty integer"),
		},
		{
			input:          "iaaae",
			byteSize:       0,
			expectedNumber: 0,
			err:            fmt.Errorf("not an integer"),
		},
		{
			input:          "i59616",
			byteSize:       0,
			expectedNumber: 0,
			err:            fmt.Errorf("malformed integer"),
		},
		{
			input:          "i0e",
			byteSize:       3,
			expectedNumber: 0,
			err:            nil,
		},
		{
			input:          "i59616e",
			byteSize:       7,
			expectedNumber: 59616,
			err:            nil,
		},
		{
			input:          "i-59616e",
			byteSize:       8,
			expectedNumber: -59616,
			err:            nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, bytesRead, err := decodeInt(tc.input)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			}

			assert.Equal(t, tc.byteSize, bytesRead)
			assert.Equal(t, tc.expectedNumber, result)
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		input    string
		byteSize int
		expected string
		err      error
	}{
		{
			input:    "7:",
			byteSize: 0,
			expected: "",
			err:      fmt.Errorf("empty string"),
		},
		{
			input:    "8",
			byteSize: 0,
			expected: "",
			err:      fmt.Errorf("empty string"),
		},
		{
			input:    "8:johndoe",
			byteSize: 0,
			expected: "",
			err:      fmt.Errorf("invalid string length"),
		},
		{
			input:    "7:johndoe",
			byteSize: 9,
			expected: "johndoe",
			err:      nil,
		},
		{
			input:    "13:creation date",
			byteSize: 16,
			expected: "creation date",
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, bytesRead, err := decodeBytes(tc.input)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			}

			assert.Equal(t, tc.byteSize, bytesRead)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		input    string
		byteSize int
		expected []interface{}
		err      error
	}{
		{
			input:    "l",
			byteSize: 0,
			expected: []interface{}{},
			err:      fmt.Errorf("empty list"),
		},
		{
			input:    "l5",
			byteSize: 1,
			expected: []interface{}(nil),
			err:      fmt.Errorf("empty string"),
		},
		{
			input:    "l5:ItemA5:ItemBe",
			byteSize: 16,
			expected: []interface{}{"ItemA", "ItemB"},
		},
		{
			input:    "l4:spami42ee",
			byteSize: 12,
			expected: []interface{}{"spam", int64(42)},
		},
		{
			input:    "l5:ItemA5:ItemB",
			byteSize: 0,
			expected: []interface{}(nil),
			err:      fmt.Errorf("malformed list"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, bytesRead, err := decodeList(tc.input)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			}

			assert.Equal(t, tc.byteSize, bytesRead)
			assert.Equal(t, tc.expected, result)
		})
	}
}
