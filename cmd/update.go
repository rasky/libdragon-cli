package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	flagUpdateLibdragonPath string
)

// findLibdragon searches for the libdragon vendored directory in the specified
// git repository. In case of success, it returns the path of the directory within
// the repo, and a boolean indicating whether the vendoring is being done via
// submodules (the alternative being subtrees).
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

// updateToolchain updates the docker toolchain image that will be used to compile
// libdragon. It is a wrapper over "docker pull".
func updateToolchain(libdragonPath string) {
	// Check if there's a reference to the needed toolchain in libdragon, otherwise
	// assume the latest is fine (hopefully...)
	image := DOCKER_IMAGE
	if imagebytes, err := os.ReadFile(filepath.Join(libdragonPath, "tools", ".docker-toolchain")); err == nil {
		image = strings.TrimSpace(string(imagebytes))
	}

	// Pull the image
	color.Greenp("\nUpdating toolchain...\n")
	spawn("docker", "pull", image)
}

func doUpdate(cmd *cobra.Command, args []string) error {
	repoRoot := findGitRoot(".")

	// Repo-relative path where libdragon is vendored
	var libdragonPath string
	var useSubmodules bool

	if flagUpdateLibdragonPath == "" {
		libdragonPath, useSubmodules = findLibdragon(repoRoot)
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
	color.Greenp("Updating libdragon...\n")
	if useSubmodules {
		spawn("git", "submodule", "update", "--remote", "--merge", filepath.Join(repoRoot, libdragonPath))
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
	rootCmd.AddCommand(cmdUpdate)
}
