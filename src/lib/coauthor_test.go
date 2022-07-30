package lib_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"testing"
)

func TestCoAuthor_String(t *testing.T) {
	assert.Equal(t, "Co-authored-by: John Doe <john@doe.com>", john.String())
}

func TestCoAuthor_UserID(t *testing.T) {
	assert.Equal(t, "John Doe <john@doe.com>", john.UserID())
}

func TestCoAuthors_String(t *testing.T) {
	coAuthors := lib.CoAuthors{john, mary}
	assert.Equal(t, "John Doe <john@doe.com>\nMary Sue <m.sue@example.com>", coAuthors.String())
}

func TestCoAuthors_From(t *testing.T) {
	input := "John Doe <john@doe.com>\nMary Sue <m.sue@example.com>"

	var coAuthors lib.CoAuthors
	assert.NoError(t, coAuthors.From([]byte(input)))

	assert.Equal(t, lib.CoAuthors{john, mary}, coAuthors)
}
