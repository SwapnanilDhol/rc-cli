package main

import (
	"fmt"
	"os"

	"revenuecat-cli/cmd"
)

func main() {
	if err := cmd.Execute(skillsFS, skillsRoot); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
