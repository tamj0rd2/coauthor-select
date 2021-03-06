package lib

import (
	"strings"
)

const COMMIT_SEPARATOR = "\n# ------------------------ >8 ------------------------"

func PrepareCommitMessage(input string, coAuthors []CoAuthor) string {
	if len(coAuthors) == 0 {
		return input
	}

	sections := strings.SplitN(input, COMMIT_SEPARATOR, 2)
	message := sections[0] + "\n"

	for _, coAuthor := range coAuthors {
		coauthorLine := coAuthor.String()
		if !strings.Contains(message, coauthorLine) {
			message += "\n" + coAuthor.String()
		}
	}

	if len(sections) > 1 {
		metadataSection := sections[1]
		message = message + COMMIT_SEPARATOR + metadataSection
	}

	return message
}
