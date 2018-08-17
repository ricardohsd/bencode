package bencode

import (
	"bytes"
	"fmt"
	"strconv"
)

var endChar = "e"

func decodeInt(v string) (int64, int, error) {
	buff := bytes.Buffer{}

	for _, b := range v[1:] {
		if string(b) == endChar {
			break
		}

		buff.WriteRune(b)
	}

	if buff.Len() == 0 {
		return 0, 0, fmt.Errorf("empty integer")
	}

	str := buff.String()

	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("not an integer")
	}

	bytesRead := len(str) + 2

	return n, bytesRead, nil
}

func decodeBytes(v string) ([]byte, int, error) {
	buff := bytes.Buffer{}

	length, err := strconv.Atoi(string(v[0]))
	if err != nil {
		return nil, 0, err
	}

	if length > len(v) {
		return nil, 0, fmt.Errorf("empty string")
	}

	if length > len(v[2:]) {
		return nil, 0, fmt.Errorf("invalid string length")
	}

	for _, b := range v[2 : length+2] {
		buff.WriteRune(b)
	}

	bt := buff.Bytes()

	bytesRead := buff.Len() + 2

	return bt, bytesRead, nil
}
