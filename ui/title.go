package ui

import (
	"fmt"
	"os/exec"
	"runtime"
)

//https://docs.microsoft.com/en-us/windows/terminal/tutorials/tab-title

//SetTitle crossplatform window title setter
func SetTitle(title string) {
	if runtime.GOOS == "windows" {
		exec.Command("TITLE \"" + title + "\"").Run()
		exec.Command("$Host.UI.RawUI.WindowTitle = \"" + title + "\"").Run()
	} else {
		fmt.Print("\033]0;" + title + "\a")
	}
}
