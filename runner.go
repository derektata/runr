package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

// CommandRunner runs a command.
type CommandRunner struct {
	Command   string
	Watcher   *fsnotify.Watcher
	Done      chan bool
	Debouncer *time.Timer
	WatchPath string
}

// NewCommandRunner creates a new CommandRunner.
// It takes in a command to run and a path to watch as a string.
// It returns a new CommandRunner and an error if there was one.
func NewCommandRunner(command, watchPath string) (*CommandRunner, error) {
	// Check if a command was provided
	if command == "" {
		return nil, fmt.Errorf("please provide a command to run")
	}

	// Create a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// If no path was provided, use the current working directory
	if watchPath == "" {
		watchPath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	// Return a new CommandRunner with the provided command, watcher, and path
	return &CommandRunner{
		Command:   command,
		Watcher:   watcher,
		Done:      make(chan bool),
		Debouncer: nil,
		WatchPath: watchPath,
	}, nil
}

// Run executes the command whenever a file changes.
func (cr *CommandRunner) Run() error {
	// Start watching for file events in a separate goroutine.
	go func() {
		defer close(cr.Done)
		cr.watchEvents()
	}()

	// Start watching the directory for file changes.
	return cr.watchDirectory()
}

// watchEvents continuously watches for file system events and triggers a command
// when a file has been changed.
func (cr *CommandRunner) watchEvents() {
	for {
		select {
		// Wait for a file system event
		case event, ok := <-cr.Watcher.Events:
			if !ok {
				return
			}
			// If a file has been changed, debounce the command
			if cr.isFileChanged(event) {
				cr.debounceCommand()
			}
		// Handle any errors that occur while watching the file system
		case err, ok := <-cr.Watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// isFileChanged checks whether a file has been changed based on the given fsnotify.Event.
// Returns true if the event indicates a write, create, or rename operation.
func (cr *CommandRunner) isFileChanged(event fsnotify.Event) bool {
	return event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename
}

// debounceCommand debounces the runCommand method if the Debouncer is not nil.
// It stops the Debouncer if it is already running and starts it again.
// The Debouncer waits for 500ms before executing the runCommand method.
func (cr *CommandRunner) debounceCommand() {
	if cr.Debouncer != nil {
		cr.Debouncer.Stop()
	}
	cr.Debouncer = time.AfterFunc(500*time.Millisecond, cr.runCommand)
}

// runCommand runs the specified command.
func (cr *CommandRunner) runCommand() {
	// Create a new command with "sh" as the executable and "-c" and the argument.
	cmd := exec.Command("sh", "-c", cr.Command)

	// Set the standard output and standard error of the command to use the current process's ones.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and wait for it to complete.
	if err := cmd.Run(); err != nil {
		// If there was an error, log it and exit with a non-zero status code.
		log.Fatal(err)
	}
}

// watchDirectory adds the directory to the watcher and waits for cr.Done signal
// cr: CommandRunner is the struct containing the Watcher and WatchPath
// Returns an error if the directory cannot be added to the watcher
func (cr *CommandRunner) watchDirectory() error {
	if err := cr.Watcher.Add(cr.WatchPath); err != nil { // Add the directory to the watcher
		return err
	}

	<-cr.Done // Wait for the done signal
	return nil
}
