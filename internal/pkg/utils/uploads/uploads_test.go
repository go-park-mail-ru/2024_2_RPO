package uploads

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestJoinFilePath(t *testing.T) {
	assert.Equal(t, JoinFilePath("1234", "jpg"), "1234.jpg")

}
