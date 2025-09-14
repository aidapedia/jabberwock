package authenticated

import (
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
