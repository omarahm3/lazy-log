package cmd

import (
	"errors"
	"fmt"
	"lazy-log/file"
	"lazy-log/utils"
	"regexp"

	"github.com/spf13/cobra"
)

type AnalyzeCommand struct {
	Filepath      string
	SearchPattern string
	Pattern       []string
	Json          bool
}

func BuildAnalyzeCommand(cmd *cobra.Command, args []string) (AnalyzeCommand, error) {
	if len(args) != 1 {
		return AnalyzeCommand{}, errors.New("filepath is required as an argument")
	}

	pattern, err := cmd.Flags().GetStringSlice("pattern")
	utils.Check(err)

	json, err := cmd.Flags().GetBool("json")
	utils.Check(err)

	search := cmd.Flag("search").Value.String()

	return AnalyzeCommand{
		Filepath:      args[0],
		SearchPattern: search,
		Pattern:       pattern,
		Json:          json,
	}, nil
}

func regexMatch(content string, pattern string) bool {
	match, err := regexp.Match(pattern, []byte(content))

	utils.Check(err)

	return match
}

func analyzeLine(analyzeCommand AnalyzeCommand, line string) {
	var processedLine string

	if analyzeCommand.SearchPattern != "" {
		match := regexMatch(line, analyzeCommand.SearchPattern)

		if !match {
			return
		}

		processedLine = line
	}

	if len(analyzeCommand.Pattern) == 0 {
		json, start, end, err := utils.ExtractJsonFromString(line)

		utils.Check(err)

		processedLine = processedLine[:start] + "[json-object]" + processedLine[end:]

		processedLine += fmt.Sprintf("\n[json-object] => %s", json)

		fmt.Println(processedLine)
	} else {
		for _, pattern := range analyzeCommand.Pattern {
			match := regexMatch(processedLine, pattern)

			if !match || !analyzeCommand.Json {
				continue
			}

			json, start, end, err := utils.ExtractJsonFromString(line)

			utils.Check(err)

			processedLine = processedLine[:start] + "[json-object]" + processedLine[end:]

			processedLine += fmt.Sprintf("\n[json-object] => %s", json)

			fmt.Println(processedLine)
		}
	}
}

func analyze(cmd *cobra.Command, args []string) {
	analyzeCommand, err := BuildAnalyzeCommand(cmd, args)

	if err != nil {
		utils.ExitGracefully(err)
	}

	if _, err := file.CheckIfValidFile(analyzeCommand.Filepath); err != nil {
		utils.ExitGracefully(err)
	}

	file.ProcessLogFile(analyzeCommand.Filepath, func(line string) {
		analyzeLine(analyzeCommand, line)
	})
}

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Begin analyzing a log file",
	Run:   analyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyzeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	analyzeCmd.Flags().String("search", "", "Search pattern to track in the logs")
	analyzeCmd.Flags().StringSlice("pattern", []string{}, "Search patterns to track in the logs")
	analyzeCmd.Flags().Bool("json", false, "Parse json objects on each log line and pretty print them")
}
