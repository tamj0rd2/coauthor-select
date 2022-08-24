package selection

import (
	"context"
	"fmt"

	"github.com/tamj0rd2/coauthor-select/src/lib"
)

type GetPairs func(ctx context.Context) (lib.CoAuthors, error)

type SavePairs func(ctx context.Context, pairs lib.CoAuthors) error

type FormatCommitMessage func(authors lib.CoAuthors) (string, error)

type SaveCommitMessage func(ctx context.Context, message string) error

type CLIApp struct {
	getPairs            GetPairs
	savePairs           SavePairs
	formatCommitMessage FormatCommitMessage
	saveCommitMessage   SaveCommitMessage
}

func NewCLIApp(
	getPairs GetPairs,
	savePairs SavePairs,
	formatCommitMessage FormatCommitMessage,
	saveCommitMessage SaveCommitMessage,
) *CLIApp {
	return &CLIApp{
		getPairs:            getPairs,
		savePairs:           savePairs,
		formatCommitMessage: formatCommitMessage,
		saveCommitMessage:   saveCommitMessage,
	}
}

func (c CLIApp) Run(ctx context.Context) error {
	pairs, err := c.getPairs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pairs: %w", err)
	}

	if err := c.savePairs(ctx, pairs); err != nil {
		// it's really not the end of the world. No need to kill the program.
		return fmt.Errorf("failed to save pairs: %w", err)
	}

	if !pairs.Any() {
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
