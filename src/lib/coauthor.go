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

func (authors CoAuthors) Any() bool {
	return len(authors) > 0
}
