package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789")

	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = chars[rnd.Intn(len(chars))]
	}

	return string(randomString)
}
