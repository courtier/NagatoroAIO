package scanner

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/ui"
	"github.com/courtier/NagatoroAIO/utils"
)

var (
	urlsTotal       []string
	saveFileName    string
	threads         int
	client          *http.Client
	receivedDomains map[string]bool
	protocol        string
)

func askForStartScanning() {
	fileName := utils.AskForSpecificFileName("url")
	file, err := utils.OpenFile(fileName)
	defer file.Close()
	if err != nil {
		logger.LogFatal("error while opening file")
	}
	urlsTotal, err = utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error occurred while loading urls", err.Error())
	}
	if len(urlsTotal) == 0 {
		logger.LogFatal("no urls found")
	}
	saveFileName = utils.CreateFileName("urls", "scannedurls")
	threadsString := utils.AskForAnything("threads")
	threads, err = strconv.Atoi(threadsString)
	if err != nil {
		logger.LogFatal("enter a number")
	}
	timeout := time.Duration(int(utils.Config.Get("timeout").(int64))) * time.Second
	tr := &http.Transport{
		MaxIdleConns:       5,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{Timeout: timeout, Transport: tr}
	startScan()
}

func startScan() {
	logger.LogAction("press enter to start scanning urls")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	logger.LogInfo("press s and enter to stop")

	urlsAmount := len(urlsTotal)

	routineSplitURLS := utils.SplitSlice(urlsTotal, threads)

	ch := make(chan map[string]interface{}, len(routineSplitURLS))

	logger.LogInfo("scanning", strconv.Itoa(urlsAmount), "urls")

	for _, urlsList := range routineSplitURLS {
		go poolURLS(ch, urlsList)
	}

	injectableCounter, totalCounter := 0, 0
	stopProgram := false

	go func() {
		option, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if option == "s\n" {
			stopProgram = true
			logger.LogSuccess("stopped, found", strconv.Itoa(injectableCounter), "injectable urls")
			os.Exit(0)
		}
	}()

	for urlResult := range ch {
		totalCounter++
		if stopProgram == true || totalCounter == urlsAmount {
			close(ch)
			break
		}
		if urlResult["vulnerable"].(bool) {
			resultString := resultToString(urlResult)
			logger.LogSuccess(resultString)
			injectableCounter++
			err := utils.AppendString(saveFileName, resultString, false)
			if err != nil {
				logger.LogError("error while appending to file")
			}
			ui.SetTitle(fmt.Sprintf("Vulnerable: %d - Non-Vulnerable: %d", injectableCounter, (totalCounter - injectableCounter)))
		} else {
			logger.LogWarning(urlResult["url"].(string))
		}
	}

	logger.LogSuccess("done, found", strconv.Itoa(injectableCounter), "injectable urls")
}

func poolURLS(ch chan map[string]interface{}, urlsList []string) {
	for _, url := range urlsList {
		scanURL(url, client, ch)
	}
}

func resultToString(result map[string]interface{}) string {
	return result["url"].(string) + " - dbms: " + result["dbms"].(string) + " - payload: " + result["payload"].(string)
}
