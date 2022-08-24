package lib

import (
	"fmt"
	"os"
	"strings"
)

const CommitSeparator = "\n# ------------------------ >8 ------------------------"

func PrepareCommitMessage(input string, coAuthors []CoAuthor) string {
	if len(coAuthors) == 0 {
		return input
	}

	sections := strings.SplitN(input, CommitSeparator, 2)
	message := sections[0] + "\n"

	for _, coAuthor := range coAuthors {
		coauthorLine := coAuthor.String()
		if !strings.Contains(message, coauthorLine) {
			message += "\n" + coAuthor.String()
		}
	}

	if len(sections) > 1 {
		metadataSection := sections[1]
		message = message + CommitSeparator + metadataSection
	}

	return message
}

func LoadCommitMessage(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read commit message file %q - %w", filePath, err)
	}
	return string(b), nil
}

func DoesCommitContainCoAuthors(input string) bool {
	return strings.Contains(input, "Co-authored-by: ")
}
