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

func parseYahooRequest(client *http.Client, ch chan string, dork string, page int) {
	req, err := buildYahooParseRequest(dork, page)
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
	newUrls, err := parseYahooResponse(response)
	if err != nil {
		logger.LogDebug("error while parsing")
		return
	}
	logger.LogDebug("parsed", strconv.Itoa(len(newUrls)), "urls from dork", dork)
	for _, url := range newUrls {
		ch <- url
	}
}

func buildYahooParseRequest(query string, page int) (*http.Request, error) {
	query = url.QueryEscape(query)
	page *= 10
	page++
	url := "https://search.yahoo.com/search?p=" + query + "&pz=10&fr=yfp-search-sb&b=" + strconv.Itoa(page) + "&pz=10&xargs=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogDebug("error while requesting suggestions")
		return req, err
	}
	req.Header.Set("Authority", "search.yahoo.com")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://search.yahoo.com/search;_ylt=AwrJQ5wg5f1f32sAgxozCQx.;_ylu=Y29sbwNpcjIEcG9zAzEEdnRpZAMEc2VjA3BhZ2luYXRpb24-?p=minecraft&pz=10&fr=yfp-search-sb&b=11&pz=10&xargs=0")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	req.Header.Set("Cookie", "B=bgku82lfr7ksv&b=3&s=ij; A1=d=AQABBKcctF8CEJ9ArhN1CG8gKPUdV8rAecAFEgABAgFgtV95YOA9b2UB9iMAAAcIn9OzXxXIU7g&S=AQAAAnC-No5Yl-xMSeXA7XW1L3s; A3=d=AQABBKcctF8CEJ9ArhN1CG8gKPUdV8rAecAFEgABAgFgtV95YOA9b2UB9iMAAAcIn9OzXxXIU7g&S=AQAAAnC-No5Yl-xMSeXA7XW1L3s; GUC=AQABAgFftWBgeUIhmQTj; APID=1Aabfdc4be-4290-11ea-aed7-12244f280cb6; A1S=d=AQABBKcctF8CEJ9ArhN1CG8gKPUdV8rAecAFEgABAgFgtV95YOA9b2UB9iMAAAcIn9OzXxXIU7g&S=AQAAAnC-No5Yl-xMSeXA7XW1L3s&j=GDPR")
	return req, nil
}

func parseYahooResponse(response string) ([]string, error) {
	urls := []string{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return urls, err
	}

	doc.Find(".td-u.fc-5th").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			link, _ = url.QueryUnescape(link)
			if strings.Contains(link, "RU=") && strings.Contains(link, "/RK=2") {
				link = strings.Split(strings.Split(link, "RU=")[1], "/RK=2")[0]
				urls = append(urls, link)
			}
		}
	})

	return urls, nil
}
