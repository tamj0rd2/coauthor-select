package main

import (
	"flag"
	"fmt"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"log"
	"os/exec"
	"strings"
)

var (
	options ValidateOptions
)

type ValidateOptions struct {
	CommitFilePath string
	TrunkName      string
	BranchName     string
	ProtectTrunk   bool
}

func init() {
	flag.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&options.TrunkName, "trunkName", "main", "name of the trunk branch")
	flag.StringVar(&options.BranchName, "branchName", "", "name of the branch you're on")
	flag.BoolVar(&options.ProtectTrunk, "protectTrunk", true, "whether you're allowed to commit to the trunk without pairs")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	handleError := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

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

	isUserOnTrunk := strings.ToLower(options.BranchName) == strings.ToLower(options.TrunkName)
	if !isUserOnTrunk || !options.ProtectTrunk {
		fmt.Println("Friendly reminder that you should get some feedback on your work occasionally because you're not pairing")
		return
	}

	log.Fatal(fmt.Errorf("you can't commit to %q without pairing", options.BranchName))
}
