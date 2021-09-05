package service

import (
	"encoding/base64"
	"math/rand"
	"strings"
)

func caesarShift(text string, shift int) (string, error) {
	var upLow string

	for _, char := range text {
		if string(char) == strings.ToUpper(string(char)) {
			upLow += strings.ToLower(string(char))
		} else if string(char) == strings.ToLower(string(char)) {
			upLow += strings.ToUpper(string(char))
		} else {
			upLow += string(char)
		}
	}

	decStr, err := base64.StdEncoding.DecodeString(reverseString(upLow))
	if err != nil {
		return "", err
	}

	text = string(decStr)[1 : len(string(decStr))-1]

	shift = (shift%26 + 26) % 26 // [0, 25]

	b := make([]byte, len(text))

	for i := 0; i < len(text); i++ {
		t := text[i]

		var a int

		switch {
		case 'a' <= t && t <= 'z':
			a = 'a'
		case 'A' <= t && t <= 'Z':
			a = 'A'
		default:
			b[i] = t
			continue
		}
		b[i] = byte(a + ((int(t)-a)+shift)%26)
	}

	return string(b), nil
}

func reverseString(str string) string {
	byteStr := []rune(str)
	for i, j := 0, len(byteStr)-1; i < j; i, j = i+1, j-1 {
		byteStr[i], byteStr[j] = byteStr[j], byteStr[i]
	}
	return string(byteStr)
}

func emailCode() string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
