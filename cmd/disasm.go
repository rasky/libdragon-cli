package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	flagDisasmFile string
)

func doDisasm(cmd *cobra.Command, args []string) error {
	if flagDisasmFile == "" {
		root := findGitRootOrCwd()

		// Look for the file. We start from the current directory and look
		// for files with ".elf" extension in either the current directory or
		// a "build" subdirectory, traversing the tree up until git root (if any).
		cwd := "."
		var matches []string
		for i := 0; i < 10; i++ {
			matches, _ = filepath.Glob(filepath.Join(cwd, "*.elf"))
			if len(matches) > 0 {
				break
			}
			matches, _ = filepath.Glob(filepath.Join(cwd, "build", "*.elf"))
			if len(matches) > 0 {
				break
			}

			cwdAbs, _ := filepath.Abs(cwd)
			if cwdAbs == root || cwdAbs == "/" {
				break
			}
			cwd = filepath.Join(cwd, "..")
		}

		if len(matches) == 0 {
			fatal("cannot find ELF file to disassemble -- use --file to specify\n")
		}
		if len(matches) != 1 {
			fatal("multiple ELF files found in %s -- use --file to specify\n", cwd)
		}

		flagDisasmFile = matches[0]
	}

	dockerArgs := []string{
		"mips64-elf-objdump",
		"-S",
	}
	if len(args) != 0 {
		dockerArgs = append(dockerArgs, "--disassemble="+args[0])
	}
	dockerArgs = append(dockerArgs, flagDisasmFile)

	spawnDockerExec(dockerArgs...)
	return nil
}

var cmdDisasm = &cobra.Command{
	Use:   "disasm [symbol]",
	Short: "Disassemble a N64 binary and show assembly source for symbols.",
	Example: `  libdragon disasm dfs_read
	-- show the disassembled code for the "dfs_read" function`,
	Args:         cobra.MaximumNArgs(1),
	RunE:         doDisasm,
	SilenceUsage: true,
}

func init() {
	cmdDisasm.Flags().StringVarP(&flagDisasmFile, "file", "F", "", "ELF binary to disassemble (default: autodiscover)")
	rootCmd.AddCommand(cmdDisasm)
}
