package dorks

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	tld "github.com/jpillora/go-tld"
)

type dorkFormat struct {
	content string
	//0 for kw, 1 for pf, 2 for pt, 3 for de
	requiredParamKeys []string
}

func parseFormatsFromFile(fileName string) ([]dorkFormat, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	formats := []dorkFormat{}
	for scanner.Scan() {
		line := scanner.Text()
		containsKeys := []string{}
		if strings.Contains(line, "(KW)") {
			containsKeys = append(containsKeys, "(KW)")
		}
		if strings.Contains(line, "(PF)") {
			containsKeys = append(containsKeys, "(PF)")
		}
		if strings.Contains(line, "(PT)") {
			containsKeys = append(containsKeys, "(PT)")
		}
		if strings.Contains(line, "(DE)") {
			containsKeys = append(containsKeys, "(DE)")
		}
		format := dorkFormat{
			content:           line,
			requiredParamKeys: containsKeys,
		}
		formats = append(formats, format)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return formats, nil
}

func fillFormat(format dorkFormat, args []string) (string, error) {
	var filled string
	if len(format.requiredParamKeys) != len(args) {
		fmt.Println(len(args))
		fmt.Println(len(format.requiredParamKeys))
		return filled, errors.New("missing field")
	}
	filled = format.content
	for i, key := range format.requiredParamKeys {
		filled = strings.ReplaceAll(filled, key, args[i])
	}
	return filled, nil
}

func parseFormatFromURL(url string) ([]string, []string, string) {
	pageFormatRegex := regexp.MustCompile("(\\.)\\w*?($|\\?|\\&)")
	pageFormats := pageFormatRegex.FindAllString(url, -1)
	pageFormats = trimFromEachEnd(pageFormats)
	pageTypeRegex := regexp.MustCompile("(\\?|\\&)(.*?)(\\=)")
	pageTypes := pageTypeRegex.FindAllString(url, -1)
	pageTypes = trimFromEachEnd(pageTypes)
	parsedURL, _ := tld.Parse(url)
	return pageFormats, pageTypes, parsedURL.TLD
	//page format string between (. or /) and ?
	//page type -> strings between "?" x "=" and "&" x "="
}

func trimFromEachEnd(arr []string) []string {
	for in, el := range arr {
		arr[in] = el[1 : len(el)-1]
	}
	return arr
}
