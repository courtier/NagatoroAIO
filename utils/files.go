package utils

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/courtier/NagatoroAIO/logger"
)

//AskForSpecificFileName asks for specific file name
func AskForSpecificFileName(specificName string) string {
	logger.LogInput(specificName + " file")
	reader := bufio.NewReader(os.Stdin)

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Replace(fileName, "\n", "", -1)
	return fileName
}

//AskForAnything asks for anything
func AskForAnything(anything string) string {
	logger.LogInput(anything)
	reader := bufio.NewReader(os.Stdin)

	response, _ := reader.ReadString('\n')
	response = strings.Replace(response, "\n", "", -1)
	return response
}

//AskForFileName asks for file name
func AskForFileName() string {
	logger.LogInput("file")
	reader := bufio.NewReader(os.Stdin)

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Replace(fileName, "\n", "", -1)
	return fileName
}

//OpenFile opens file
func OpenFile(fileName string) (*os.File, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

//ReadFileLines read file line by line
func ReadFileLines(file *os.File) ([]string, error) {
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

//SaveArrayOfString saves array of string to file
func SaveArrayOfString(fileName string, lines []string, overwrite bool) error {
	totalString := []byte{}
	for _, s := range lines {
		sBytes := []byte(s + "\n")
		totalString = append(totalString, sBytes...)
	}

	if overwrite {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(file.Name(), totalString, 0644); err != nil {
			return err
		}
	} else {
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.Write(totalString); err != nil {
			return err
		}
		file.Sync()
	}

	return nil
}

//AppendString appends string to file
func AppendString(fileName string, line string, overwrite bool) error {
	sBytes := []byte(line + "\n")

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(sBytes); err != nil {
		return err
	}
	file.Sync()

	return nil
}

//CreateFileName creates file name with time appended
func CreateFileName(folder, words string) string {
	time := time.Now().Format("15:04:05")
	folderPath := filepath.Join("results", folder)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		logger.LogFatal(err.Error())
	}
	return filepath.Join(folderPath, words+time+".txt")
}
