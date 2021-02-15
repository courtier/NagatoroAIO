package logger

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gookit/color"
)

var (
	red     = color.FgRed.Render
	green   = color.FgGreen.Render
	blue    = color.FgBlue.Render
	cyan    = color.FgCyan.Render
	magenta = color.FgMagenta.Render
	yellow  = color.FgYellow.Render
)

var (
	//DebugMode enable debug mode
	DebugMode bool
)

//LogInfo logs info
func LogInfo(input string, args ...string) {
	msg := cyan("info >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogDebug logs debug
func LogDebug(input string, args ...string) {
	if !DebugMode {
		return
	}
	msg := blue("debug >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogWarning logs warning
func LogWarning(input string, args ...string) {
	msg := yellow("warning >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogInput logs input
func LogInput(input string, args ...string) {
	msg := magenta("input >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	msg += magenta(" >> ")
	fmt.Print(msg)
}

//LogAction logs action
func LogAction(input string, args ...string) {
	msg := magenta("action >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogError logs error
func LogError(input string, args ...string) {
	msg := red("error >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogSuccess logs success
func LogSuccess(input string, args ...string) {
	msg := green("success >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogCategory logs category
func LogCategory(input string, args ...string) {
	msg := blue("category >> ") + input
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogOption logs option
func LogOption(id int, args ...string) {
	msg := cyan(strconv.Itoa(id) + " >>")
	for _, arg := range args {
		msg += " " + arg
	}
	fmt.Println(msg)
}

//LogFatal throw fatal error and quit with (1)
func LogFatal(input string, args ...string) {
	msg := "fatal >> " + input
	for _, arg := range args {
		msg += " " + arg
	}
	color.Error.Println(msg)
	os.Exit(1)
}

//LogCustom log with custom colors use html
func LogCustom(input string, args ...string) {
	msg := input
	for _, arg := range args {
		msg += " " + arg
	}
	color.Println(msg)
}
