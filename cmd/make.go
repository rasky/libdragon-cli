package cmd

import (
	"github.com/spf13/cobra"
)

func doMake(cmd *cobra.Command, args []string) error {
	// "libdragon make" is just a shortcut for "libdragon exec make"
	args = append([]string{"make"}, args...)
	return spawnDockerExec(args...)
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
	cmdMake.Flags().SetInterspersed(false)
	rootCmd.AddCommand(cmdMake)
}
