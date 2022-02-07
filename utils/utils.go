package utils

import (
	"fmt"
	"os"
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
