package validate

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/tamj0rd2/coauthor-select/src/lib"
)

type ValidateOptions struct {
	CommitFilePath string
	TrunkName      string
	BranchName     string
	ProtectTrunk   bool
}

func Validate(args []string) {
	flags := flag.NewFlagSet("validate", flag.ExitOnError)

	var options ValidateOptions
	flags.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flags.StringVar(&options.TrunkName, "trunkName", "main", "name of the trunk branch")
	flags.StringVar(&options.BranchName, "branchName", "", "name of the branch you're on")
	flags.BoolVar(&options.ProtectTrunk, "protectTrunk", true, "whether you're allowed to commit to the trunk without pairs")

	handleError := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	handleError(flags.Parse(args))

	if options.BranchName == "" {
		b, err := exec.Command("git", "branch", "--show-current").Output()
		handleError(err)

		options.BranchName = strings.TrimSpace(string(b))
		if options.BranchName == "" {
			log.Fatal("failed to get branch name - are you in a detached head state or in the middle of a rebase?")
		}
	}

	commitFile, err := lib.LoadCommitMessage(options.CommitFilePath)
	handleError(err)

	if lib.DoesCommitContainCoAuthors(commitFile) {
		return
	}

	isUserOnTrunk := strings.EqualFold(options.BranchName, options.TrunkName)
	if !isUserOnTrunk || !options.ProtectTrunk {
		fmt.Println("Friendly reminder that 🫵 you should get some feedback on your work occasionally because you're not pairing")
		return
	}

	log.Fatal(newPairsRequiredError(options.TrunkName))
}

func newPairsRequiredError(trunkName string) error {
	message := fmt.Sprintf("🚨 - Can't commit to %s without a pair", trunkName)
	message += "\nOptions:"
	message += "\n  - get someone to quickly jump in and review your changes so you can select them as a pair for this commit"
	message += "\n  - checkout a branch, make commits on there and make a PR when you're ready for review"
	return errors.New(message)
}
