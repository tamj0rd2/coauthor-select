package blackboxtests

import (
	"fmt"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
)

func Test_InteractiveSelectHook_WhenSomeoneIs_WorkingAlone(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)

	err := runInteractiveSelectHook(t, []string{"No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func Test_InteractiveSelectHook_WhenSomeoneIs_Pairing_ForTheFirstTime_WithASinglePerson(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{tam}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsNotAPairsFile()

	err := runInteractiveSelectHook(t, []string{"Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func Test_InteractiveSelectHook_WhenSomeoneIs_Pairing_ForTheFirstTime_WithMultiplePeople(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{tam, pete}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsNotAPairsFile()

	err := runInteractiveSelectHook(t, []string{"Tam", "Pete", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func Test_InteractiveSelectHook_WhenSomeoneIs_Pairing_WithTheSamePersonAsLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		expectedPairs = lib.CoAuthors{pete}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, expectedPairs.Names())

	err := runInteractiveSelectHook(t, []string{"Yes"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func Test_InteractiveSelectHook_WhenSomeoneIs_Pairing_WithDifferentPeopleThanLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		previousPairs = lib.CoAuthors{pete}
		expectedPairs = lib.CoAuthors{tam}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, previousPairs.Names())

	err := runInteractiveSelectHook(t, []string{"No", "Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func Test_InteractiveSelectHook_WhenSomeoneIs_Pairing_ButWasWorkingAloneLastTime(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		previousPairs = lib.CoAuthors{}
		expectedPairs = lib.CoAuthors{tam}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, previousPairs.Names())

	err := runInteractiveSelectHook(t, []string{"Tam", "No one else"})
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, expectedPairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, expectedPairs)
}

func runInteractiveSelectHook(t *testing.T, textToSubmit []string) error {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../...", "select",
		fmt.Sprintf("--commitFile=%s", commitFilePath),
		fmt.Sprintf("--authorsFile=%s", authorsFilePath),
		fmt.Sprintf("--pairsFile=%s", pairsFilePath),
		fmt.Sprintf("--forceSearchPrompts=%t", true),
		fmt.Sprintf("--interactive=%t", true),
	)

	cmdStdin, err := cmd.StdinPipe()
	assert.NoError(t, err)

	go func() {
		defer func() {
			_ = cmdStdin.Close()
		}()

		maxIndex := len(textToSubmit) - 1
		for i, text := range textToSubmit {
			if _, err := io.WriteString(cmdStdin, text+"\n"); err != nil {
				panic(fmt.Errorf("failed to write %q to stdin: %v\n", text, err))
			}

			if i < maxIndex {
				// the console thing promptui uses is apparently too slow to read inputs so quickly
				time.Sleep(time.Second)
			}
		}
	}()

	b, err := cmd.CombinedOutput()
	t.Log("CLI output:\n", string(b))
	return err
}
