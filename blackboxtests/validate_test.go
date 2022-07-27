package blackboxtests

import (
	"github.com/alecthomas/assert/v2"
	"os/exec"
	"strconv"
	"testing"
)

func Test_ValidateHook_WorkingAlone_OnTrunk(t *testing.T) {
	givenThereIsACommitMessageFile(t, "feat-376 Did some work")

	output, err := runValidateHook(t, "trunk", "trunk", false)
	assert.Error(t, err)
	assert.Contains(t, output, `you can't commit to "trunk" without pairing`)
}

func Test_ValidateHook_WorkingAlone_OnBranch(t *testing.T) {
	givenThereIsACommitMessageFile(t, "feat-376 Did some work")

	output, err := runValidateHook(t, "trunk", "not-trunk", false)
	assert.NoError(t, err)
	assert.Contains(t, output, `you should get some feedback on your work occasionally`)
}

func Test_ValidateHook_Pairing_OnTrunk(t *testing.T) {
	givenThereIsACommitMessageFile(t, "feat-376 Did some work\n"+tam.String())

	_, err := runValidateHook(t, "trunk", "trunk", false)
	assert.NoError(t, err)
}

func Test_ValidateHook_Pairing_OnBranch(t *testing.T) {
	givenThereIsACommitMessageFile(t, "feat-376 Did some work\n"+tam.String())

	_, err := runValidateHook(t, "trunk", "not-trunk", false)
	assert.NoError(t, err)
}

func runValidateHook(t *testing.T, trunkName string, branchName string, protectTrunk bool) (string, error) {
	cmd := exec.Command("go", "run", "../cmd/validate/...",
		"--commitFile", commitFilePath,
		"--trunkName", trunkName,
		"--branchName", branchName,
		"--protectTrunk", strconv.FormatBool(protectTrunk),
	)

	b, err := cmd.CombinedOutput()
	output := string(b)
	t.Log("CLI output:\n", output)
	return output, err
}
