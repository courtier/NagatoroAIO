package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/courtier/NagatoroAIO/dorks"
	"github.com/courtier/NagatoroAIO/keywords"
	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/parser"
	"github.com/courtier/NagatoroAIO/proxy"
	"github.com/courtier/NagatoroAIO/scanner"
	"github.com/courtier/NagatoroAIO/textutils"
	"github.com/courtier/NagatoroAIO/ui"
	"github.com/courtier/NagatoroAIO/utils"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		logger.DebugMode = true
		logger.LogDebug("debug mode enabled")
	} else {
		logger.DebugMode = false
	}
	logger.LogCustom("<magenta>███╗░░██╗░█████╗░░██████╗░░█████╗░████████╗░█████╗░██████╗░░█████╗░\n" +
		"████╗░██║██╔══██╗██╔════╝░██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗\n" +
		"██╔██╗██║███████║██║░░██╗░███████║░░░██║░░░██║░░██║██████╔╝██║░░██║\n" +
		"██║╚████║██╔══██║██║░░╚██╗██╔══██║░░░██║░░░██║░░██║██╔══██╗██║░░██║\n" +
		"██║░╚███║██║░░██║╚██████╔╝██║░░██║░░░██║░░░╚█████╔╝██║░░██║╚█████╔╝\n" +
		"╚═╝░░╚══╝╚═╝░░╚═╝░╚═════╝░╚═╝░░╚═╝░░░╚═╝░░░░╚════╝░╚═╝░░╚═╝░╚════╝░</>")
	logger.LogCustom("<yellow>" + utils.GenRandomSentence() + "</>")
	ui.SetTitle("Nagatoro - courtier#5443")
	utils.LoadConfigToMemory()
	AskForOption()
}

//AskForOption choose from main modules
func AskForOption() {
	logger.LogCategory("modules")
	logger.LogOption(1, "utils")
	logger.LogOption(2, "keywords")
	logger.LogOption(3, "dorks")
	logger.LogOption(4, "parser")
	logger.LogOption(5, "scanner")
	logger.LogOption(6, "proxies")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		textutils.AskForUtils()
	} else if option == "2" {
		keywords.AskForKeywords()
	} else if option == "3" {
		dorks.AskForDorks()
	} else if option == "4" {
		parser.AskForParser()
	} else if option == "5" {
		scanner.AskForScanner()
	} else if option == "6" {
		proxy.AskForProxy()
	} else {
		logger.LogError("coming soon))")
	}
}
