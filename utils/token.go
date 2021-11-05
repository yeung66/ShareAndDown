package utils

import (
	"crypto/rand"
	"fmt"
)

func TokenGenerator() string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
