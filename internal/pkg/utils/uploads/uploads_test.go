package uploads

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestJoinFileName(t *testing.T) {
	assert.Equal(t, JoinFileName("1234", "jpg", "default"), "1234.jpg")

}
