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

func TestHookWhenSomeoneIs_PairingOnTheTrunk_WithANewPair(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{tam}
		options       = newOptions().WorkingOnTrunkWithProtectionSetTo(true).Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)

	_, err := runHook(t, options, []string{"Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_PairingOnTheTrunk_WithTheSamePersonAsLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{pete}
		options       = newOptions().WorkingOnTrunkWithProtectionSetTo(true).Build()
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

func TestHookWhenSomeoneIs_PairingOnTheTrunk_WithDifferentPeopleThanLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		previousPairs = lib.CoAuthors{pete}
		expectedPairs = lib.CoAuthors{tam}
		options       = newOptions().WorkingOnTrunkWithProtectionSetTo(true).Build()
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

func TestHookWhenSomeoneIs_WorkingAlone_OnTheTrunk_AndBranchProtectionIsEnabled(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs lib.CoAuthors
		options       = newOptions().WorkingOnTrunkWithProtectionSetTo(true).Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)

	output, err := runHook(t, options, []string{"No one else"})
	assert.Error(t, err)
	assert.Contains(t, output, fmt.Sprintf("can't commit to %s without a pair", options.TrunkName))
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_WorkingAlone_OnTheTrunk_AndBranchProtectionIsDisabled(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs lib.CoAuthors
		options       = newOptions().WorkingOnTrunkWithProtectionSetTo(false).Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)

	_, err := runHook(t, options, []string{"No one else"})
	assert.NoError(t, err)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func TestHookWhenSomeoneIs_WorkingAlone_OnABranch(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs lib.CoAuthors
		options       = newOptions().WorkingOnBranchWithProtectionSetTo(true).Build()
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)

	_, err := runHook(t, options, []string{"No one else"})
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
	bytes, err := json.Marshal(pairs)
	assert.NoError(t, err, "could not marshall pairs")

	err = os.WriteFile(pairsFilePath, bytes, 0666)
	assert.NoError(t, err, "could not write pairs file")
}

func assertPairsFileHasEqualPairs(t *testing.T, expectedPairs lib.CoAuthors) {
	t.Helper()
	b, err := os.ReadFile(pairsFilePath)
	assert.NoError(t, err, "could not read file %q", pairsFilePath)

	var actualPairs lib.CoAuthors
	assert.NoError(t, json.Unmarshal(b, &actualPairs))
	assert.Equal(t, expectedPairs, actualPairs)
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

type inputType interface {
	string | []byte
}

func runHook(t *testing.T, options src.Options, textToSubmit []string) (string, error) {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../main.go",
		fmt.Sprintf("--commitFile=%s", options.CommitFilePath),
		fmt.Sprintf("--authorsFile=%s", options.AuthorsFilePath),
		fmt.Sprintf("--pairsFile=%s", options.PairsFilePath),
		fmt.Sprintf("--trunkName=%s", options.TrunkName),
		fmt.Sprintf("--branchName=%s", options.BranchName),
		fmt.Sprintf("--protectTrunk=%t", options.ProtectTrunk),
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
	options src.Options
}

func newOptions() optionsBuilder {
	return optionsBuilder{
		options: src.Options{
			CommitFilePath:     commitFilePath,
			AuthorsFilePath:    authorsFilePath,
			PairsFilePath:      pairsFilePath,
			TrunkName:          "trunk",
			BranchName:         "trunk",
			ProtectTrunk:       true,
			ForceSearchPrompts: true,
		},
	}
}

func (b optionsBuilder) Build() src.Options {
	return b.options
}

func (b optionsBuilder) WorkingOnTrunkWithProtectionSetTo(protect bool) optionsBuilder {
	b.options.TrunkName = "trunk"
	b.options.BranchName = "trunk"
	b.options.ProtectTrunk = protect
	return b
}

func (b optionsBuilder) WorkingOnBranchWithProtectionSetTo(protect bool) optionsBuilder {
	b.options.TrunkName = "trunk"
	b.options.BranchName = "not-trunk"
	b.options.ProtectTrunk = protect
	return b
}
