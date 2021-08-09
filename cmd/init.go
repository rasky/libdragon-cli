package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	flagInitForce         bool
	flagInitUseSubmodules bool
)

//go:embed prj-skeleton
var skeleton embed.FS

func doInit(cmd *cobra.Command, args []string) error {

	// Check if we're inside a git repository
	rootdir, err := getOutput("git", "rev-parse", "--show-toplevel")
	if err != nil {
		fatal("error: this command must be run within a git repository")
	}

	// Extract the project skeleton
	progress("Create project skeleton...\n")
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
		fatal("%v", err)
	}

	// Create a submodule for libdragon
	progress("Download libdragon...\n")

	if flagInitUseSubmodules {
		spawn("git", "submodule", "add", "--force",
			"--name", LIBDRAGON_SUBMODULE,
			"--branch", LIBDRAGON_BRANCH,
			LIBDRAGON_GIT)
	} else {
		// Reconstruct relative path in repo wrt the current directory, so that
		// we will be able to tell git subtree where to create the subtree folder.
		prefix := "libdragon"
		abspwd, err := filepath.Abs(".")
		if err == nil {
			reldir, err := filepath.Rel(rootdir[0], abspwd)
			if err == nil {
				prefix = path.Join(filepath.ToSlash(reldir), prefix)
			}
		}

		// git subtree does not work on empty repository (one with zero commits). The error
		// message is obscure. Since this is a very common case with "libdragon init",
		// verify whether HEAD exists and if it doesn't, create an initial empty commit.
		if _, err := getOutput("git", "rev-parse", "HEAD"); err != nil {
			mustRun("git", "commit", "--allow-empty", "-n", "-m", "Initial commit.")
		}

		// Add the subtree
		spawn("git", "-C", rootdir[0], "subtree", "add", "--prefix", prefix, LIBDRAGON_GIT, LIBDRAGON_BRANCH, "--squash")
	}

	updateToolchain("libdragon")

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
	cmdInit.Flags().BoolVarP(&flagInitUseSubmodules, "submodule", "m", false, "to vendor libdragon, use git submodule instead of git subtree")
	rootCmd.AddCommand(cmdInit)
}
