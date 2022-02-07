package utils_test

import (
	"fmt"
	tests_helpers "lazy-log/tests"
	"lazy-log/utils"
	"testing"
)

func Test_ExtractJsonFromString(t *testing.T) {
  complexJson := `{"name":"John Doe", "email":"john@example.com", "info": [{"date": "Today", "group": {"name": "sports", "id": 1}}]}`
  complexBadJson := `{"name":"John Doe", "email":"john@example.com", "info": [{"date": "Today", "group": {"name": "sports", "id": 1}}}`
  simpleJson := `{"name":"John Doe", "email":"john@example.com"}`
  simpleBadJson := `{"name":"John Doe", "email":"john@example.com"`

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:  "Extract simple JSON from middle of a string",
			input: fmt.Sprintf("Data was saved %s to DB", simpleJson),
			expected: `{
  "name": "John Doe",
  "email": "john@example.com"
}`,
			expectError: false,
		},
		{
			name:  "Extract simple JSON from end of a string",
			input: fmt.Sprintf("Data was saved to DB %s", simpleJson),
			expected: `{
  "name": "John Doe",
  "email": "john@example.com"
}`,
			expectError: false,
		},
		{
			name:  "Extract JSON from start of a string",
			input: fmt.Sprintf("%s Data was saved to DB", simpleJson),
			expected: `{
  "name": "John Doe",
  "email": "john@example.com"
}`,
			expectError: false,
		},
		{
			name:  "Extract complex JSON from start of a string",
			input: fmt.Sprintf("%s Data was saved to DB", complexJson),
			expected: `{
  "name": "John Doe",
  "email": "john@example.com",
  "info": [
    {
      "date": "Today",
      "group": {
        "name": "sports",
        "id": 1
      }
    }
  ]
}`,
			expectError: false,
		},
		{
			name:  "Extract complex JSON from middle of a string",
			input: fmt.Sprintf("Data was saved %s to DB", complexJson),
			expected: `{
  "name": "John Doe",
  "email": "john@example.com",
  "info": [
    {
      "date": "Today",
      "group": {
        "name": "sports",
        "id": 1
      }
    }
  ]
}`,
			expectError: false,
		},
		{
			name:  "Extract complex JSON from end of a string",
			input: fmt.Sprintf("Data was saved to DB %s", complexJson),
			expected: `{
  "name": "John Doe",
  "email": "john@example.com",
  "info": [
    {
      "date": "Today",
      "group": {
        "name": "sports",
        "id": 1
      }
    }
  ]
}`,
			expectError: false,
		},
		{
			name:  "Invalid JSON with only opening curly bracket",
			input: fmt.Sprintf("Data was saved %s to DB", simpleBadJson),
			expected: "String is not JSON, not completed",
			expectError: true,
		},
		{
			name:  "Invalid JSON with only opening bracket",
			input: fmt.Sprintf("Data was saved %s to DB", complexBadJson),
			expected: "Invalid JSON input, opening [ is not closed",
			expectError: true,
		},
	}

	for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      json, err := utils.ExtractJsonFromString(test.input)

      if err != nil {
        if !test.expectError {
          t.Errorf("Unexpected error: [%v]", err)
          return
        }

        if err.Error() != test.expected {
          tests_helpers.TestError(t, err.Error(), test.expected)
          return
        }
      }


      if err == nil && json != test.expected {
        tests_helpers.TestError(t, json, test.expected)
        return
      }
    })
	}
}
