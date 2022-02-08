package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ExitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func Check(e error) {
	if e != nil {
		ExitGracefully(e)
	}
}

func sliceToString(slice []string) string {
	return strings.Join(slice, "")
}

func prettyJsonString(content string) string {
	buffer := &bytes.Buffer{}

	err := json.Indent(buffer, []byte(content), "", "  ")

	Check(err)

	return buffer.String()
}

// This will actually do the JSON extraction from the string
// It will make sure that it parses an string
// And then it will return a pretty JSON string
func extractJsonFromSubString(substring string) (string, int, error) {
	stack := []string{}
	var jsonString string

	for _, char := range substring {
		if strings.ContainsAny(string(char), "{[") {
			stack = append(stack, string(char))
			jsonString += string(char)
			continue
		}

		if len(stack) == 0 {
			return prettyJsonString(jsonString), len(jsonString), nil
		}

		switch string(char) {
		case "}":
			check := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if check == "[" {
				return "", 0, errors.New("Invalid JSON input, opening [ is not closed")
			}
		case "]":
			check := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if check == "{" {
				return "", 0, errors.New("Invalid JSON input, opening { is not closed")
			}
		}
		jsonString += string(char)
	}

	if len(stack) > 0 {
		return "", 0, errors.New("String is not JSON, not completed")
	}

	return prettyJsonString(jsonString), len(jsonString), nil
}

func ExtractJsonFromString(content string) (string, int, int, error) {
	startPosition := strings.Index(content, "{")
	potentialJson := content[startPosition:]

	json, originalStringLength, err := extractJsonFromSubString(potentialJson)

	endPosition := originalStringLength + startPosition

	if err != nil {
		startPosition = 0
		endPosition = 0
	}

	return json, startPosition, endPosition, err
}
