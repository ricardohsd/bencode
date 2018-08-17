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
