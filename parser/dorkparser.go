package parser

import (
	"bufio"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apoorvam/goterminal"
	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/proxy"
	"github.com/courtier/NagatoroAIO/ui"
	"github.com/courtier/NagatoroAIO/utils"
	"github.com/jpillora/go-tld"
)

type (
	queryFuncType func(*http.Client, chan string, string, int)
)

var (
	dorksTotal      []string
	proxiesTotal    []string
	dorkClients     map[string]*http.Client
	queryFunction   queryFuncType
	saveFileName    string
	threads         int
	engine          string
	writer          *goterminal.Writer
	pagePerDork     int
	pageFormatRegex *regexp.Regexp
	pageTypeRegex   *regexp.Regexp
	receivedDomains map[string]bool
	protocol        string
)

func askForStartParsing(searchEngine string, queryFunc queryFuncType) {
	fileName := utils.AskForSpecificFileName("dorks")
	file, err := utils.OpenFile(fileName)
	defer file.Close()
	if err != nil {
		logger.LogFatal("error while opening file")
	}
	dorksTotal, err = utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error occurred while loading dorks", err.Error())
	}
	if len(dorksTotal) == 0 {
		logger.LogFatal("no dorks found")
	}
	saveFileName = utils.CreateFileName("urls", "parsedurls")
	queryFunction = queryFunc
	threadsString := utils.AskForAnything("threads")
	threads, err = strconv.Atoi(threadsString)
	if err != nil {
		logger.LogFatal("enter a number")
	}
	protocol = utils.AskForAnything("proxy (http(s)/socks4(a)/socks5) - enter for no proxy")
	if protocol != "" && !utils.IsProxyProtocolValid(protocol) {
		logger.LogFatal("https(s)/socks4/socks4a/socks5")
	}
	if protocol != "" {
		fileName = utils.AskForSpecificFileName("proxy")
		file, err = utils.OpenFile(fileName)
		if err != nil {
			logger.LogFatal("error while opening file")
		}
		proxiesTotal, err = utils.ReadFileLines(file)
		if err != nil {
			logger.LogFatal("error occurred while loading proxies", err.Error())
		}
		logger.LogDebug("mapping proxies to dorks")
		dorkClients = proxy.MapStringToClient(dorksTotal, proxiesTotal, protocol)
		logger.LogDebug("mapped proxies to dorks")
	} else {
		dorkClients = make(map[string]*http.Client, len(dorksTotal))
	}
	engine = searchEngine
	writer = ui.SetupGeneralUI()
	pagePerDork = int(utils.Config.Get("page-per-dork").(int64))
	pageFormatRegex = regexp.MustCompile("(\\.)\\w*?($|\\?)")
	pageTypeRegex = regexp.MustCompile("(\\?|\\&)(.*?)(\\=)")
	receivedDomains = make(map[string]bool)
	startParse()
}

func startParse() {
	if protocol == "" {
		timeout := time.Duration(int(utils.Config.Get("timeout").(int64))) * time.Second
		client := &http.Client{Timeout: timeout}
		for _, dork := range dorksTotal {
			dorkClients[dork] = client
		}
	}
	parseDorks()
}

func parseDorks() {
	logger.LogAction("press enter to start parsing dorks")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	logger.LogInfo("press s and enter to stop")

	dorksAmount := len(dorksTotal)

	routineSplitDorks := utils.SplitSlice(dorksTotal, threads)

	ch := make(chan string, len(routineSplitDorks))

	logger.LogInfo("parsing", strconv.Itoa(dorksAmount), "dorks")

	for _, dorksList := range routineSplitDorks {
		go poolDorks(queryFunction, ch, dorksList)
	}

	totalURLCounter, validURLCounter, doneDorkCounter := 0, 0, 0
	batchURLCounter := 0

	batchURLs := []string{}

	stopProgram := false

	go func() {
		option, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if option == "s\n" {
			ui.StopGeneralUI(writer)
			logger.LogInfo("stopping")

			if len(batchURLs) != 0 {
				err := utils.SaveArrayOfString(saveFileName, batchURLs, false)
				if err != nil {
					logger.LogDebug("error while saving to file")
				}
			}

			logger.LogSuccess("stopped, parsed", strconv.Itoa(validURLCounter), "urls")
			os.Exit(0)
		}
	}()

	for newURL := range ch {
		if stopProgram == true {
			close(ch)
			break
		}
		if newURL == "done" {
			doneDorkCounter++
			ui.UpdateParseUI(writer, engine, dorksAmount, totalURLCounter, validURLCounter, doneDorkCounter)
			if dorksAmount == doneDorkCounter {
				close(ch)
				break
			}
			continue
		}
		totalURLCounter++
		batchURLCounter++
		u, err := tld.Parse(newURL)
		cleanDomain := newURL
		if err == nil {
			cleanDomain = u.Subdomain + "." + u.Domain + "." + u.TLD
		}
		if receivedDomains[cleanDomain] != true && checkValidURL(newURL) {
			logger.LogDebug("received valid url", newURL)
			validURLCounter++
			batchURLs = append(batchURLs, newURL)
		}
		receivedDomains[cleanDomain] = true
		if batchURLCounter == 10 {
			err := utils.SaveArrayOfString(saveFileName, batchURLs, false)
			if err != nil {
				logger.LogDebug("error while saving to file")
			}
			ui.UpdateParseUI(writer, engine, dorksAmount, totalURLCounter, validURLCounter, doneDorkCounter)
			batchURLCounter = 0
			batchURLs = []string{}
		}
	}

	ui.StopGeneralUI(writer)

	if len(batchURLs) != 0 {
		err := utils.SaveArrayOfString(saveFileName, batchURLs, false)
		if err != nil {
			logger.LogDebug("error while saving to file")
		}
	}

	logger.LogSuccess("done, parsed", strconv.Itoa(validURLCounter), "urls")

}

func poolDorks(queryFunction queryFuncType, ch chan string, dorkList []string) {
	for _, dork := range dorkList {
		for i := 0; i < pagePerDork; i++ {
			queryFunction(dorkClients[dork], ch, dork, i)
		}
		ch <- "done"
	}
}

func checkValidURL(url string) bool {
	if strings.Contains(url, "?") && pageFormatRegex.FindString(url) != "" && pageTypeRegex.FindString(url) != "" {
		return true
	}
	return false
}
