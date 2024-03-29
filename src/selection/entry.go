package selection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/tamj0rd2/coauthor-select/src/lib"
)

func MakeSelection(args []string) {
	options, err := parseOptions(args)
	if err != nil {
		log.Fatal(err)
	}

	cliApp := NewCLIApp(
		func(ctx context.Context) (lib.CoAuthors, error) {
			authors, err := readAuthorsFile(options.AuthorsFilePath)
			if err != nil {
				return nil, err
			}

			if !options.Interactive {
				return getCoAuthorsNonInteractive(options.PairsFilePath, authors)
			}

			return getCoAuthorsInteractive(options, authors)
		},
		func(ctx context.Context, pairs lib.CoAuthors) error {
			b, err := json.Marshal(pairs.Names())
			if err != nil {
				return fmt.Errorf("failed to marshal pairs - %w", err)
			}

			if err := os.WriteFile(options.PairsFilePath, b, 0o644); err != nil {
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
			if err := os.WriteFile(options.CommitFilePath, []byte(message), 0o644); err != nil {
				return fmt.Errorf("failed to write commit message file to %q: %w", options.CommitFilePath, err)
			}
			return nil
		},
	)

	if err := cliApp.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func getCoAuthorsInteractive(options selectOptions, authors lib.CoAuthors) ([]lib.CoAuthor, error) {
	previousPairs, wantsToUsePreviousPairs, err := getPreviousPairsInteractive(options)
	if err != nil {
		return nil, err
	}

	if wantsToUsePreviousPairs {
		return authors.Subset(previousPairs), nil
	}

	selectedPairs, err := getPairNamesInteractive(authors.Names())
	if err != nil {
		return nil, err
	}
	return authors.Subset(selectedPairs), nil
}

func getPreviousPairsInteractive(options selectOptions) ([]string, bool, error) {
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

func getPairNamesInteractive(authorNames []string) ([]string, error) {
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

func getCoAuthorsNonInteractive(pairsFilePath string, authors lib.CoAuthors) ([]lib.CoAuthor, error) {
	pairNames, err := readJSON[[]string](pairsFilePath)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "failed to read file") {
			return nil, err
		}
	}

	var coAuthors []lib.CoAuthor
	for _, name := range pairNames {
		author, err := authors.Get(name)
		if err != nil {
			return nil, err
		}

		coAuthors = append(coAuthors, author)
	}

	return coAuthors, nil
}

func readAuthorsFile(filePath string) (lib.CoAuthors, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read authors file %q - %w", filePath, err)
	}

	var authors lib.CoAuthors
	if err := authors.From(b); err != nil {
		return nil, fmt.Errorf("failed to parse authors file %q - %w", filePath, err)
	}

	return authors, nil
}

func readJSON[T any](filePath string) (T, error) {
	var result T
	b, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("failed to read file %q - %w", filePath, err)
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal file %q - %w", filePath, err)
	}

	return result, nil
}
