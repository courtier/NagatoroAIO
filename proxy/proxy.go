package proxy

import (
	"bufio"
	"os"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
)

//AskForProxy choose from proxy modules
func AskForProxy() {
	logger.LogCategory("modules (http/s only)")
	logger.LogOption(1, "scrape proxies")
	logger.LogOption(2, "check proxies")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askForStartScraping()
	} else if option == "2" {
		askForStartChecking()
	} else {
		logger.LogError("coming soon))")
	}
}
