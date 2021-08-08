package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagUpdateLibdragonPath string
	flagUpdateRev           string
)

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func findLibdragon(repoRoot string) (string, bool) {
	// Check if we're using submodules
	if path, err := getOutput("git", "config",
		"--file", filepath.Join(repoRoot, ".gitmodules"),
		"--get", "submodule."+LIBDRAGON_SUBMODULE+".path"); err == nil && path[0] != "" {
		return path[0], true
	}

	// If we are using subtree, grep the logs to find the path
	logs := mustOutput("git", "log", "--grep", "git-subtree-dir:", "--format=tformat:%b")
	for _, logline := range logs {
		if strings.HasPrefix(logline, "git-subtree-dir:") && strings.HasSuffix(logline, "libdragon") {
			fields := strings.SplitN(logline, ":", 2)
			return strings.TrimSpace(fields[1]), false
		}
	}

	// Not found
	return "", false
}

func updateToolchain(libdragonPath string) {
	// Check if there's a reference to the needed toolchain in libdragon, otherwise
	// assume the latest is fine (hopefully...)
	image := DOCKER_IMAGE
	if imagebytes, err := os.ReadFile(filepath.Join(libdragonPath, "tools", ".docker-toolchain")); err == nil {
		image = strings.TrimSpace(string(imagebytes))
	}

	// Pull the image
	color.Green("\nUpdating toolchain...\n")
	spawn("docker", "pull", image)
}

func doUpdate(cmd *cobra.Command, args []string) error {
	repoRoot := findGitRoot(".")

	libdragonPath := flagUpdateLibdragonPath
	useSubmodules := false
	if libdragonPath == "" {
		libdragonPath, useSubmodules = findLibdragon(repoRoot)
		if libdragonPath == "" {
			fatal("cannot find libdragon in this repository\nuse --directory to specify the location")
		}
	} else {
		if !isDir(libdragonPath) {
			fatal("%s: not a directory", libdragonPath)
		}
	}

	vprintf("found libdragon: %v ", filepath.Join(repoRoot, libdragonPath))
	if useSubmodules {
		vprintf("(submodule)\n")
	} else {
		vprintf("(subtree)\n")
	}

	// Update libdragon
	color.Green("Updating libdragon...\n")
	if useSubmodules {
		spawn("git", "submodule", "update", "--remote", "--merge", libdragonPath)
	} else {
		spawn("git", "subtree", "pull", "--prefix", libdragonPath, LIBDRAGON_GIT, LIBDRAGON_BRANCH, "--squash")
	}

	updateToolchain(libdragonPath)

	return nil
}

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update libdragon in the current repository, and its toolchain.",
	Example: `  libdragon update
	-- update libdragon`,
	RunE:         doUpdate,
	SilenceUsage: true,
}

func init() {
	cmdUpdate.Flags().StringVarP(&flagUpdateLibdragonPath, "directory", "d", "", "specify where libdragon is located (default: autodetect)")
	cmdUpdate.Flags().StringVarP(&flagUpdateRev, "revision", "r", "", "specify the commit/branch/tag to update libdragon to (default: update current branch)")
	rootCmd.AddCommand(cmdUpdate)
}
