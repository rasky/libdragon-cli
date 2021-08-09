package cmd

import (
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	DOCKER_IMAGE        = "anacierdem/libdragon"
	CONTAINER_FILE      = "libdragon-docker-container"
	VOLUME_ROOT         = "/app"
	LIBDRAGON_GIT       = "https://github.com/DragonMinded/libdragon"
	LIBDRAGON_BRANCH    = "trunk"
	LIBDRAGON_SUBMODULE = "libdragon"
)

var (
	flagVerbose  bool
	flagChdir    string
	flagColorize bool
)

var rootCmd = &cobra.Command{
	Use:   "libdragon",
	Short: "libdragon command line tool",
	Long:  "libdragon tool - help managing development of Nintendo 64 ROMs using libdragon",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "be verbose")
	rootCmd.PersistentFlags().StringVarP(&flagChdir, "chdir", "C", "", "work in the specified directory")
	rootCmd.PersistentFlags().BoolVarP(&flagColorize, "color", "", true, "use colorful output")

	cobra.OnInitialize(func() {
		if flagChdir != "" {
			if err := os.Chdir(flagChdir); err != nil {
				fatal("%v\n", err)
			}
			vprintf("chdir to: %s\n", flagChdir)
		}
		if !flagColorize {
			color.Disable()
		}
	})

	// Colorize usage prompt
	cobra.AddTemplateFunc("StyleHeading", color.Green.Sprint)
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
		`Flags:`, `{{StyleHeading "Flags:"}}`,
	).Replace(usageTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
