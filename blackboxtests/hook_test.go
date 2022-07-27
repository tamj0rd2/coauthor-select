package blackboxtests

import (
	"encoding/json"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestHookWhenSomeoneIs_WorkingAlone(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{}
		options       = newOptions().Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)

	_, err := runHook(t, options, []string{"No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_Pairing_ForTheFirstTime_WithASinglePerson(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{tam}
		options       = newOptions().Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsNotAPairsFile()

	_, err := runHook(t, options, []string{"Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_Pairing_ForTheFirstTime_WithMultiplePeople(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{tam, pete}
		options       = newOptions().Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsNotAPairsFile()

	_, err := runHook(t, options, []string{"Tam", "Pete", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_Pairing_WithTheSamePersonAsLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{pete}
		options       = newOptions().Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, expectedPairs.Names())

	_, err := runHook(t, options, []string{"Yes"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_Pairing_WithDifferentPeopleThanLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		previousPairs = lib.CoAuthors{pete}
		expectedPairs = lib.CoAuthors{tam}
		options       = newOptions().Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, previousPairs.Names())

	_, err := runHook(t, options, []string{"No", "Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_Pairing_ButWasWorkingAloneLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		previousPairs = lib.CoAuthors{}
		expectedPairs = lib.CoAuthors{tam}
		options       = newOptions().Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, previousPairs.Names())

	_, err := runHook(t, options, []string{"Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func givenThereIsACommitMessageFile(t *testing.T, message string) {
	t.Helper()
	err := os.WriteFile(commitFilePath, []byte(message), 0666)
	assert.NoError(t, err)
}

func givenThereIsAnAuthorsFile(t *testing.T, authors lib.CoAuthors) {
	t.Helper()
	bytes, err := json.Marshal(authors)
	assert.NoError(t, err, "could not marshall authors")

	err = os.WriteFile(authorsFilePath, bytes, 0666)
	assert.NoError(t, err, "could not write authors file")
}

func givenThereIsAPairsFile(t *testing.T, pairs []string) {
	t.Helper()
	if len(pairs) == 0 {
		pairs = []string{}
	}

	b, err := json.Marshal(pairs)
	assert.NoError(t, err, "could not marshall pairs")

	err = os.WriteFile(pairsFilePath, b, 0666)
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

var (
	tam  = lib.CoAuthor{Name: "Tam", Email: "t@am.com"}
	pete = lib.CoAuthor{Name: "Pete", Email: "p@ete.com"}
)

const (
	commitFilePath  = "test_commit_file"
	authorsFilePath = "test_authors.json"
	pairsFilePath   = "test_pairs.json"
)

func runHook(t *testing.T, options src.SelectOptions, textToSubmit []string) (string, error) {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../cmd/select/main.go",
		fmt.Sprintf("--commitFile=%s", options.CommitFilePath),
		fmt.Sprintf("--authorsFile=%s", options.AuthorsFilePath),
		fmt.Sprintf("--pairsFile=%s", options.PairsFilePath),
		fmt.Sprintf("--forceSearchPrompts=%t", options.ForceSearchPrompts),
	)

	cmdStdin, err := cmd.StdinPipe()
	assert.NoError(t, err)

	go func() {
		defer func() {
			assert.NoError(t, cmdStdin.Close())
		}()

		maxIndex := len(textToSubmit) - 1
		for i, text := range textToSubmit {
			if _, err := io.WriteString(cmdStdin, text+"\n"); err != nil {
				err := fmt.Errorf("failed to write %q to stdin: %v\n", text, err)
				fmt.Println(err)
				assert.NoError(t, err)
			}

			if i < maxIndex {
				// the console thing promptui uses is apparently too slow to read inputs so quickly
				time.Sleep(time.Second)
			}
		}
	}()

	b, err := cmd.CombinedOutput()
	t.Log("CLI output:\n", string(b))
	return stripansi.Strip(string(b)), err
}

func cleanup() {
	_ = os.Remove(commitFilePath)
	_ = os.Remove(authorsFilePath)
	_ = os.Remove(pairsFilePath)
}

type optionsBuilder struct {
	options src.SelectOptions
}

func newOptions() optionsBuilder {
	return optionsBuilder{
		options: src.SelectOptions{
			CommitFilePath:     commitFilePath,
			AuthorsFilePath:    authorsFilePath,
			PairsFilePath:      pairsFilePath,
			ForceSearchPrompts: true,
		},
	}
}

func (b optionsBuilder) Build() src.SelectOptions {
	return b.options
}
