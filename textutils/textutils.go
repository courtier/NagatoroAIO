package textutils

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
)

//AskForUtils choose from util modules
func AskForUtils() {
	logger.LogCategory("modules")
	logger.LogOption(1, "replace all")
	logger.LogOption(2, "remove dupes")
	logger.LogOption(3, "mix all words")
	logger.LogOption(4, "remove numbers")
	logger.LogOption(5, "clean spaces")
	logger.LogOption(6, "surround with quotation marks")

	reader := bufio.NewReader(os.Stdin)

	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)

	if option == "1" {
		replaceAllSomething()
	} else if option == "2" {
		removeDupes()
	} else if option == "3" {
		randomizer()
	} else if option == "4" {
		removeNumbers()
	} else if option == "5" {
		cleanSpaces()
	} else if option == "6" {
		placeQuotations()
	} else {
		logger.LogError("coming soon))")
	}
}

func replaceAllSomething() {
	fileName := utils.AskForFileName()

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)
	logger.LogInput("what to replace")
	toReplace, _ := reader.ReadString('\n')
	toReplace = strings.ReplaceAll(toReplace, "\n", "")

	logger.LogInput("what to replace with")
	replaceWith, _ := reader.ReadString('\n')
	replaceWith = strings.ReplaceAll(replaceWith, "\n", "")

	totalString := []byte{}
	for _, line := range lines {
		line = line + "\n"
		line = strings.ReplaceAll(line, toReplace, replaceWith)
		totalString = append(totalString, []byte(line)...)
	}

	err = ioutil.WriteFile(fileName, []byte(totalString), 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("Done")
}

func removeDupes() {
	fileName := utils.AskForFileName()

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}
	defer file.Close()

	lines = utils.RemoveDupesOffSlice(lines)

	totalString := ""
	for _, line := range lines {
		totalString += line + "\n"
	}

	err = ioutil.WriteFile(fileName, []byte(totalString), 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("done")
}

func removeNumbers() {
	fileName := utils.AskForFileName()

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}
	defer file.Close()

	reg := regexp.MustCompile("[0-9]")
	totalString := ""
	for _, line := range lines {
		line = reg.ReplaceAllString(line, "")
		totalString += line + "\n"
	}

	err = ioutil.WriteFile(fileName, []byte(totalString), 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("done")
}

func cleanSpaces() {
	fileName := utils.AskForFileName()

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}
	defer file.Close()

	totalString := ""
	space := regexp.MustCompile(`\s+`)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = space.ReplaceAllString(line, " ")
		totalString += line + "\n"
	}

	err = ioutil.WriteFile(fileName, []byte(totalString), 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("done")
}

func placeQuotations() {
	fileName := utils.AskForFileName()

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}
	defer file.Close()

	totalString := ""
	for _, line := range lines {
		totalString += "\"" + line + "\"\n"
	}

	err = ioutil.WriteFile(fileName, []byte(totalString), 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("done")
}

func randomizer() {

	fileName := utils.AskForFileName()

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}
	defer file.Close()

	words := []string{}

	for _, line := range lines {
		wordsInLine := strings.Split(line, " ")
		for _, word := range wordsInLine {
			words = append(words, word)
		}
	}

	words = utils.RemoveDupesOffSlice(words)

	totalString := []byte{}
	for _, wordX := range words {
		wordX = wordX + "\n"
		totalString = append(totalString, []byte(wordX)...)
	}
	for _, wordX := range words {
		for _, wordY := range words {
			s := append([]byte(wordX), []byte(wordY)...)
			s = append(s, []byte("\n")...)
			totalString = append(totalString, []byte(s)...)
		}
	}
	for _, wordX := range words {
		for _, wordY := range words {
			for _, wordZ := range words {
				s := append([]byte(wordX), []byte(wordY)...)
				s = append(s, []byte(wordZ)...)
				s = append(s, []byte("\n")...)
				totalString = append(totalString, []byte(s)...)
			}
		}
	}

	err = ioutil.WriteFile(fileName, totalString, 0644)
	if err != nil {
		logger.LogFatal("error while writing file", err.Error())
	}

	logger.LogSuccess("done")
}
