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
