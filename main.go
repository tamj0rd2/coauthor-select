package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tamj0rd2/coauthor-select/lib"
)

var (
	authorsFilePath string
	commitFilePath  string
	pairsFilePath   string
)

func init() {
	flag.StringVar(&authorsFilePath, "authorsFile", "authors.json", "names & emails of teammates")
	flag.StringVar(&commitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&pairsFilePath, "pairsFile", "pairs.json", "path to pairs file")
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	file, err := os.ReadFile(commitFilePath)
	if err != nil {
		log.Fatal(err)
	}

	coAuthors, err := getCoAuthors(authorsFilePath, pairsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if len(coAuthors) == 0 {
		return
	}

	updatedCommitMessage := lib.PrepareCommitMessage(string(file), coAuthors)
	if err = os.WriteFile(commitFilePath, []byte(updatedCommitMessage), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added co-authors:", coAuthors)
}

func getCoAuthors(authorsFilePath string, pairsFilePath string) ([]lib.CoAuthor, error) {
	authorsFile, err := os.ReadFile(authorsFilePath) // TODO: make filepath configurable
	if err != nil {
		return nil, err
	}

	var authors lib.CoAuthors
	err = json.NewDecoder(bytes.NewReader(authorsFile)).Decode(&authors)
	if err != nil {
		return nil, err
	}

	// TODO: be able to load the pairs from somewhere other than a file. i.e discord
	pairFile, err := os.ReadFile(pairsFilePath)
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
