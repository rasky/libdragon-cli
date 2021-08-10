package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	flagUpdateDockerImage   string
	flagUpdateLibdragonPath string
	flagUpdateWhat          string
)

// updateToolchain updates the docker toolchain image that will be used to compile
// libdragon. It is a wrapper over "docker pull".
func updateToolchain() {
	var image string

	if flagUpdateDockerImage != "" {
		image = flagUpdateDockerImage

		// Persist the image name by creating a local file in the repo root.
		// Users might want to commit this file to persist their custom
		// toolchain selection
		repoRoot := findGitRoot()
		if repoRoot != "" {
			repoRoot = "."
		}

		if err := os.WriteFile(filepath.Join(repoRoot, CACHED_IMAGE_FILE), []byte(image+"\n"), 0666); err != nil {
			fatal("error persisting toolchain change: %v\n", err)
		}
	} else {
		image = findDockerImage()
	}

	// Pull the requested image
	spawn("docker", "pull", image)

}

// updateLibdragon updates the vendored libdragon copy within the repository.
// It works with both the submodule and subtree vendoring strategy.
func updateLibdragon() {
	repoRoot := mustFindGitRoot()

	// Repo-relative path where libdragon is vendored
	var libdragonPath string
	var useSubmodules bool

	if flagUpdateLibdragonPath == "" {
		// If the directory was not specified on a command line, search for it.
		libdragonPath, useSubmodules = findLibdragon()
		if libdragonPath == "" {
			fatal("cannot find libdragon in this repository\nuse --directory to specify the location\n")
		}
	} else {
		// If the directory was specified on the command line, we assume that
		// it's a cwd-relative path (so a path that makes sense for the user
		// in the context of where the command was launched). Convert it to
		// absolute (if not already).
		var err error
		libdragonPath, err = filepath.Abs(flagUpdateLibdragonPath)
		if err != nil {
			fatal("error convert directory to absolute path: %v\n", err)
		}
		// Check if it's really an existing directory
		if !isDir(libdragonPath) {
			fatal("%s: not a directory\n", libdragonPath)
		}

		// Try to detect the vendoring strategy. Look for an existing .git
		// file that would hint at submodule.
		useSubmodules = isFile(filepath.Join(libdragonPath, ".git"))

		// Now convert to repo-relative path.
		libdragonPath, err = filepath.Rel(repoRoot, libdragonPath)
		if err != nil {
			fatal("error convert directory to repo-relative path: %v\n", err)
		}
	}

	vprintf("found libdragon: %v ", filepath.Join(repoRoot, libdragonPath))
	if useSubmodules {
		vprintf("(submodule)\n")
	} else {
		vprintf("(subtree)\n")
	}

	// Update libdragon
	if useSubmodules {
		spawn("git", "submodule", "update", "--remote", "--merge", filepath.Join(repoRoot, libdragonPath))
	} else {
		spawn("git", "subtree", "pull", "--prefix", libdragonPath, LIBDRAGON_GIT, LIBDRAGON_BRANCH, "--squash")
	}

}

func doUpdate(cmd *cobra.Command, args []string) error {
	what := "libdragon\000toolchain"
	if len(args) > 0 {
		what = strings.Join(args, "\000")
	}

	// First update libdragon, as that might change the required toolchain.
	if strings.Contains(what, "libdragon") {
		color.Greenp("Updating libdragon...\n")
		updateLibdragon()
	}

	// Update the toolchain, if requested
	if strings.Contains(what, "toolchain") {
		color.Greenp("Updating toolchain...\n")
		updateToolchain()
	}

	return nil
}

var cmdUpdate = &cobra.Command{
	Use:   "update [what]",
	Short: "Update libdragon in the current repository, and its toolchain.",
	Long: `This command can be used to perform updates on the vendored libdragon version, 
the toolchain, or both. By default, it updates everything. Otherwise, you can specify
as argument a single item to update, which can be either "libdragon" or "docker"`,
	Example: `  libdragon update
	-- update libdragon and toolchain
  libdragon update toolchain
      -- update toolchain only`,
	ValidArgs:    []string{"libdragon", "toolchain"},
	Args:         cobra.OnlyValidArgs,
	RunE:         doUpdate,
	SilenceUsage: true,
}

func init() {
	cmdUpdate.Flags().StringVarP(&flagUpdateLibdragonPath, "directory", "d", "", "specify where libdragon is located (default: autodetect)")
	cmdUpdate.Flags().StringVarP(&flagUpdateDockerImage, "image", "i", "", "specify the Docker image to use as a toolchain")
	rootCmd.AddCommand(cmdUpdate)
}
