package keywords

import (
	"bufio"
	"os"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/translator"
)

//AskForKeywords choose from keyword modules
func AskForKeywords() {
	logger.LogCategory("modules")
	logger.LogOption(1, "generate keywords")
	logger.LogOption(2, "translate keywords")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askForGenerator()
	} else if option == "2" {
		translator.AskForStartTranslating()
	} else {
		logger.LogError("coming soon))")
	}
}

func askForGenerator() {
	logger.LogCategory("generators")
	logger.LogOption(1, "startpage")
	logger.LogOption(2, "google")
	logger.LogOption(3, "qwant")
	logger.LogOption(4, "amazon")
	logger.LogOption(5, "yahoo")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askForStartScraping(doStartpageRequest)
	} else if option == "2" {
		askForStartScraping(doGoogleRequest)
	} else if option == "3" {
		askForStartScraping(doQwantRequest)
	} else if option == "4" {
		askForStartScraping(doAmazonRequest)
	} else if option == "5" {
		askForStartScraping(doYahooRequest)
	} else {
		logger.LogError("coming soon))")
	}
}
