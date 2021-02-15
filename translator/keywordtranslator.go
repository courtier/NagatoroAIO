package translator

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
)

type (
	queryFuncType func(*http.Client, chan []string, string)
)

var (
	saveFileName string
	targetLang   string
)

//AskForStartTranslating start translation option
func AskForStartTranslating() {
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

	logger.LogInfo("uses bing translator, every language there is available here")
	logger.LogInput("target language")

	targetLang, _ = reader.ReadString('\n')
	targetLang = strings.Replace(targetLang, "\n", "", -1)

	startTranslate(keywords)
}

func startTranslate(queries []string) []string {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	loopQueries(client, queries)
	return nil
}

func loopQueries(client *http.Client, queries []string) {
	if len(queries) == 0 {
		logger.LogError("no query keywords found")
		return
	}
	logger.LogAction("press enter to start translating keywords")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ch := make(chan string, len(queries))
	translatedWords := []string{}

	for _, query := range queries {
		go doBingTranslateRequest(client, ch, query, targetLang)
		logger.LogDebug("translating", query)
	}

	counter := 0

	for translatedWord := range ch {
		translatedWords = append(translatedWords, translatedWord)
		counter++
		if counter == len(queries) {
			close(ch)
			break
		}
	}

	currKeywordLength := len(translatedWords)

	logger.LogSuccess("done, translated", strconv.Itoa(currKeywordLength), "keywords")

	logger.LogDebug("saving to file")
	saveKeywords(translatedWords)
	logger.LogInfo("stopping")
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
	return filepath.Join(folderPath, "translatedkeywords"+time+".txt")
}

func saveKeywords(translations []string) error {
	file, err := os.OpenFile(saveFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogFatal(err.Error())
	}
	defer file.Close()
	totalString := []byte{}
	for _, s := range translations {
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
