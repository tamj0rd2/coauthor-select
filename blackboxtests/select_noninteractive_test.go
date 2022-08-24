package blackboxtests

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
)

func Test_NonInteractiveSelectHook_WhenSomeoneIs_WorkingAlone(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		pairs         = lib.CoAuthors{}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, pairs.Names())

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, pairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, pairs)
}

func Test_NonInteractiveSelectHook_WhenSomeoneIs_WorkingAlone_AndThereIsNoPairsFile(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		pairs         = lib.CoAuthors{}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsNotAPairsFile()

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, pairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, pairs)
}

func Test_NonInteractiveSelectHook_WhenSomeoneIs_Pairing_WithASinglePerson(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		pairs         = lib.CoAuthors{tam}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, pairs.Names())

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, pairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, pairs)
}

func Test_NonInteractiveSelectHook_WhenSomeoneIs_Pairing_WithMultiplePeople(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		authors       = lib.CoAuthors{tam, pete}
		pairs         = lib.CoAuthors{tam, pete}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, pairs.Names())

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, pairs)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertPairsFileHasEqualPairs(t, pairs)
}

func Test_NonInteractiveSelectHook_WhenSomeoneIs_Pairing_WithSomeoneWhoIsAlreadyListedAsACoAuthor(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work\n" + pete.String()
		authors       = lib.CoAuthors{tam, pete}
		pairs         = lib.CoAuthors{tam}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, pairs.Names())

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)
	assertCommitMessageFileContainsContents(t, pete.String())
	assertCommitMessageFileContainsContents(t, tam.String())
}

func Test_NonInteractiveSelectHook_WhenSomeoneIs_Pairing_WithSomeoneWhoIsUnspecifiedInTheAuthorsFile(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work\n" + pete.String()
		authors       = lib.CoAuthors{tam}
		pairs         = lib.CoAuthors{pete}
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnAuthorsFile(t, authors)
	givenThereIsAPairsFile(t, pairs.Names())

	output, err := runNonInteractiveSelectHook(t)
	assert.Error(t, err)
	assert.Contains(t, output, `author "Pete" is not specified in the authors file`)
}

func runNonInteractiveSelectHook(t *testing.T) (string, error) {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../cmd/select/...",
		fmt.Sprintf("--commitFile=%s", commitFilePath),
		fmt.Sprintf("--authorsFile=%s", authorsFilePath),
		fmt.Sprintf("--pairsFile=%s", pairsFilePath),
	)

	b, err := cmd.CombinedOutput()
	t.Log("CLI output:\n", string(b))
	return stripansi.Strip(string(b)), err
}
