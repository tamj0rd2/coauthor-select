package domain

import "fmt"

type CoAuthor struct {
	Name  string
	Email string
}

func NewCoAuthor(name string, email string) CoAuthor {
	return CoAuthor{
		Name:  name,
		Email: email,
	}
}

func (c CoAuthor) String() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Email)
}
