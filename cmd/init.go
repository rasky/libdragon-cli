package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

var flagInitForce bool

//go:embed prj-skeleton
var skeleton embed.FS

func doInit(cmd *cobra.Command, args []string) error {

	// Check if we're inside a git repository
	_, err := getOutput("git", "rev-parse", "--show-toplevel")
	if err != nil {
		fatal("error: this command must be run within a git repository")
	}

	// Create a submodule for libdragon
	mustRun("git", "submodule", "add", "--force", "https://github.com/DragonMinded/libdragon")

	// Extract the project skeleton
	skfs, _ := fs.Sub(skeleton, "prj-skeleton")
	err = fs.WalkDir(skfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err == nil {
			if d.IsDir() {
				vprintf("creating: %s\n", path)
				os.Mkdir(path, 0777)
			} else {
				if !flagInitForce {
					if _, err := os.Stat(path); err == nil {
						return fmt.Errorf("file already exists: %v (use --force to overwrite)", path)
					}
				}
				vprintf("extracting: %s\n", path)
				data, _ := fs.ReadFile(skfs, path)
				err = os.WriteFile(path, data, 0666)
			}
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Create a skeleton libdragon application in the current directory.",
	Example: `  libdragon init
	-- create skeleton project`,
	RunE:         doInit,
	SilenceUsage: true,
}

func init() {
	cmdInit.Flags().BoolVarP(&flagInitForce, "force", "f", false, "force overwriting")
	rootCmd.AddCommand(cmdInit)
}
