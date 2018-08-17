package bencode

import (
	"bytes"
	"fmt"
	"strconv"
)

var endChar = "e"
var stringLimiter = ":"

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

	length, prefixDigits, err := parseStringLength(v)
	if err != nil {
		return nil, 0, err
	}

	if length > len(v) {
		return nil, 0, fmt.Errorf("empty string")
	}

	if length > len(v[prefixDigits:]) {
		return nil, 0, fmt.Errorf("invalid string length")
	}

	for _, b := range v[prefixDigits : length+prefixDigits] {
		buff.WriteRune(b)
	}

	bt := buff.Bytes()

	bytesRead := buff.Len() + prefixDigits

	return bt, bytesRead, nil
}

func parseStringLength(v string) (int, int, error) {
	lenBuff := bytes.Buffer{}

	for i := 0; i < len(v); i++ {
		if string(v[i]) == stringLimiter {
			break
		}

		lenBuff.WriteByte(v[i])
	}

	digits := lenBuff.Len()

	length, err := strconv.Atoi(lenBuff.String())
	if err != nil {
		return 0, 0, err
	}

	return length, digits + 1, nil
}
