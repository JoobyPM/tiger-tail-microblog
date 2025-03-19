package util

import (
	"strings"
	"testing"
)

func TestGeneratePostID(t *testing.T) {
	// Test that the ID has the correct prefix
	id := GeneratePostID()
	if !strings.HasPrefix(id, "post_") {
		t.Errorf("Expected ID to start with 'post_', got %s", id)
	}

	// Test that the ID is not empty after the prefix
	if len(id) <= 5 { // "post_" is 5 characters
		t.Errorf("Expected ID to have content after 'post_', got %s", id)
	}

	// Test that multiple calls generate different IDs
	id2 := GeneratePostID()
	if id == id2 {
		t.Errorf("Expected different IDs, got %s and %s", id, id2)
	}
}
