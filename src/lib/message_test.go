package lib_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"testing"
)

var (
	tam    = lib.CoAuthor{Name: "tam", Email: "t@am.com"}
	john   = lib.CoAuthor{Name: "John Doe", Email: "john@doe.com"}
	mary   = lib.CoAuthor{Name: "Mary Sue", Email: "m.sue@example.com"}
	rizzle = lib.CoAuthor{Name: "rizzle", Email: "rizzle@kicks.co"}
)

func TestAddingCoAuthorsToPlainMessage(t *testing.T) {
	commitMessage := "Hello world :D"
	coAuthors := []lib.CoAuthor{tam, john}

	expectedMessage := fmt.Sprintf("Hello world :D\n\n%s\n%s", tam, john)

	preparedMessage := lib.PrepareCommitMessage(commitMessage, coAuthors)
	assert.Equal(t, expectedMessage, preparedMessage)
}

func TestDoesNotAddCoauthorThatAlreadyExists(t *testing.T) {
	commitMessage := "Hello world :D\n\n" + john.String()
	coAuthors := []lib.CoAuthor{tam, john, rizzle}

	expectedMessage := fmt.Sprintf("%s\n\n%s\n%s", commitMessage, tam, rizzle)

	preparedMessage := lib.PrepareCommitMessage(commitMessage, coAuthors)
	assert.Equal(t, expectedMessage, preparedMessage)
}

func TestAddingCoAuthorsToTemplatedMessage(t *testing.T) {
	inputMessage := "Hello world :D" + lib.COMMIT_SEPARATOR + "\nother stuff"
	coAuthors := []lib.CoAuthor{tam, john}

	expectedMessage := fmt.Sprintf("Hello world :D\n\n%s\n%s%s\nother stuff", tam, john, lib.COMMIT_SEPARATOR)

	actualMessage := lib.PrepareCommitMessage(inputMessage, coAuthors)
	assert.Equal(t, expectedMessage, actualMessage)
}
