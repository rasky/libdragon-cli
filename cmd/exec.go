package cmd

import (
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

func spawnDockerExec(args ...string) error {
	root := findGitRootOrCwd()
	container := searchContainer(root, true)

	// Reconstruct the relative path within the git root, so that it can be
	// set as working directory in the docker container. If this fails,
	// just avoid setting a working directory and hope for the best.
	workdir := VOLUME_ROOT
	if root != "." {
		abspwd, err := filepath.Abs(".")
		if err == nil {
			reldir, err := filepath.Rel(root, abspwd)
			if err == nil {
				workdir = path.Join(workdir, filepath.ToSlash(reldir))
			}
		}
	}

	docker_args := []string{
		"exec",
		"--workdir", workdir,
		container,
	}
	docker_args = append(docker_args, args...)

	spawn("docker", docker_args...)
	return nil
}

func doExec(cmd *cobra.Command, args []string) error {
	return spawnDockerExec(args...)
}

var cmdExec = &cobra.Command{
	Use:   "exec <command> [args...]",
	Short: "Run a command using the libdragon toolchain.",
	Example: `  libdragon exec makedfs game.dfs assets/
	-- build the DFS filesystem from the assets subdirectory`,
	Args:         cobra.MinimumNArgs(1),
	RunE:         doExec,
	SilenceUsage: true,
}

func init() {
	cmdExec.Flags().SetInterspersed(false)
	rootCmd.AddCommand(cmdExec)
}
