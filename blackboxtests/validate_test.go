package blackboxtests

import (
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

const (
	trunk  = "trunk"
	branch = "not-trunk"
)

func Test_ValidateHook_WorkingAlone_OnTrunk(t *testing.T) {
	t.Cleanup(cleanup)

	givenThereIsACommitMessageFile(t, "feat-376 Did some work")

	output, err := runValidateHook(t, "trunk", true)
	assert.Error(t, err)
	assert.Contains(t, output, `Can't commit to trunk without a pair`)
}

func Test_ValidateHook_WorkingAlone_OnTrunk_WithProtectionOn(t *testing.T) {
	t.Cleanup(cleanup)

	givenThereIsACommitMessageFile(t, "feat-376 Did some work")

	output, err := runValidateHook(t, "trunk", false)
	assert.NoError(t, err)
	assert.Contains(t, output, `you should get some feedback on your work occasionally`)
}

func Test_ValidateHook_WorkingAlone_OnBranch(t *testing.T) {
	t.Cleanup(cleanup)

	givenThereIsACommitMessageFile(t, "feat-376 Did some work")

	output, err := runValidateHook(t, branch, false)
	assert.NoError(t, err)
	assert.Contains(t, output, `you should get some feedback on your work occasionally`)
}

func Test_ValidateHook_Pairing_OnTrunk(t *testing.T) {
	t.Cleanup(cleanup)

	givenThereIsACommitMessageFile(t, "feat-376 Did some work\n"+tam.String())

	_, err := runValidateHook(t, "trunk", false)
	assert.NoError(t, err)
}

func Test_ValidateHook_Pairing_OnBranch(t *testing.T) {
	t.Cleanup(cleanup)

	givenThereIsACommitMessageFile(t, "feat-376 Did some work\n"+tam.String())

	_, err := runValidateHook(t, branch, false)
	assert.NoError(t, err)
}

func runValidateHook(t *testing.T, branchName string, protectTrunk bool) (string, error) {
	cmd := exec.Command("go", "run", "../...", "validate",
		fmt.Sprintf("--commitFile=%s", commitFilePath),
		fmt.Sprintf("--trunkName=%s", trunk),
		fmt.Sprintf("--branchName=%s", branchName),
		fmt.Sprintf("--protectTrunk=%s", strconv.FormatBool(protectTrunk)),
	)

	b, err := cmd.CombinedOutput()
	output := string(b)
	t.Log("CLI output:\n", output)
	return output, err
}
