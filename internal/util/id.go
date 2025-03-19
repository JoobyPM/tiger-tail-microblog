package util

import (
	"github.com/google/uuid"
)

// GeneratePostID generates a unique post ID using UUID v4
// This ensures uniqueness even under high concurrency or across multiple instances
func GeneratePostID() string {
	return "post_" + uuid.New().String()
}
