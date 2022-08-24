package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tamj0rd2/coauthor-select/src/selection"
	"github.com/tamj0rd2/coauthor-select/src/validate"
)

func main() {
	log.SetFlags(0)

	printHelp := func() {
		fmt.Println("Usage: coauthor-select <command> [<args>]")
		fmt.Println("Commands:")
		fmt.Println("  select\t\tselect who you're working with")
		fmt.Println("  validate\tvalidates that pushing is allowed")
	}

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "select":
		selection.MakeSelection(os.Args[2:])
	case "validate":
		validate.Validate(os.Args[2:])
	default:
		printHelp()
		os.Exit(1)
	}
}
