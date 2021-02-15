package parser

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/courtier/NagatoroAIO/logger"
)

func parseOkeanoRequest(client *http.Client, ch chan string, dork string, page int) {
	req, err := buildOkeanoParseRequest(dork, page)
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
	newUrls, err := parseOkeanoResponse(response)
	if err != nil {
		logger.LogDebug("error while parsing")
		return
	}
	logger.LogDebug("parsed", strconv.Itoa(len(newUrls)), "urls from dork", dork)
	for _, url := range newUrls {
		ch <- url
	}
}

func buildOkeanoParseRequest(query string, page int) (*http.Request, error) {
	query = url.QueryEscape(query)
	page++
	url := "https://okeano.com/search?&q=" + query + "&p=" + strconv.Itoa(page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogDebug("error while requesting suggestions")
		return req, err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://okeano.com/search?q=minecraft")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	req.Header.Set("Cookie", "sc=1")
	return req, nil
}

func parseOkeanoResponse(response string) ([]string, error) {
	urls := []string{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return urls, err
	}

	doc.Find(".result-url").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			urls = append(urls, link)
		}
	})

	return urls, nil
}
