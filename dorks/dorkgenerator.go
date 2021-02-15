package dorks

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
)

func askStartDorkGenerator() {
	logger.LogInput("dork formats file")
	reader := bufio.NewReader(os.Stdin)

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Replace(fileName, "\n", "", -1)
	dorkFormats, err := parseFormatsFromFile(fileName)
	if err != nil {
		logger.LogFatal("error while loading formats", err.Error())
	}
	containedKeys := make(map[string]bool)
	for _, dorkFormat := range dorkFormats {
		for _, key := range dorkFormat.requiredParamKeys {
			containedKeys[key] = true
		}
	}
	allParamsMapList := make(map[string][]string)
	for requiredKey := range containedKeys {
		paramName := paramShortToLong(requiredKey)
		fileName = askForFileName(reader, paramName)
		paramList, err := loadParameterFile(fileName)
		if err != nil {
			logger.LogFatal("error while loading", paramName)
		}
		allParamsMapList[requiredKey] = paramList
	}
	logger.LogInfo("dork generator starting")
	filledResults := []string{}
	counter := 0
	for _, dorkFormat := range dorkFormats {
		fillinParamsList := [][]string{}
		for _, paramKey := range dorkFormat.requiredParamKeys {
			fillinParamsList = append(fillinParamsList, allParamsMapList[paramKey])
		}
		combinedChannel := utils.Iter(fillinParamsList)
		for product := range combinedChannel {
			stringProduct := utils.InterfaceArrToStringArr(product)
			filled, err := fillFormat(dorkFormat, stringProduct)
			if err != nil {
				logger.LogError("error while filling in format", err.Error())
				continue
			}
			counter++
			filledResults = append(filledResults, filled)
		}
	}
	logger.LogSuccess("generated", strconv.Itoa(counter), "dorks")

	saveFileName := utils.CreateFileName("dorks", "genneddorks")
	logger.LogInfo("saving to", saveFileName)

	file, err := os.OpenFile(saveFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogFatal("error while opening file", err.Error())
	}
	writer := bufio.NewWriter(file)
	var stringBuilder strings.Builder
	for _, result := range filledResults {
		stringBuilder.WriteString(result + "\n")
	}
	_, err = writer.WriteString(stringBuilder.String())
	if err != nil {
		logger.LogFatal("error while saving to file", err.Error())
	}
	writer.Flush()
}

func loadParameterFile(fileName string) ([]string, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	list := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		list = append(list, line)
	}
	return list, nil
}

func askForFileName(reader *bufio.Reader, name string) string {
	logger.LogInput(name, "file")

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Replace(fileName, "\n", "", -1)
	return fileName
}

func paramShortToLong(key string) string {
	switch key {
	case "(KW)":
		return "keyword"
	case "(PF)":
		return "page format"
	case "(PT)":
		return "page type"
	case "(DE)":
		return "domain extension"
	default:
		return "keyword"
	}
}
