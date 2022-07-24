package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"strings"
)

type GetPairs func(ctx context.Context) (lib.CoAuthors, error)

type FormatCommitMessage func(authors lib.CoAuthors) (string, error)

type SaveCommitMessage func(ctx context.Context, message string) error

type CLIApp struct {
	getPairs            GetPairs
	formatCommitMessage FormatCommitMessage
	saveCommitMessage   SaveCommitMessage
}

func NewCLIApp(
	getPairs GetPairs,
	formatCommitMessage FormatCommitMessage,
	saveCommitMessage SaveCommitMessage,
) *CLIApp {
	return &CLIApp{
		getPairs:            getPairs,
		formatCommitMessage: formatCommitMessage,
		saveCommitMessage:   saveCommitMessage,
	}
}

func (c CLIApp) Run(ctx context.Context, trunkName string, branchName string) error {
	pairs, err := c.getPairs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pairs: %w", err)
	}

	if !pairs.Any() {
		if strings.ToLower(branchName) == strings.ToLower(trunkName) {
			return newPairsRequiredError(branchName)
		}
		return nil
	}

	commitMessage, err := c.formatCommitMessage(pairs)
	if err != nil {
		return fmt.Errorf("failed to format commit message: %w", err)
	}

	if err = c.saveCommitMessage(ctx, commitMessage); err != nil {
		return fmt.Errorf("failed to save commit message: %w", err)
	}

	fmt.Println("Added co-authors:", pairs)
	return nil
}

func newPairsRequiredError(trunkName string) error {
	message := fmt.Sprintf("can't commit to %s without a pair", trunkName)
	message += "\nOptions:"
	message += "\n  - get someone to quickly jump in and review your changes so you can select them as a pair for this commit"
	message += "\n  - checkout a branch, make commits on there and make a PR when you're ready for review"
	return errors.New(message)
}
