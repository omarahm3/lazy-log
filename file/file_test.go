package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
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
			function: func(cmd *cobra.Command, args []string) {
				analyzeCommand, _ := BuildAnalyzeCommand(cmd, args)
				expectedCommand := AnalyzeCommand{
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
			function: func(cmd *cobra.Command, args []string) {
				analyzeCommand, err := BuildAnalyzeCommand(cmd, args)
				expectedCommand := AnalyzeCommand{}
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
			function: func(cmd *cobra.Command, args []string) {
				analyzeCommand, err := BuildAnalyzeCommand(cmd, args)
				expectedCommand := AnalyzeCommand{}
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
			function: func(cmd *cobra.Command, args []string) {
				analyzeCommand, err := BuildAnalyzeCommand(cmd, args)
				expectedCommand := AnalyzeCommand{}
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
			result, err := CheckIfValidFile(test.filename)

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
