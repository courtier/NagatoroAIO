package proxy

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/courtier/NagatoroAIO/logger"
)

func scrapeProxyRequest(client *http.Client, ch chan string, source string) {
	req, err := buildScrapeRequest(source)
	if err != nil {
		logger.LogDebug("error while building request")
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.LogDebug("error while processing response")
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogDebug("error while reading body")
		return
	}
	response := string(respBody)
	newProxies, err := parseScrapeResponse(response)
	if err != nil {
		logger.LogDebug("error while parsing")
		return
	}
	logger.LogDebug("parsed", strconv.Itoa(len(newProxies)), "proxies from source", source)
	for _, url := range newProxies {
		ch <- url
	}
}

func buildScrapeRequest(source string) (*http.Request, error) {
	req, err := http.NewRequest("GET", source, nil)
	if err != nil {
		return req, err
	}
	return req, nil
}

func parseScrapeResponse(response string) ([]string, error) {
	proxies := []string{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return proxies, err
	}

	response = cleanRegex.ReplaceAllLiteralString(response, "")
	response = strings.ReplaceAll(response, "\n", " ")
	words := strings.Split(response, " ")
	for _, word := range words {
		if len(word) < 9 {
			continue
		}
		if ipPortRegex.MatchString(word) {
			proxies = append(proxies, word)
		} else if ipRegex.MatchString(word) {
			//'td:contains(word)'
			doc.Find("td:contains(" + word + ")").Each(func(_ int, s *goquery.Selection) {
				nextSibling := s.Siblings().Nodes[0]
				potentialPort := nextSibling.Data
				if portRegex.MatchString(potentialPort) {
					proxy := word + ":" + potentialPort
					proxies = append(proxies, proxy)
				}
			})
		}
	}

	return proxies, nil
}
