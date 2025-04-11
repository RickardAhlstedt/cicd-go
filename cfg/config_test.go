package cfg

import (
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	config := &BuildConfig{
		Ignore: []string{"cfg/**", "test/**", "*.log"},
	}

	tests := []struct {
		path     string
		expected bool
	}{
		{"cfg/", true},
		{"cfg/file.txt", true},
		{"cfg/subdir/file.txt", true},
		{"test/file.txt", true},
		{"test/subdir/file.txt", true},
		{"logfile.log", true},
		{"src/main.go", false},
		{"testfile.txt", false},
	}

	for _, test := range tests {
		result := config.ShouldIgnore(test.path)
		if result != test.expected {
			t.Errorf("ShouldIgnore(%q) = %v; want %v", test.path, result, test.expected)
		}
	}
}
