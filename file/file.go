package file

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

type AnalyzeCommand struct {
	Filepath      string
	SearchPattern string
}

func BuildAnalyzeCommand(cmd *cobra.Command, args []string) (AnalyzeCommand, error) {
	if len(args) != 1 {
		return AnalyzeCommand{}, errors.New("filepath is required as an argument")
	}

	search := cmd.Flag("search").Value.String()

	if search == "" {
		return AnalyzeCommand{}, errors.New("Search text is empty")
	}

	return AnalyzeCommand{
		Filepath:      args[0],
		SearchPattern: search,
	}, nil
}

func CheckIfValidFile(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err != nil && os.IsNotExist(err) {
		return false, errors.New("Log file doesn't exist")
	}

	return true, nil
}
