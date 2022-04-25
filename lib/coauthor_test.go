package lib_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/lib"
	"testing"
)

func TestCoAuthor_String(t *testing.T) {
	coauthor := lib.CoAuthor{
		Name:  "John Doe",
		Email: "john@doe.com",
	}

	assert.Equal(t, "Co-authored-by: John Doe <john@doe.com>", coauthor.String())
}
