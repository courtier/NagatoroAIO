package proxy

import (
	"bufio"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/apoorvam/goterminal"
	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/ui"
	"github.com/courtier/NagatoroAIO/utils"
)

var (
	sourcesTotal    []string
	saveFileName    string
	threads         int
	timeout         int
	writer          *goterminal.Writer
	receivedProxies map[string]bool
	ipPortRegex     *regexp.Regexp
	ipRegex         *regexp.Regexp
	portRegex       *regexp.Regexp
	cleanRegex      *regexp.Regexp
)

func askForStartScraping() {
	fileName := utils.AskForSpecificFileName("source")
	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}
	sourcesTotal, err = utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error occurred while loading sources", err.Error())
	}
	saveFileName = utils.CreateFileName("proxies", "scrapedproxies")
	threadsString := utils.AskForAnything("threads")
	threads, err = strconv.Atoi(threadsString)
	if err != nil {
		logger.LogFatal("enter a number")
	}
	timeoutString := utils.AskForAnything("timeout (seconds)")
	timeout, err = strconv.Atoi(timeoutString)
	if err != nil {
		logger.LogFatal("enter a number")
	}
	writer = ui.SetupProxyScrapeUI()
	ipPortRegex = regexp.MustCompile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$")
	ipRegex = regexp.MustCompile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
	portRegex = regexp.MustCompile("^([0-9]{1,5})$")
	cleanRegex = regexp.MustCompile("[^0-9.:\040\n]")
	receivedProxies = make(map[string]bool)
	startParse()
}

func startParse() {
	tr := &http.Transport{
		MaxIdleConns:       threads,
		IdleConnTimeout:    time.Duration(timeout) * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	scrapeProxies(client)
}

func scrapeProxies(client *http.Client) {
	logger.LogAction("press enter to start scraping proxies")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	logger.LogInfo("press s and enter to stop")

	sourcesAmount := len(sourcesTotal)

	if sourcesAmount == 0 {
		logger.LogFatal("no sources found")
	}

	routineSplitSources := utils.SplitSlice(sourcesTotal, threads)

	ch := make(chan string, len(routineSplitSources))

	logger.LogDebug("scraping", strconv.Itoa(sourcesAmount), "sources")

	for _, sourcesList := range routineSplitSources {
		go poolSources(client, ch, sourcesList)
	}

	totalProxyCounter, validProxyCounter, doneSourceCounter := 0, 0, 0
	batchProxyCounter := 0

	batchProxies := []string{}

	go func() {
		option, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if option == "s\n" {
			ui.StopGeneralUI(writer)
			logger.LogInfo("stopping")

			if len(batchProxies) != 0 {
				err := utils.SaveArrayOfString(saveFileName, batchProxies, false)
				if err != nil {
					logger.LogError("error while saving to file")
				}
			}

			logger.LogSuccess("done, scraped", strconv.Itoa(totalProxyCounter), "proxies")
			os.Exit(0)
		}
	}()

	for newProxy := range ch {
		if newProxy == "done" {
			doneSourceCounter++
			ui.UpdateProxyScrapeUI(writer, sourcesAmount, doneSourceCounter, totalProxyCounter, validProxyCounter)
			if sourcesAmount == doneSourceCounter {
				close(ch)
				break
			}
			continue
		}
		totalProxyCounter++
		if receivedProxies[newProxy] != true {
			batchProxyCounter++
			validProxyCounter++
			batchProxies = append(batchProxies, newProxy)
			receivedProxies[newProxy] = true
		}
		if batchProxyCounter == 10 {
			err := utils.SaveArrayOfString(saveFileName, batchProxies, false)
			if err != nil {
				logger.LogError("error while saving to file")
			}
			ui.UpdateProxyScrapeUI(writer, sourcesAmount, doneSourceCounter, totalProxyCounter, validProxyCounter)
			batchProxyCounter = 0
			batchProxies = []string{}
		}
	}

	ui.StopProxyScrapeUI(writer)

	if len(batchProxies) != 0 {
		err := utils.SaveArrayOfString(saveFileName, batchProxies, false)
		if err != nil {
			logger.LogError("error while saving to file")
		}
	}

	logger.LogSuccess("done, scraped", strconv.Itoa(totalProxyCounter), "proxies")

}

func poolSources(client *http.Client, ch chan string, sourceList []string) {
	for _, source := range sourceList {
		scrapeProxyRequest(client, ch, source)
		ch <- "done"
	}
}
