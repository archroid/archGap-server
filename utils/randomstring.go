package utils

import "crypto/rand"

// GenerateRandomString creates a random string of fixed length
func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "rnd"
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
