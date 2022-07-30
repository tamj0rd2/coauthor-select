package lib

import (
	"fmt"
	"regexp"
	"strings"
)

type CoAuthor struct {
	Name  string
	Email string
}

func newCoAuthor(name string, email string) CoAuthor {
	return CoAuthor{Name: name, Email: email}
}

func (c CoAuthor) String() string {
	return fmt.Sprintf("Co-authored-by: %s", c.UserID())
}

func (c CoAuthor) UserID() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Email)
}

type CoAuthors []CoAuthor

func (authors CoAuthors) String() string {
	var coAuthorStrings []string
	for _, author := range authors {
		coAuthorStrings = append(coAuthorStrings, author.UserID())
	}
	return strings.Join(coAuthorStrings, "\n")
}

func (authors CoAuthors) Get(name string) (CoAuthor, error) {
	for _, author := range authors {
		if author.Name == name {
			return author, nil
		}
	}

	validNames := strings.Join(authors.Names(), ", ")
	return CoAuthor{}, fmt.Errorf("author %q is not specified in the authors file. valid options: [%s]", name, validNames)
}

func (authors CoAuthors) Names() []string {
	names := []string{} // this is instantiated so that it won't be nil.
	for _, author := range authors {
		names = append(names, author.Name)
	}

	return names
}

func (authors CoAuthors) Any() bool {
	return len(authors) > 0
}

func (authors CoAuthors) Subset(names []string) []CoAuthor {
	var subset CoAuthors
	for _, name := range names {
		if author, err := authors.Get(name); err == nil {
			subset = append(subset, author)
		}
	}
	return subset
}

func (authors *CoAuthors) From(bytes []byte) error {
	rxp := regexp.MustCompile(`(.*) <(.*)>`)
	matchGroups := rxp.FindAllStringSubmatch(string(bytes), -1)
	for _, group := range matchGroups {
		if len(group) != 3 {
			return fmt.Errorf("invalid author format: %s. Should be like: Name <email@example.com>", group[0])
		}
		*authors = append(*authors, newCoAuthor(group[1], group[2]))
	}
	return nil
}
