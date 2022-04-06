package lib

import (
	"fmt"
	"strings"
)

const COMMIT_SEPARATOR = "\n# ------------------------ >8 ------------------------"

type CoAuthor struct {
	Name  string
	Email string
}

func (c CoAuthor) String() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Email)
}

func PrepareCommitMessage(input string, coAuthors []CoAuthor) string {
	if len(coAuthors) == 0 {
		return input
	}

	sections := strings.SplitN(input, COMMIT_SEPARATOR, 2)

	message := sections[0] + "\n"
	for _, coAuthor := range coAuthors {
		message += fmt.Sprintf("\nCo-authored-by: %s", coAuthor.String())
	}

	if len(sections) > 1 {
		metadataSection := sections[1]
		message = message + COMMIT_SEPARATOR + metadataSection
	}

	return message
}
