package cmd_test

import (
	"lazy-log/cmd"
	tests_helpers "lazy-log/tests"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newCommand(analyzeFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := tests_helpers.NewRootCommand()

	analyzeCommand := &cobra.Command{
		Use: "analyze command test",
		Run: analyzeFunc,
	}

	analyzeCommand.Flags().String("search", "", "Search pattern to track in the logs")
	analyzeCommand.Flags().StringSlice("pattern", []string{}, "Search patterns to track in the logs")
	analyzeCommand.Flags().Bool("json", false, "Parse json objects on each log line and pretty print them")

	cmd.AddCommand(analyzeCommand)

	return cmd
}

func Test_BuildAnalyzeCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		function func(cmd *cobra.Command, args []string)
	}{
		{
			name: "Default parameters",
			args: []string{"analyze", "file.log", "--search=test"},
			function: func(command *cobra.Command, args []string) {
				analyzeCommand, _ := cmd.BuildAnalyzeCommand(command, args)
				expectedCommand := cmd.AnalyzeCommand{
					Filepath:      "file.log",
					SearchPattern: "test",
					Pattern:       []string{},
					Json:          false,
				}

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					tests_helpers.TestError(t, analyzeCommand, expectedCommand)
				}
			},
		},
		{
			name: "No search pattern",
			args: []string{"analyze", "file.log"},
			function: func(command *cobra.Command, args []string) {
				analyzeCommand, err := cmd.BuildAnalyzeCommand(command, args)
				expectedCommand := cmd.AnalyzeCommand{
          Filepath: "file.log",
          Pattern: []string{},
          Json: false,
        }

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					tests_helpers.TestError(t, analyzeCommand, expectedCommand)
					return
				}

				if err != nil {
          t.Errorf("Unexpected error: {%v}", err)
				}
			},
		},
		{
			name: "Empty search pattern",
			args: []string{"analyze", "file.log", "--search="},
			function: func(command *cobra.Command, args []string) {
				analyzeCommand, err := cmd.BuildAnalyzeCommand(command, args)
				expectedCommand := cmd.AnalyzeCommand{
          Filepath: "file.log",
          Pattern: []string{},
          Json: false,
        }

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					tests_helpers.TestError(t, analyzeCommand, expectedCommand)
					return
				}

				if err != nil {
          t.Errorf("Unexpected error: {%v}", err)
				}
			},
		},
		{
			name: "No filepath argument",
			args: []string{"analyze"},
			function: func(command *cobra.Command, args []string) {
				analyzeCommand, err := cmd.BuildAnalyzeCommand(command, args)
				expectedCommand := cmd.AnalyzeCommand{}
				errorMessage := "filepath is required as an argument"

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					tests_helpers.TestError(t, analyzeCommand, expectedCommand)
					return
				}

				if err != nil && string(err.Error()) != errorMessage {
					tests_helpers.TestError(t, err.Error(), errorMessage)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := newCommand(test.function)

			tests_helpers.ExecuteCommand(command, test.args...)
		})
	}
}
