package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	str := "ie"
	result, bytesRead, err := decodeInt(str)
	assert.Equal(t, err.Error(), "empty integer")
	assert.Equal(t, 0, bytesRead)
	assert.Equal(t, int64(0), result)

	str = "iaaae"
	result, bytesRead, err = decodeInt(str)
	assert.Equal(t, err.Error(), "not an integer")
	assert.Equal(t, 0, bytesRead)
	assert.Equal(t, int64(0), result)

	str = "i0e"
	result, bytesRead, err = decodeInt(str)
	assert.Nil(t, err)
	assert.Equal(t, 3, bytesRead)
	assert.Equal(t, int64(0), result)

	str = "i59616e"
	result, bytesRead, err = decodeInt(str)
	assert.Nil(t, err)
	assert.Equal(t, 7, bytesRead)
	assert.Equal(t, int64(59616), result)

	str = "i-59616e"
	result, bytesRead, err = decodeInt(str)
	assert.Nil(t, err)
	assert.Equal(t, 8, bytesRead)
	assert.Equal(t, int64(-59616), result)
}
