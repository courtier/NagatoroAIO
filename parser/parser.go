package parser

import (
	"bufio"
	"os"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
)

//AskForParser choose from parser modules
func AskForParser() {
	logger.LogCategory("modules")
	logger.LogOption(1, "parse dorks")
	logger.LogCustom("<red>2 >> filter dorks by amount of search results - soon</>")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askForEngine()
	} else {
		logger.LogDebug("coming soon))")
	}
}

func askForEngine() {
	logger.LogCategory("engines")
	logger.LogOption(1, "startpage")
	logger.LogOption(2, "bing")
	logger.LogOption(3, "okeano")
	logger.LogOption(4, "yahoo")
	logger.LogOption(5, "aol")
	logger.LogOption(6, "google")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askForStartParsing("startpage", parseStartPageRequest)
	} else if option == "2" {
		askForStartParsing("bing", parseBingRequest)
	} else if option == "3" {
		askForStartParsing("okeano", parseOkeanoRequest)
	} else if option == "4" {
		askForStartParsing("yahoo", parseYahooRequest)
	} else if option == "5" {
		askForStartParsing("aol", parseAolRequest)
	} else if option == "6" {
		askForStartParsing("google", parseGoogleRequest)
	} else {
		logger.LogDebug("coming soon))")
	}
}
