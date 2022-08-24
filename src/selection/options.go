package selection

import (
	"flag"
	"os"

	"github.com/mattn/go-isatty"
)

type selectOptions struct {
	CommitFilePath     string
	AuthorsFilePath    string
	PairsFilePath      string
	ForceSearchPrompts bool
	Interactive        bool
}

func parseOptions(args []string) (selectOptions, error) {
	options := selectOptions{
		Interactive: isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
	}

	flags := flag.NewFlagSet("select", flag.ExitOnError)

	flags.StringVar(&options.AuthorsFilePath, "authorsFile", ".coauthors", "names & emails of teammates")
	flags.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flags.StringVar(&options.PairsFilePath, "pairsFile", "pairs.json", "path to pairs file")
	flags.BoolVar(&options.ForceSearchPrompts, "forceSearchPrompts", false, "makes all prompts searches for ease of testing")
	flags.BoolVar(&options.Interactive, "interactive", options.Interactive, "whether you're using an interactive terminal")
	err := flags.Parse(args)

	return options, err
}
