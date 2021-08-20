package cmd

import (
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

func doMake(cmd *cobra.Command, args []string) error {
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
		"make"}
	docker_args = append(docker_args, args...)

	spawn("docker", docker_args...)
	return nil
}

var cmdMake = &cobra.Command{
	Use:   "make",
	Short: "Run the libdragon build system.",
	Example: `  libdragon make
	-- build the current application`,
	RunE:         doMake,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(cmdMake)
}
