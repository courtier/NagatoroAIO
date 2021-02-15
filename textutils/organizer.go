package textutils

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
)

func categorize() {
	logger.LogInput("file")
	reader := bufio.NewReader(os.Stdin)

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Replace(fileName, "\n", "", -1)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		logger.LogFatal("file does not exist")
	}

	file, err := os.Open(fileName)
	if err != nil {
		logger.LogFatal("error while opening file", err.Error())
	}
	defer file.Close()

	logger.LogInfo("use regex to categorize the file")
	logger.LogInfo("say end to start categorizing after entering regexes")
	logger.LogInput("first category")

	category, _ := reader.ReadString('\n')
	category = strings.Replace(category, "\n", "", -1)

	categories := []regexp.Regexp{}

	for category != "end" {
		categoryRegex := regexp.MustCompile(category)
		categories = append(categories, *categoryRegex)
		category, _ = reader.ReadString('\n')
		category = strings.Replace(category, "\n", "", -1)
	}

	scanner := bufio.NewScanner(file)
	categorizedLines := make(map[*regexp.Regexp][]string)
	for scanner.Scan() {
		line := scanner.Text()
		for _, reg := range categories {
			if reg.Match([]byte(line)) {
				categorizedLines[&reg] = append(categorizedLines[&reg], line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		logger.LogFatal("error while reading file", err.Error())
	}

	totalString := ""
	for _, reg := range categorizedLines {
		for _, line := range reg {
			totalString += line + "\n"
		}
	}

	err = ioutil.WriteFile(fileName, []byte(totalString), 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("done")
}
