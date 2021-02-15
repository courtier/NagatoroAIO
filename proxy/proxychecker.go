package proxy

import (
	"bufio"
	"os"
	"strconv"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/ui"
	"github.com/courtier/NagatoroAIO/utils"
)

var (
	proxiesTotal []string
	protocol     string
)

func askForStartChecking() {
	fileName := utils.AskForSpecificFileName("proxy")
	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}
	proxiesTotal, err = utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error occurred while loading proxies", err.Error())
	}
	saveFileName = utils.CreateFileName("proxies", "checkedproxies")
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
	protocol = utils.AskForAnything("protocol (http(s)/socks4/socks5)")
	if !utils.IsProxyProtocolValid(protocol) {
		logger.LogFatal("https(s)/socks4/socks4a/socks5")
	}

	checkRequest, err = buildCheckRequest()
	if err != nil {
		logger.LogFatal("error while building check request", err.Error())
	}

	writer = ui.SetupProxyScrapeUI()
	checkProxies()
}

func checkProxies() {
	logger.LogAction("press enter to start checking proxies")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	logger.LogInfo("press s and enter to stop")

	proxyAmount := len(proxiesTotal)

	if proxyAmount == 0 {
		logger.LogFatal("no proxies found")
	}

	routineSplitProxies := utils.SplitSlice(proxiesTotal, threads)

	ch := make(chan string, len(routineSplitProxies))

	logger.LogDebug("checking", strconv.Itoa(proxyAmount), "proxies")

	for _, proxyList := range routineSplitProxies {
		go poolProxies(ch, proxyList)
	}

	workingProxyCounter, brokenProxyCounter, totalCheckedCounter := 0, 0, 0

	go func() {
		option, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if option == "s\n" {
			ui.StopGeneralUI(writer)
			logger.LogInfo("stopping")

			logger.LogSuccess("done, checked", strconv.Itoa(totalCheckedCounter), "proxies")
			os.Exit(0)
		}
	}()

	for newProxy := range ch {
		logger.LogDebug("new proxy", newProxy)
		totalCheckedCounter++
		if newProxy == "broken" {
			brokenProxyCounter++
		} else {
			workingProxyCounter++
			err := utils.AppendString(saveFileName, newProxy, false)
			if err != nil {
				logger.LogError("error while saving to file")
			}
		}
		ui.UpdateProxyCheckUI(writer, proxyAmount, totalCheckedCounter, workingProxyCounter, brokenProxyCounter)
		if totalCheckedCounter == proxyAmount {
			close(ch)
			break
		}
	}

	ui.StopGeneralUI(writer)

	logger.LogSuccess("done, checked", strconv.Itoa(totalCheckedCounter), "proxies")

}

func poolProxies(ch chan string, proxyList []string) {
	for _, proxy := range proxyList {
		checkProxyRequest(ch, proxy)
	}
}
