package main

import (
	"flag"
	"github.com/mattn/go-isatty"
	"os"
)

type selectOptions struct {
	CommitFilePath     string
	AuthorsFilePath    string
	PairsFilePath      string
	ForceSearchPrompts bool
	Interactive        bool
}

func parseOptions() selectOptions {
	var (
		options = selectOptions{
			Interactive: isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
		}
	)

	flag.StringVar(&options.AuthorsFilePath, "authorsFile", ".coauthors", "names & emails of teammates")
	flag.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&options.PairsFilePath, "pairsFile", "pairs.json", "path to pairs file")
	flag.BoolVar(&options.ForceSearchPrompts, "forceSearchPrompts", false, "makes all prompts searches for ease of testing")
	flag.BoolVar(&options.Interactive, "interactive", options.Interactive, "whether you're using an interactive terminal")
	flag.Parse()

	return options
}
