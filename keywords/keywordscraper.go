package keywords

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
)

type (
	queryFuncType func(*http.Client, chan []string, string)
)

var (
	suggestionsTotal []string
	queryFunction    queryFuncType
	saveFileName     string
)

func askForStartScraping(queryFunc queryFuncType) {
	logger.LogInfo("place initial keywords in a file")
	logger.LogInput("file")

	reader := bufio.NewReader(os.Stdin)

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Replace(fileName, "\n", "", -1)

	keywords, err := loadKeywords(fileName)
	if err != nil {
		logger.LogFatal("error occurred while loading keywords", err.Error())
	}
	saveFileName = createFileName()
	queryFunction = queryFunc
	startScrape(keywords)
}

func startScrape(queries []string) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	suggestionsTotal = append(suggestionsTotal, queries...)
	loopQueries(client, queries)
}

func loopQueries(client *http.Client, queries []string) {
	if len(queries) == 0 {
		logger.LogError("no query keywords found")
		return
	}
	logger.LogAction("press enter to start scraping keywords")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ch := make(chan []string, len(queries))
	newBatchWords := []string{}

	for _, query := range queries {
		go queryFunction(client, ch, query)
		logger.LogDebug("scraping", query)
	}

	counter := 0

	for newWords := range ch {
		newBatchWords = append(newBatchWords, newWords...)
		counter++
		if counter == len(queries) {
			close(ch)
			break
		}
	}

	currKeywordLength := len(newBatchWords)

	completelyNewWords := utils.SubtractSlices(suggestionsTotal, newBatchWords)
	suggestionsTotal = append(suggestionsTotal, completelyNewWords...)

	logger.LogSuccess("done, scraped", strconv.Itoa(currKeywordLength), "keywords")
	logger.LogInfo("total keywords", strconv.Itoa(len(suggestionsTotal)))

	logger.LogDebug("saving to file")
	saveKeywords()

	logger.LogAction("press s and enter to stop here")
	logger.LogAction("press c and enter to continue generating")

	choiceBytes, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')
	choice := string(choiceBytes)

	if choice == "s\n" {
		logger.LogInfo("stopping")
		return
	}

	if len(completelyNewWords) == 0 {
		logger.LogInfo("no new keywords found, stopping")
		return
	}
	loopQueries(client, completelyNewWords)
}

func loadKeywords(fileName string) ([]string, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keywordsFound := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		keywordsFound = append(keywordsFound, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keywordsFound, nil
}

func createFileName() string {
	time := time.Now().Format("15:04:05")
	folderPath := filepath.Join("results", "keywords")
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		logger.LogFatal(err.Error())
	}
	return filepath.Join(folderPath, "gennedkeywords"+time+".txt")
}

func saveKeywords() error {
	file, err := os.OpenFile(saveFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogFatal(err.Error())
	}
	defer file.Close()
	totalString := []byte{}
	for _, s := range suggestionsTotal {
		totalString = append(totalString, s...)
		totalString = append(totalString, "\n"...)
	}

	err = ioutil.WriteFile(file.Name(), totalString, 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("saved to", saveFileName)
	return nil
}
