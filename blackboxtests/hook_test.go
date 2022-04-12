package blackboxtests

import (
	"encoding/json"
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/lib"
	"os"
	"os/exec"
	"testing"
)

func TestHook(t *testing.T) {
	var (
		commitFilePath = "test_commit_file"
		commitMessage  = "feat-376 Did some work"

		authors = []lib.CoAuthor{{
			Name:  "Tam",
			Email: "tam@tam.com",
		}}
		authorsFilePath = "test_authors.json"

		pairs         = []string{"Tam"}
		pairsFilePath = "test_pairs.json"
	)
	givenThereIsACommitMessageFile(t, commitFilePath, commitMessage)
	givenThereIsAnAuthorsFile(t, authorsFilePath, authors)
	givenThereIsAPairsFile(t, pairsFilePath, pairs)

	output, err := exec.Command(
		"go", "run", "../main.go",
		"--commitFile", commitFilePath,
		"--authorsFile", authorsFilePath,
		"--pairsFile", pairsFilePath,
	).CombinedOutput()
	t.Log("CLI output:\n", string(output))
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, authors)
	assertCommitMessageFileHasContents(t, commitFilePath, expectedMessage)
}

func givenThereIsACommitMessageFile(t *testing.T, filePath string, message string) {
	t.Helper()
	err := os.WriteFile(filePath, []byte(message), 0666)
	assert.NoError(t, err)
}

func givenThereIsAnAuthorsFile(t *testing.T, filePath string, authors []lib.CoAuthor) {
	t.Helper()
	bytes, err := json.Marshal(authors)
	assert.NoError(t, err, "could not marshall authors")

	err = os.WriteFile(filePath, bytes, 0666)
	assert.NoError(t, err, "could not write authors file")
}

func givenThereIsAPairsFile(t *testing.T, filePath string, pairs []string) {
	t.Helper()
	bytes, err := json.Marshal(pairs)
	assert.NoError(t, err, "could not marshall pairs")

	err = os.WriteFile(filePath, bytes, 0666)
	assert.NoError(t, err, "could not write pairs file")
}

func assertCommitMessageFileHasContents(t *testing.T, filePath string, message string) {
	t.Helper()
	fileContent, err := os.ReadFile(filePath)
	assert.NoError(t, err, "could not read commit file", filePath)
	assert.Equal(t, message, string(fileContent))
}
