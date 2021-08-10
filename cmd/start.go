package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// searchContainer searches for a libdragon container associated to a certain
// path, and optionally autostarts it if not found. Returns the container ID.
func searchContainer(path string, autostart bool) string {
	// Check whether there is a container file. This is the fastest and safest
	// way to find the container associated with this directory, but only works
	// for directories which are git roots.
	// Otherwise, if we are requested to mount in a directory which is not a repo
	// root, there will be no container file to speed up further usage later.
	containerFile := filepath.Join(path, ".git", CACHED_CONTAINER_FILE)

	if containerBytes, err := os.ReadFile(containerFile); err == nil {
		container := strings.TrimSpace(string(containerBytes))

		// We want to check whether the container still exists and optionally
		// start it with the smallest possible amount of docker commands, so
		// that execution is as fast as possible (for libdragon make).
		// In fact, the file might be stale (eg: the container has been purged).
		if autostart {
			// In case of autostart, we just run "docker start". If that fails,
			// we assume that the container does not exist anymore (stale file).
			if err := run("docker", "container", "start", container); err == nil {
				vprintf("container found: %v\n", container)
				return container
			}
		} else {
			// If no autostart was requested, we use "docker container ls" to
			// check whether it still exists.
			if out := mustOutput("docker", "container", "ls", "-qa", "-f", "id="+container); out[0] != "" {
				vprintf("container found: %v\n", out[0])
				return out[0]
			}
		}
	}

	// Fallback: look for containers that mount the same volume
	out := mustOutput("docker", "container", "ls", "-qa", "-f", "volume="+string(path))
	if out[0] != "" {
		vprintf("container found: %v\n", out[0])
		if autostart {
			// Restart the container in case it's not running.
			mustRun("docker", "container", "start", out[0])
		}
		return out[0]
	}

	// No container found. If we were not asked to autostart, there's nothing
	// more to do.
	if !autostart {
		return ""
	}

	// Start a new container
	out = mustOutput("docker", "run",
		"-e", "IS_DOCKER=true",
		"-d",                                                       // detached
		"--mount", "type=bind,source="+path+",target="+VOLUME_ROOT, // mount
		"-w", "/app", // working dir
		DOCKER_IMAGE,
		"tail", "-f", "/dev/null",
	)

	// Try writing the container file. This will fail if this is not a git
	// root because the .git subdir will not exist. Ignore the error.
	os.WriteFile(containerFile, []byte(out[0]+"\n"), 0666)

	return out[0]
}

func doStart(cmd *cobra.Command, args []string) error {
	path := findGitRoot()
	if path == "" {
		// If we're not in a git repository, create a container that
		// mounts the current directory.
		path = "."
	}
	searchContainer(path, true)
	return nil
}

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start a libdragon container for the current repository.",
	Example: `  libdragon start
	-- start the libdragon container`,
	RunE:         doStart,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(cmdStart)
}
