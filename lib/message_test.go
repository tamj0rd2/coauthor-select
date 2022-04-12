package lib_test

import (
	"testing"

	"github.com/tamj0rd2/coauthor-select/lib"
)

func TestAddingCoAuthorsToPlainMessage(t *testing.T) {
	inputMessage := "Hello world :D"
	coAuthors := []lib.CoAuthor{{Name: "tam", Email: "t@am.com"}, {Name: "tam2", Email: "t@am2.com"}}

	expectedMessage := "Hello world :D\n\nCo-authored-by: tam <t@am.com>\nCo-authored-by: tam2 <t@am2.com>"

	actualMessage := lib.PrepareCommitMessage(inputMessage, coAuthors)

	if actualMessage != expectedMessage {
		t.Fatalf("\nEXPECTED:\n%v\n\nGOT:\n%v", expectedMessage, actualMessage)
	}
}

func TestAddingCoAuthorsToTemplatedMessage(t *testing.T) {
	inputMessage := "Hello world :D" + lib.COMMIT_SEPARATOR + "\nother stuff"
	coAuthors := []lib.CoAuthor{{Name: "tam", Email: "t@am.com"}}

	expectedMessage := "Hello world :D\n\nCo-authored-by: tam <t@am.com>" + lib.COMMIT_SEPARATOR + "\nother stuff"

	actualMessage := lib.PrepareCommitMessage(inputMessage, coAuthors)

	if actualMessage != expectedMessage {
		t.Fatalf("\nEXPECTED:\n%v\n\nGOT:\n%v", expectedMessage, actualMessage)
	}
}
