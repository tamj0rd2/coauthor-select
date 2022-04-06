package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tamj0rd2/coauthor-select/git"
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

	coAuthors := []git.CoAuthor{{Name: "tamj0rd2", Email: "tam@tam.com"}}

	output := git.PrepareCommitMessage(string(file), coAuthors)

	err = os.WriteFile(commitFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added co-authors:", coAuthors)
}
