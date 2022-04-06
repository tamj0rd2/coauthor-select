package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tamj0rd2/coauthor-select/domain"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage:\ncoauthor-apply <messageFilePath>")
	}

	commitFilePath := os.Args[1]
	file, err := os.ReadFile(commitFilePath)
	if err != nil {
		log.Fatal(err)
	}

	coAuthors, err := getCoAuthors()
	if err != nil {
		log.Fatal(err)
	}

	output := domain.PrepareCommitMessage(string(file), coAuthors)

	err = os.WriteFile(commitFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added co-authors:", coAuthors)
}

func getCoAuthors() ([]domain.CoAuthor, error) {
	authorsFile, err := os.ReadFile("authors.json") // TODO: make filepath configurable
	if err != nil {
		return nil, err
	}

	var authors Authors
	err = json.NewDecoder(bytes.NewReader(authorsFile)).Decode(&authors)
	if err != nil {
		return nil, err
	}

	// TODO: be able to load the pairs from somewhere other than a file. i.e discord
	pairFile, err := os.ReadFile("pair.json")
	if err != nil {
		return nil, err
	}

	var pairs []string
	err = json.NewDecoder(bytes.NewReader(pairFile)).Decode(&pairs)
	if err != nil {
		return nil, err
	}

	var coAuthors []domain.CoAuthor
	for _, name := range pairs {
		author, err := authors.Get(name)
		if err != nil {
			return nil, err
		}

		coAuthors = append(coAuthors, author)
	}

	return coAuthors, nil
}

type Authors []domain.CoAuthor

func (authors Authors) Get(name string) (domain.CoAuthor, error) {
	for _, author := range authors {
		if author.Name == name {
			return author, nil
		}
	}

	return domain.CoAuthor{}, fmt.Errorf("author %s not present in the authors file", name)
}
