package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tamj0rd2/coauthor-select/lib"
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

	output := lib.PrepareCommitMessage(string(file), coAuthors)

	err = os.WriteFile(commitFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added co-authors:", coAuthors)
}

func getCoAuthors() ([]lib.CoAuthor, error) {
	authorsFile, err := os.ReadFile("authors.json") // TODO: make filepath configurable
	if err != nil {
		return nil, err
	}

	var authors lib.CoAuthors
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

	var coAuthors []lib.CoAuthor
	for _, name := range pairs {
		author, err := authors.Get(name)
		if err != nil {
			return nil, err
		}

		coAuthors = append(coAuthors, author)
	}

	return coAuthors, nil
}
