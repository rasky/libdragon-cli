package cmd

import (
	"github.com/spf13/cobra"
)

func doMake(cmd *cobra.Command, args []string) error {
	path := findGitRoot(".")
	container := searchContainer(path, true)

	args = append([]string{
		"exec",
		container,
		"make"},
		args...)

	spawn("docker", args...)
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
