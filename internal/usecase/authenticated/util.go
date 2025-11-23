package authenticated

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// GenerateElementID generates the element ID for the given method and path
func GenerateElementID(method, path string) string {
	return fmt.Sprintf("%s|%s", method, path)
}

// ParseElementID parses the method and path from the element ID
func ParseElementID(elementID string) (method, path string) {
	parts := strings.Split(elementID, "|")
	return parts[0], parts[1]
}

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	const otpChars = "1234567890"
	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
