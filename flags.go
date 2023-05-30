package main

import (
	"github.com/spf13/pflag"
)

// ParseFlags parses the command-line flags and returns them as strings
// Returns an error if any issue occurred while parsing
func ParseFlags() (string, string, error) {
	// commandFlag represents the command to run when a file changes
	var commandFlag string
	// watchFlag represents the file or directory to watch
	var watchFlag string

	// Bind the commandFlag and watchFlag to the passed flags
	pflag.StringVarP(&commandFlag, "command", "c", "", "The command to run when a file changes")
	pflag.StringVarP(&watchFlag, "watch", "w", "", "The file or directory to watch")

	// Parse the flags
	pflag.Parse()

	return commandFlag, watchFlag, nil
}
