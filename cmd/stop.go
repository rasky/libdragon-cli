package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func doStop(cmd *cobra.Command, args []string) error {
	path := findGitRoot(".")
	out := searchContainer(path, false)
	if out != "" {
		mustRun("docker", "container", "rm", "--force", out)

		// Remove the container file if it exists
		os.Remove(filepath.Join(path, ".git", CONTAINER_FILE))
	}
	return nil
}

var cmdStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop a libdragon container for the current repository.",
	Example: `  libdragon stop
	-- stop the libdragon container`,
	RunE:         doStop,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(cmdStop)
}
