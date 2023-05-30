# runr

runr is a Go program designed to watch a directory or a file for changes and execute a given command whenever a change is detected. This is particularly useful in development environments where you want to automatically run tests, rebuild your application, or perform any other action when your code changes. The program efficiently uses Golang channels for managing when the command should be run in response to file changes.

## Installation

You can install runr by either downloading a pre-compiled binary from the [releases page] or by manually compiling it from the source code. 

## Usage

To run runr, use the following flags:

- `-c` or `--command`: This is the command you want to execute when a file changes. This flag is required.
- `-w` or `--watch`: This is the file or directory you want to watch. If no path is provided, the current working directory is used.

Here's an example of using runr to automatically run tests when any file changes in the current directory:

```bash
./runr -c "go test" -w .
```

## Manual Compilation

If you prefer to manually compile runr, you need to have Go installed on your machine. 

First, clone this repository to your local machine using the following command:

```bash
git clone https://github.com/derektata/runr.git
```

Next, navigate into the cloned repository and install the required dependencies. runr uses the following external Go packages:

1. `github.com/fsnotify/fsnotify`: For file system notifications.
2. `github.com/spf13/pflag`: For command-line flag parsing.

You can install these packages using the following commands:

```bash
go get github.com/fsnotify/fsnotify
go get github.com/spf13/pflag
```

Once the dependencies are installed, you can compile the project using:

```bash
go build -o runr main.go
```

This will create an executable file named `runr` in the current directory.

## How It Works

The `CommandRunner` struct contains all the necessary information to watch a directory and run a command when a file changes. It includes a `fsnotify.Watcher` to watch for file changes, a command to run, a `time.Timer` for debouncing multiple rapid file changes, and a `chan bool` to signal when watching should stop.

The `NewCommandRunner` function is used to create a new `CommandRunner`. It checks whether a command and a path to watch were provided, creates a new file watcher, and returns a new `CommandRunner`.

The `Run` method of `CommandRunner` starts watching the provided directory for file changes and launches a goroutine that listens for file events. When an event indicating a file change is received, it debounces and runs the specified command.

Debouncing is used to ensure that the command is not run too often. If a file changes and then quickly changes again, the command will only be run once after a 500ms delay. This is achieved through the use of Golang's channels and timers, providing a way to effectively manage when the command should be run in response to file changes.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE] file for details.


[releases page]: https://github.com/username/runr/releases
[LICENSE]: https://github.com/derektata/runr/blob/main/LICENSE
