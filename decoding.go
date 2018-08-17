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

	if string(v[len(v)-1]) != endChar {
		return 0, 0, fmt.Errorf("malformed integer")
	}

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

func decodeList(v string) ([]interface{}, int, error) {
	buff := []interface{}{}
	length := len(v)

	pos := 1
	bytesRead := 1

	if len(v) == 1 {
		return buff, 0, fmt.Errorf("empty list")
	}

	for {
		if pos == length && string(v[pos-1]) != endChar {
			return nil, 0, fmt.Errorf("malformed list")
		}

		if pos >= length {
			break
		}

		switch string(v[pos]) {
		case "e":
			bytesRead++
			return buff, bytesRead, nil
		case "i":
			str, btr, err := decodeInt(v[pos:])
			if err != nil {
				return nil, bytesRead, err
			}
			buff = append(buff, str)
			bytesRead += btr

			pos = pos + btr
		default:
			str, btr, err := decodeBytes(v[pos:])
			if err != nil {
				return nil, bytesRead, err
			}
			buff = append(buff, string(str))
			bytesRead += btr

			pos = pos + btr
		}
	}

	return buff, bytesRead, nil
}

func decodeDict(v string) (map[string]interface{}, int, error) {
	buff := make(map[string]interface{})
	length := len(v)

	// pos and bytesRead start as 1 because it skips the first char (d)
	pos := 1
	bytesRead := 1

	if length == 1 {
		return nil, 0, fmt.Errorf("malformed dictionary")
	}

	for {
		if pos >= length {
			return nil, 0, fmt.Errorf("malformed dictionary")
		}

		if string(v[pos]) == endChar {
			bytesRead++
			return buff, bytesRead, nil
		}

		key, btr, err := decodeBytes(v[pos:])
		if err != nil {
			return nil, bytesRead, err
		}

		pos = pos + btr
		bytesRead += btr

		switch string(v[pos]) {
		case "i":
			str, btr, err := decodeInt(v[pos:])
			if err != nil {
				return nil, 0, err
			}

			pos = pos + btr
			bytesRead += btr

			buff[string(key)] = fmt.Sprintf("%v", str)
		case "l":
			list, btr, err := decodeList(v[pos : length-1])
			if err != nil {
				return nil, 0, err
			}

			pos = pos + btr
			bytesRead += btr

			buff[string(key)] = list
		case "d":
			dic, btr, err := decodeDict(v[pos : len(v)-1])
			if err != nil {
				return nil, 0, err
			}

			pos = pos + btr
			bytesRead += btr

			buff[string(key)] = dic
		default:
			str, btr, err := decodeBytes(v[pos:])
			if err != nil {
				return nil, 0, err
			}

			pos = pos + btr
			bytesRead += btr

			buff[string(key)] = string(str)
		}
	}

	return buff, bytesRead, nil
}
