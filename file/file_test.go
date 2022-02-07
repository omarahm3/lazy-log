package file_test

import (
	"bytes"
	"io/ioutil"
	"lazy-log/cmd"
	"lazy-log/file"
	"lazy-log/utils"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newLogAnalyzerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "log-analyzer",
	}
}

func newCommand(analyzeFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := newLogAnalyzerCommand()

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

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buffer := new(bytes.Buffer)
	root.SetOut(buffer)
	root.SetErr(buffer)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buffer.String(), err
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func testError(test *testing.T, actual interface{}, expected interface{}) {
	test.Errorf("\nActual: %v\nExpected: %v", actual, expected)
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
				}

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					testError(t, analyzeCommand, expectedCommand)
				}
			},
		},
		{
			name: "No search pattern",
			args: []string{"analyze", "file.log"},
			function: func(command *cobra.Command, args []string) {
				analyzeCommand, err := cmd.BuildAnalyzeCommand(command, args)
				expectedCommand := cmd.AnalyzeCommand{}
				errorMessage := "Search text is empty"

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					testError(t, analyzeCommand, expectedCommand)
					return
				}

				if err != nil && string(err.Error()) != errorMessage {
					testError(t, err.Error(), errorMessage)
				}
			},
		},
		{
			name: "Empty search pattern",
			args: []string{"analyze", "file.log", "--search="},
			function: func(command *cobra.Command, args []string) {
				analyzeCommand, err := cmd.BuildAnalyzeCommand(command, args)
				expectedCommand := cmd.AnalyzeCommand{}
				errorMessage := "Search text is empty"

				if !reflect.DeepEqual(analyzeCommand, expectedCommand) {
					testError(t, analyzeCommand, expectedCommand)
					return
				}

				if err != nil && string(err.Error()) != errorMessage {
					testError(t, err.Error(), errorMessage)
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
					testError(t, analyzeCommand, expectedCommand)
					return
				}

				if err != nil && string(err.Error()) != errorMessage {
					testError(t, err.Error(), errorMessage)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := newCommand(test.function)

			executeCommand(command, test.args...)
		})
	}
}

func Test_CheckIfValidFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test*.log")

	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name        string
		filename    string
		expected    bool
		expectError bool
	}{
		{
			name:        "File does exist",
			filename:    tmpFile.Name(),
			expected:    true,
			expectError: false,
		},
		{
			name:        "File does not exist",
			filename:    "nothing/test.log",
			expected:    false,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := file.CheckIfValidFile(test.filename)

			if (err != nil) != test.expectError {
				t.Errorf("Error: %v, Expect error: %v", err, test.expectError)
				return
			}

			if result != test.expected {
				testError(t, result, test.expectError)
				return
			}
		})
	}
}

func Test_ScanFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expected    bool
		expectError bool
	}{
		{
			name:        "Scanner scan each line",
			fileContent: "1\n2\n3\n4\n5",
			expected:    true,
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "test.log")

			utils.Check(err)

			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(test.fileContent)

			utils.Check(err)

			// Persist data on disk
			err = tmpFile.Sync()

			utils.Check(err)

      var lines []string

			file.ProcessLogFile(tmpFile.Name(), func(line string) {
        lines = append(lines, line)
			})

      expectedOutput := strings.Split(test.fileContent, "\n")

      if !reflect.DeepEqual(lines, expectedOutput) {
				testError(t, lines, expectedOutput)
				return
      }
		})
	}
}
