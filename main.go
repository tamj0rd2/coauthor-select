package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/tamj0rd2/coauthor-select/src"
	"github.com/tamj0rd2/coauthor-select/src/app"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	options src.Options
)

func init() {
	flag.StringVar(&options.AuthorsFilePath, "authorsFile", "authors.json", "names & emails of teammates")
	flag.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&options.PairsFilePath, "pairsFile", "pairs.json", "path to pairs file")
	flag.StringVar(&options.TrunkName, "trunkName", "main", "the name of the trunk branch")
	flag.StringVar(&options.BranchName, "branchName", "", "the branch you're currently on")
	flag.BoolVar(&options.ForceSearchPrompts, "forceSearchPrompts", false, "makes all prompts searches for ease of testing")
	flag.BoolVar(&options.ProtectTrunk, "protectTrunk", true, "whether or not to allow working solo on the trunk")
}

func main() {
	var err error
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()
	if options.BranchName == "" {
		options.BranchName, err = getBranchName(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	cliApp := app.NewCLIApp(
		func(ctx context.Context) (lib.CoAuthors, error) {
			return getCoAuthors()
		},
		func(ctx context.Context, pairs lib.CoAuthors) error {
			b, err := json.Marshal(pairs.Names())
			if err != nil {
				return fmt.Errorf("failed to marshal pairs - %w", err)
			}

			if err := os.WriteFile(options.PairsFilePath, b, 0644); err != nil {
				return fmt.Errorf("failed to write pairs file - %w", err)
			}
			return nil
		},
		func(authors lib.CoAuthors) (string, error) {
			file, err := os.ReadFile(options.CommitFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to read commit message file: %w", err)
			}
			return lib.PrepareCommitMessage(string(file), authors), nil
		},
		func(ctx context.Context, message string) error {
			if err := os.WriteFile(options.CommitFilePath, []byte(message), 0644); err != nil {
				return fmt.Errorf("failed to write commit message file to %q: %w", options.CommitFilePath, err)
			}
			return nil
		},
	)

	if err := cliApp.Run(ctx, options.TrunkName, options.BranchName, options.ProtectTrunk); err != nil {
		log.Fatal(err)
	}
}

func getBranchName(ctx context.Context) (string, error) {
	b, err := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	output := strings.TrimSpace(string(b))
	if err != nil {
		return "", fmt.Errorf("failed to get branch name: %w - %s", err, output)
	}
	return output, nil
}

func getCoAuthors() ([]lib.CoAuthor, error) {
	authorsFile, err := os.ReadFile(options.AuthorsFilePath) // TODO: make filepath configurable
	if err != nil {
		return nil, err
	}

	var authors lib.CoAuthors
	err = json.NewDecoder(bytes.NewReader(authorsFile)).Decode(&authors)
	if err != nil {
		return nil, err
	}

	previousPairs, wantsToUsePreviousPairs, err := getPreviousPairs()
	if err != nil {
		return nil, err
	}

	if wantsToUsePreviousPairs {
		return authors.Subset(previousPairs), nil
	}

	selectedPairs, err := getPairNames(authors.Names())
	if err != nil {
		return nil, err
	}
	return authors.Subset(selectedPairs), nil
}

func getPreviousPairs() ([]string, bool, error) {
	var pairs []string
	pairFile, err := os.ReadFile(options.PairsFilePath)
	if err != nil {
		return nil, false, nil
	}

	if err = json.NewDecoder(bytes.NewReader(pairFile)).Decode(&pairs); err != nil {
		return nil, false, fmt.Errorf("failed to decode pairs file %q - %w", options.PairsFilePath, err)
	}

	if len(pairs) == 0 {
		return nil, false, nil
	}

	yesOrNo := []string{"Yes", "No"}
	prompt := promptui.Select{
		Label:             fmt.Sprintf("Are you still working with these exact people? [%s]", strings.Join(pairs, ", ")),
		Items:             []string{"Yes", "No"},
		StartInSearchMode: options.ForceSearchPrompts,
		Searcher:          newSearcher(yesOrNo),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, false, fmt.Errorf("failed to figure out if you're still pairing with the same people: %w", err)
	}

	return pairs, result == "Yes", nil
}

func getPairNames(authorNames []string) ([]string, error) {
	const noOneElse = "No one else"
	authorNamesToChooseFrom := append([]string{noOneElse}, authorNames...)
	var selectedPairs []string

	for {
		pairSelection := promptui.Select{
			Label:             "Who else are you working with?",
			Items:             authorNamesToChooseFrom,
			StartInSearchMode: true,
			Searcher:          newSearcher(authorNamesToChooseFrom),
		}

		_, pairName, err := pairSelection.Run()
		if err != nil {
			return nil, fmt.Errorf("failed to select pair - %w", err)
		}

		if pairName == noOneElse {
			return selectedPairs, nil
		}

		selectedPairs = append(selectedPairs, pairName)
	}
}

func newSearcher(items []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		name := strings.ToLower(items[index])
		return strings.Contains(name, strings.ToLower(input))
	}
}
