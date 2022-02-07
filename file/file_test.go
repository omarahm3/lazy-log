package file_test

import (
	"io/ioutil"
	"lazy-log/file"
	tests_helpers "lazy-log/tests"
	"lazy-log/utils"
	"os"
	"reflect"
	"strings"
	"testing"
)

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
				tests_helpers.TestError(t, result, test.expectError)
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
				tests_helpers.TestError(t, lines, expectedOutput)
				return
			}
		})
	}
}
