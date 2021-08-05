package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// vprintf is a printf that prints only when the verbose mode is activated
func vprintf(s string, args ...interface{}) {
	if flagVerbose {
		fmt.Printf(s, args...)
	}
}

// fatal exits the process printing a formatted message to stderr
func fatal(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}

// getOutput runs the specified command with the specified arguments,
// and acquires its output (stdout). The output is returned as a list
// of strings, one per line.
func getOutput(command string, args ...string) ([]string, error) {
	if flagVerbose {
		fmt.Println("launching:", command, args)
	}

	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(out), "\n"), nil
}

// run runs the specified command with the specified arguments. stdout/stderr
// is shown only if verbose.
func run(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	if flagVerbose {
		fmt.Println("launching:", command, args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// spawn runs the specified command with the specified arguments, always showing
// stdout/stderr (attaching it to the parent console). If the command exits with
// an error, the parent process is exited as well with the same error.
func spawn(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if flagVerbose {
		fmt.Println("launching:", command, args)
	}

	err := cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok {
		os.Exit(ee.ExitCode())
	}
}

// mustOutput is like getOutput, but aborts the process with fatal if the
// command fails its execution.
func mustOutput(command string, args ...string) []string {
	out, err := getOutput(command, args...)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok && ee.Stderr != nil {
			fmt.Fprintf(os.Stderr, "%s", ee.Stderr)
		}
		fatal("error running command: %v", err)
	}
	return out
}

// mustRun is like run, but aborts the process with fatal if the command
// fails its execution.
func mustRun(command string, args ...string) {
	if err := run(command, args...); err != nil {
		fatal("error running command: %v", err)
	}
}

// findGitRoot checks whether path is part of a git repo; if so, returns its
// root, otherwise returns path itself.
func findGitRoot(path string) string {
	// FIXME: we could rewrite this without relying on git in PATH for Windows
	// users, and simply go through parent folders looking for ".git" directory.
	rootdir1, err := getOutput("git", "rev-parse", "--show-toplevel")
	if err == nil {
		path, err = filepath.Abs(rootdir1[0])
		if err != nil {
			fatal("error getting absolute path: %v", err)
		}
	}
	return path
}
