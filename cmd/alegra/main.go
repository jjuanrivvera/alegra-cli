// Command alegra is a command-line interface for the Alegra accounting API.
package main

import (
	"fmt"
	"os"

	"github.com/jjuanrivvera/alegra-cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
