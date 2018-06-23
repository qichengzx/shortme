package utils

import (
	"bytes"
	"math"
)

const codeLen = 62
const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Encode(number int64) string {
	if number == 0 {
		return string(charset[0])
	}

	chars := make([]byte, 0)
	for number > 0 {
		remainder := number % codeLen
		chars = append(chars, charset[remainder])
		number = number / codeLen
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

func Decode(token string) int {
	number := 0
	chars := []byte(charset)
	tokenLength := len(token)

	for idx, c := range []byte(token) {
		power := tokenLength - (idx + 1)
		index := bytes.IndexByte(chars, c)
		number += index * int(math.Pow(codeLen, float64(power)))
	}

	return number
}
