package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

var RestoreConsoleMode = func() {}

// vprintf is a printf that prints only when the verbose mode is activated
func vprintf(s string, args ...interface{}) {
	if flagVerbose {
		fmt.Printf(s, args...)
	}
}

func progress(s string, args ...interface{}) {
	text := color.Green.Render(fmt.Sprintf(s, args...))
	fmt.Fprint(os.Stdout, text)
}

func critical(s string, args ...interface{}) {
	text := color.Red.Render(fmt.Sprintf(s, args...))
	fmt.Fprint(os.Stderr, text)
}

// fatal exits the process printing a formatted message to stderr
func fatal(s string, args ...interface{}) {
	critical(s, args...)
	os.Exit(1)
}

func fatal_exitproc(err error, command string, args []string) {
	critical("error running command: ")
	fmt.Fprintf(os.Stderr, "%s %v\n", command, args)
	if ee, ok := err.(*exec.ExitError); ok && ee.Stderr != nil {
		critical("%s", ee.Stderr)
	}
	fatal("%v\n", err)
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
	cmd.Stderr = os.Stderr

	if flagVerbose {
		fmt.Println("launching:", command, args)
		cmd.Stdout = os.Stdout
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
	if runtime.GOOS == "windows" {
		if command == "git" {
			// Workaround for git for windows bug
			RestoreConsoleMode()
		}
	}

	if err != nil {
		fatal_exitproc(err, command, args)
	}

}

// mustOutput is like getOutput, but aborts the process with fatal if the
// command fails its execution.
func mustOutput(command string, args ...string) []string {
	out, err := getOutput(command, args...)
	if err != nil {
		fatal_exitproc(err, command, args)
	}
	return out
}

// mustRun is like run, but aborts the process with fatal if the command
// fails its execution.
func mustRun(command string, args ...string) {
	if err := run(command, args...); err != nil {
		fatal_exitproc(err, command, args)
	}
}

// isDir returns true if the specified path is a file
func isFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// isDir returns true if the specified path is a directory
func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

var (
	cachedGitRoot     string
	cachedGitRootOnce bool

	cachedLibdragonPath          string
	cachedLibdragonPathSubmodule bool
	cachedLibdragonPathOnce      bool
)

// findGitRoot checks whether the current directory is part of a git repo; if so,
// returns its root, otherwise returns empty string.
func findGitRoot() string {
	if !cachedGitRootOnce {
		cachedGitRootOnce = true

		rootdir1, err := getOutput("git", "rev-parse", "--show-toplevel")
		if err == nil {
			cachedGitRoot, err = filepath.Abs(rootdir1[0])
			if err != nil {
				fatal("error getting absolute path: %v", err)
			}
		}
	}

	return cachedGitRoot
}

// mustFindGitRoot is like findGitRoot, but aborts with fatal if no git repository
// is found.
func mustFindGitRoot() string {
	path := findGitRoot()
	if path == "" {
		fatal("error: this command must be run from a git repository\n")
	}
	return path
}

// findLibdragon searches for the libdragon vendored directory in the git repo
// which contains teh current directory. In case of success, it returns the path
// of the directory within the repo, and a boolean indicating whether the
// vendoring is being done via submodules (the alternative being subtrees).
func findLibdragon() (string, bool) {
	if !cachedLibdragonPathOnce {
		cachedLibdragonPathOnce = true

		repoRoot := mustFindGitRoot()

		// Check if we're using submodules
		if path, err := getOutput("git", "config",
			"--file", filepath.Join(repoRoot, ".gitmodules"),
			"--get", "submodule."+LIBDRAGON_SUBMODULE+".path"); err == nil && path[0] != "" {

			cachedLibdragonPath = path[0]
			cachedLibdragonPathSubmodule = true
		} else {
			// If we are using subtree, grep the logs to find the path
			logs := mustOutput("git", "log", "--grep", "git-subtree-dir:", "--format=tformat:%b")
			for _, logline := range logs {
				if strings.HasPrefix(logline, "git-subtree-dir:") && strings.HasSuffix(logline, "libdragon") {
					fields := strings.SplitN(logline, ":", 2)

					cachedLibdragonPath = strings.TrimSpace(fields[1])
					cachedLibdragonPathSubmodule = false
					break
				}
			}
		}
	}

	return cachedLibdragonPath, cachedLibdragonPathSubmodule
}

func findDockerImage() string {
	libdragonPath, _ := findLibdragon()
	if libdragonPath != "" {
		// Check if there's a reference to the needed toolchain in libdragon
		if imagebytes, err := os.ReadFile(filepath.Join(libdragonPath, "tools", ".docker-toolchain")); err == nil {
			return strings.TrimSpace(string(imagebytes))
		}
	}

	repoRoot := findGitRoot()
	if repoRoot != "" {
		if imagebytes, err := os.ReadFile(filepath.Join(repoRoot, ".git", CACHED_IMAGE_FILE)); err == nil {
			return strings.TrimSpace(string(imagebytes))
		}
	}

	return DOCKER_IMAGE
}
