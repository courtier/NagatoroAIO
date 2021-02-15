package dorks

import (
	"bufio"
	"os"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
)

//AskForDorks choose from dork modules
func AskForDorks() {
	logger.LogCategory("modules")
	logger.LogOption(1, "generate dorks from templates")
	logger.LogOption(2, "generate DE, PT and PFs from urls")
	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askStartDorkGenerator()
	} else if option == "2" {
		askForURLToParamsConversion()
	} else {
		logger.LogError("nah")
	}
}
