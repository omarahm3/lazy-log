package tests

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use: "root",
	}
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buffer := new(bytes.Buffer)
	root.SetOut(buffer)
	root.SetErr(buffer)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buffer.String(), err
}

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return output, err
}

func TestError(test *testing.T, actual interface{}, expected interface{}) {
	test.Errorf("\nActual:\t\t%v\nExpected:\t%v", actual, expected)
}
