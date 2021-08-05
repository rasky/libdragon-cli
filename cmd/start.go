package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// searchContainer searches for a libdragon container associated to a certain
// path, and optionally autostarts it if not found. Returns the container ID.
func searchContainer(path string, autostart bool) string {
	// Check whether there is a container file. This is the fastest and safest
	// way to find the container associated with this directory, but only works
	// for directories which are git roots.
	containerFile := filepath.Join(path, ".git", CONTAINER_FILE)

	container, err := os.ReadFile(containerFile)
	if err == nil {
		out := mustOutput("docker", "container", "ls", "-qa", "-f", "id="+string(container))
		if out[0] != "" {
			vprintf("container found: %v\n", out[0])
			if autostart {
				// Restart the container in case it's not running.
				mustRun("docker", "container", "start", out[0])
			}
			return out[0]
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
	os.WriteFile(containerFile, []byte(out[0]), 0666)

	return out[0]
}

func doStart(cmd *cobra.Command, args []string) error {
	path := findGitRoot(".")
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
