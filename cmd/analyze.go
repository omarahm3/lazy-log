package cmd

import (
	"fmt"
	"lazy-log/file"
	"lazy-log/utils"

	"github.com/spf13/cobra"
)

func analyze(cmd *cobra.Command, args []string) {
	analyzeCommand, err := file.BuildAnalyzeCommand(cmd, args)

	if err != nil {
		utils.ExitGracefully(err)
	}

	if _, err := file.CheckIfValidFile(analyzeCommand.Filepath); err != nil {
		utils.ExitGracefully(err)
	}

  file.ProcessLogFile(analyzeCommand, func(line string) {
  })

	fmt.Println(analyzeCommand)
}

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Begin analyzing a log file",
	Run: analyze,
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
}
