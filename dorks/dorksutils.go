package dorks

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
)

func askForURLToParamsConversion() {
	fileName := utils.AskForSpecificFileName("url")

	file, err := utils.OpenFile(fileName)
	if err != nil {
		logger.LogFatal("error while opening file")
	}
	defer file.Close()

	lines, err := utils.ReadFileLines(file)
	if err != nil {
		logger.LogFatal("error while reading file")
	}

	logger.LogInfo("starting conversion")

	pageFormats := []string{}
	pageTypes := []string{}
	domainExtensions := []string{}

	for _, line := range lines {
		foundFormat, foundTypes, foundExtension := parseFormatFromURL(line)
		pageFormats = append(pageFormats, foundFormat...)
		pageTypes = append(pageTypes, foundTypes...)
		domainExtensions = append(domainExtensions, foundExtension)
		logger.LogDebug("converted", line)
	}

	logger.LogSuccess("converted", strconv.Itoa(len(pageFormats)), "page formats,",
		strconv.Itoa(len(pageTypes)), "page types, and",
		strconv.Itoa(len(domainExtensions)), "domain extensions")

	time := time.Now().Format("15:04:05")
	folder := filepath.Join("converted", time)

	logger.LogInfo("deduping")

	pageFormats = utils.RemoveDupesOffSlice(pageFormats)
	pageTypes = utils.RemoveDupesOffSlice(pageTypes)
	domainExtensions = utils.RemoveDupesOffSlice(domainExtensions)

	logger.LogSuccess("deduped")

	logger.LogInfo("saving to files in", folder)

	formatsFileName := utils.CreateFileName(folder, "pageformats")
	typesFileName := utils.CreateFileName(folder, "pagetypes")
	extensionsFileName := utils.CreateFileName(folder, "domainextensions")
	if err = utils.SaveArrayOfString(formatsFileName, pageFormats, true); err != nil {
		logger.LogError("error while saving page formats", err.Error())
	}
	if err = utils.SaveArrayOfString(typesFileName, pageTypes, true); err != nil {
		logger.LogError("error while saving page types", err.Error())
	}
	if err = utils.SaveArrayOfString(extensionsFileName, domainExtensions, true); err != nil {
		logger.LogError("error while saving domain extensions", err.Error())
	}
	if err == nil {
		logger.LogSuccess("saved to file")
	}
}
