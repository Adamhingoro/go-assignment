package helper

import (
	"testing"
)

func TestHashString(t *testing.T) {
	// Test cases
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"world", "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}

	for _, test := range tests {
		result := HashString(test.input)
		if result != test.expected {
			t.Errorf("HashString(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}
