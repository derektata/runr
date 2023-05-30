package main

import (
	"log"
)

func main() {
	// Parse the flags
	commandFlag, watchFlag, err := ParseFlags()
	if err != nil {
		// If an error occurred, log and exit
		log.Fatal(err)
	}

	// Create a new CommandRunner with the passed commandFlag and watchFlag
	cr, err := NewCommandRunner(commandFlag, watchFlag)
	if err != nil {
		// If an error occurred, log and exit
		log.Fatal(err)
	}

	// Run the CommandRunner
	if err = cr.Run(); err != nil {
		// If an error occurred, log and exit
		log.Fatal(err)
	}
}
