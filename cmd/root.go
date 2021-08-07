package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	DOCKER_IMAGE     = "anacierdem/libdragon"
	CONTAINER_FILE   = "libdragon-docker-container"
	VOLUME_ROOT      = "/app"
	LIBDRAGON_GIT    = "https://github.com/DragonMinded/libdragon"
	LIBDRAGON_BRANCH = "trunk"
)

var (
	flagVerbose bool
	flagChdir   string
)

var rootCmd = &cobra.Command{
	Use:   "libdragon",
	Short: "libdragon command line tool",
	Long:  "libdragon tool - help managing development of Nintendo 64 ROMs using libdragon",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "be verbose")
	rootCmd.PersistentFlags().StringVarP(&flagChdir, "chdir", "C", "", "work in the specified directory")

	cobra.OnInitialize(func() {
		if flagChdir != "" {
			if err := os.Chdir(flagChdir); err != nil {
				fatal("%v\n", err)
			}
			vprintf("chdir to: %s\n", flagChdir)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
