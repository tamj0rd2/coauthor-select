package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tamj0rd2/coauthor-select/src/app"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	authorsFilePath string
	commitFilePath  string
	pairsFilePath   string
	trunkName       string
	branchName      string
)

func init() {
	flag.StringVar(&authorsFilePath, "authorsFile", "authors.json", "names & emails of teammates")
	flag.StringVar(&commitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&pairsFilePath, "pairsFile", "pairs.json", "path to pairs file")
	flag.StringVar(&trunkName, "trunkName", "main", "the name of the trunk branch")
	flag.StringVar(&branchName, "branchName", "", "the branch you're currently on")
}

func main() {
	var err error
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	cliApp := app.NewCLIApp(
		func(ctx context.Context) (lib.CoAuthors, error) {
			return getCoAuthors(authorsFilePath, pairsFilePath)
		},
		func(authors lib.CoAuthors) (string, error) {
			file, err := os.ReadFile(commitFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to read commit message file: %w", err)
			}
			return lib.PrepareCommitMessage(string(file), authors), nil
		},
		func(ctx context.Context, message string) error {
			if err := os.WriteFile(commitFilePath, []byte(message), 0644); err != nil {
				return fmt.Errorf("failed to write commit message file to %q: %w", commitFilePath, err)
			}
			return nil
		},
	)

	ctx := context.Background()
	if branchName == "" {
		branchName, err = getBranchName(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cliApp.Run(ctx, trunkName, branchName); err != nil {
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

func getCoAuthors(authorsFilePath string, pairsFilePath string) ([]lib.CoAuthor, error) {
	authorsFile, err := os.ReadFile(authorsFilePath) // TODO: make filepath configurable
	if err != nil {
		return nil, err
	}

	var authors lib.CoAuthors
	err = json.NewDecoder(bytes.NewReader(authorsFile)).Decode(&authors)
	if err != nil {
		return nil, err
	}

	// TODO: be able to load the pairs from somewhere other than a file. i.e discord
	pairFile, err := os.ReadFile(pairsFilePath)
	if err != nil {
		return nil, err
	}

	var pairs []string
	err = json.NewDecoder(bytes.NewReader(pairFile)).Decode(&pairs)
	if err != nil {
		return nil, err
	}

	var coAuthors []lib.CoAuthor
	for _, name := range pairs {
		author, err := authors.Get(name)
		if err != nil {
			return nil, err
		}

		coAuthors = append(coAuthors, author)
	}

	return coAuthors, nil
}
