package uploads

import (
	"os"
	"testing"
)

func TestJoinFileURL(t *testing.T) {
	os.Setenv("USER_UPLOADS_URL", "http://example.com/uploads/")

	tests := []struct {
		fileUUID      string
		fileExtension string
		defaultValue  string
		expected      string
	}{
		{"", "", DefaultAvatarURL, DefaultAvatarURL},
		{"12345", "", DefaultAvatarURL, "http://example.com/uploads/12345"},
		{"12345", "jpg", DefaultAvatarURL, "http://example.com/uploads/12345.jpg"},
	}

	for _, test := range tests {
		result := JoinFileURL(test.fileUUID, test.fileExtension, test.defaultValue)
		if result != test.expected {
			t.Errorf("JoinFileURL(%q, %q, %q) = %q; want %q", test.fileUUID, test.fileExtension, test.defaultValue, result, test.expected)
		}
	}
}

func TestExtractFileExtension(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"file.jpg", "jpg"},
		{"file.tar.gz", "gz"},
		{"file", ""},
		{"file.", ""},
	}

	for _, test := range tests {
		result := ExtractFileExtension(test.fileName)
		if result != test.expected {
			t.Errorf("ExtractFileExtension(%q) = %q; want %q", test.fileName, result, test.expected)
		}
	}
}

func TestJoinFilePath(t *testing.T) {
	tests := []struct {
		fileUUID      string
		fileExtension string
		expected      string
	}{
		{"12345", "", "12345"},
		{"12345", "jpg", "12345.jpg"},
	}

	for _, test := range tests {
		result := JoinFilePath(test.fileUUID, test.fileExtension)
		if result != test.expected {
			t.Errorf("JoinFilePath(%q, %q) = %q; want %q", test.fileUUID, test.fileExtension, result, test.expected)
		}
	}
}
