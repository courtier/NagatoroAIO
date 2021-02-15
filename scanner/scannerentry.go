package scanner

import (
	"bufio"
	"os"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
)

//AskForScanner choose from scanner modules
func AskForScanner() {
	logger.LogCategory("modules")
	logger.LogOption(1, "scan urls for injections")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		askForStartScanning()
	} else {
		logger.LogError("coming soon))")
	}
}
