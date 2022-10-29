package utility

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	buffer := make([]rune, length)
	for i := range buffer {
		buffer[i] = letters[rand.Intn(len(letters))]
	}

	return string(buffer)
}
