package parser

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/courtier/NagatoroAIO/logger"
	"golang.org/x/net/html"
)

func parseStartPageRequest(client *http.Client, ch chan string, dork string, page int) {
	req, err := buildStartpageParseRequest(dork, page)
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
	newUrls, err := parseStartpageResponse(response)
	if err != nil {
		logger.LogDebug("error while parsing")
		return
	}
	logger.LogDebug("parsed", strconv.Itoa(len(newUrls)), "urls from dork", dork)
	for _, url := range newUrls {
		ch <- url
	}
}

func buildStartpageParseRequest(query string, page int) (*http.Request, error) {
	query = url.QueryEscape(query)
	url := "https://www.startpage.com/sp/search?lui=english&language=english&cat=web&sc=EL1X9eGHcRFl20&abp=-1&page=" + strconv.Itoa(page) + "&query=" + query
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogDebug("error while requesting suggestions")
		return req, err
	}
	req.Header.Set("Authority", "www.startpage.com")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://www.startpage.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.startpage.com/")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	return req, nil
}

func parseStartpageResponse(response string) ([]string, error) {
	urls := []string{}

	//TODO
	//w-gl__result-url result-link
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return urls, err
	}

	doc.Find(".w-gl__result-url.result-link").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			urls = append(urls, link)
		}
	})

	return urls, nil
}

func getHref(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
