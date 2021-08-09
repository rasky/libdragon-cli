package cmd

import (
	"os"

	"golang.org/x/sys/windows"
)

var consoleMode uint32

func init() {
	fd := os.Stderr.Fd()
	windows.GetConsoleMode(windows.Handle(fd), &consoleMode)
}

// RestoreConsoleMode restore the console to the same mode it had at process startup.
// This is used to workaround a bug in git for Windows that leaves the console
// in a different state (in particular, it disables the EnableVirtualTerminalProcessingMode=0x4
// flag, which in turns breaks ANSI colors).
// Reference: https://github.com/git-for-windows/git/issues/2661
func RestoreConsoleMode() {
	fd := os.Stderr.Fd()
	windows.SetConsoleMode(windows.Handle(fd), consoleMode)
}
