package git_test

import (
	"testing"

	"github.com/tamj0rd2/coauthor-select/git"
)

func TestAddingCoAuthorsToPlainMessage(t *testing.T) {
	inputMessage := "Hello world :D"
	coAuthors := []git.CoAuthor{{Name: "tam", Email: "t@am.com"}, {Name: "tam2", Email: "t@am2.com"}}

	expectedMessage := "Hello world :D\n\nCo-authored-by: tam <t@am.com>\nCo-authored-by: tam2 <t@am2.com>"

	actualMessage := git.PrepareCommitMessage(inputMessage, coAuthors)

	if actualMessage != expectedMessage {
		t.Fatalf("\nEXPECTED:\n%v\n\nGOT:\n%v", expectedMessage, actualMessage)
	}
}

func TestAddingCoAuthorsToTemplatedMessage(t *testing.T) {
	inputMessage := "Hello world :D" + git.COMMIT_SEPARATOR + "\nother stuff"
	coAuthors := []git.CoAuthor{{Name: "tam", Email: "t@am.com"}}

	expectedMessage := "Hello world :D\n\nCo-authored-by: tam <t@am.com>" + git.COMMIT_SEPARATOR + "\nother stuff"

	actualMessage := git.PrepareCommitMessage(inputMessage, coAuthors)

	if actualMessage != expectedMessage {
		t.Fatalf("\nEXPECTED:\n%v\n\nGOT:\n%v", expectedMessage, actualMessage)
	}
}
