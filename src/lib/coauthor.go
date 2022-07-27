package lib

import "fmt"

type CoAuthor struct {
	Name  string
	Email string
}

func (c CoAuthor) String() string {
	return fmt.Sprintf("Co-authored-by: %s <%s>", c.Name, c.Email)
}

type CoAuthors []CoAuthor

func (authors CoAuthors) Get(name string) (CoAuthor, error) {
	for _, author := range authors {
		if author.Name == name {
			return author, nil
		}
	}

	return CoAuthor{}, fmt.Errorf("author %s not present in the authors file", name)
}

func (authors CoAuthors) Names() []string {
	var names []string
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
