package blackboxtests

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
)

func TestMain(m *testing.M) {
	cleanup()
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func givenThereIsACommitMessageFile(t *testing.T, message string) {
	t.Helper()
	err := os.WriteFile(commitFilePath, []byte(message), 0o666)
	assert.NoError(t, err)
}

func givenThereIsAnAuthorsFile(t *testing.T, authors lib.CoAuthors) {
	t.Helper()
	bytes := []byte(authors.String())

	err := os.WriteFile(authorsFilePath, bytes, 0o666)
	assert.NoError(t, err, "could not write authors file")
}

func givenThereIsAPairsFile(t *testing.T, pairs []string) {
	t.Helper()
	if len(pairs) == 0 {
		pairs = []string{}
	}

	b, err := json.Marshal(pairs)
	assert.NoError(t, err, "could not marshall pairs")

	err = os.WriteFile(pairsFilePath, b, 0o666)
	assert.NoError(t, err, "could not write pairs file")
}

func givenThereIsNotAPairsFile() {
	_ = os.Remove(pairsFilePath)
}

func assertPairsFileHasEqualPairs(t *testing.T, expectedPairs lib.CoAuthors) {
	t.Helper()
	b, err := os.ReadFile(pairsFilePath)
	assert.NoError(t, err, "could not read file %q", pairsFilePath)

	var actualPairNames []string
	assert.NoError(t, json.Unmarshal(b, &actualPairNames))
	assert.Equal(t, expectedPairs.Names(), actualPairNames)
}

func assertCommitMessageFileHasContents(t *testing.T, message string) {
	t.Helper()
	fileContent, err := os.ReadFile(commitFilePath)
	assert.NoError(t, err, "could not read commit file %q", commitFilePath)
	assert.Equal(t, message, string(fileContent))
}

func assertCommitMessageFileContainsContents(t *testing.T, message string) {
	t.Helper()
	fileContent, err := os.ReadFile(commitFilePath)
	assert.NoError(t, err, "could not read commit file %q", commitFilePath)
	assert.Contains(t, string(fileContent), message)
}

var (
	tam  = lib.CoAuthor{Name: "Tam", Email: "t@am.com"}
	pete = lib.CoAuthor{Name: "Pete", Email: "p@ete.com"}
)

const (
	commitFilePath  = "test_commit_file"
	authorsFilePath = "test_authors.json"
	pairsFilePath   = "test_pairs.json"
)

func cleanup() {
	_ = os.Remove(commitFilePath)
	_ = os.Remove(authorsFilePath)
	_ = os.Remove(pairsFilePath)
}
