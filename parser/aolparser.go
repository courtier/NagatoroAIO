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

func parseAolRequest(client *http.Client, ch chan string, dork string, page int) {
	req, err := buildAolParseRequest(dork, page)
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
	newUrls, err := parseAolResponse(response)
	if err != nil {
		logger.LogDebug("error while parsing")
		return
	}
	logger.LogDebug("parsed", strconv.Itoa(len(newUrls)), "urls from dork", dork)
	for _, url := range newUrls {
		ch <- url
	}
}

func buildAolParseRequest(query string, page int) (*http.Request, error) {
	query = url.QueryEscape(query)
	page *= 10
	page++
	url := "https://search.aol.com/aol/search?p=" + query + "&pz=10&fr=yfp-search-sb&b=" + strconv.Itoa(page) + "&pz=10&xargs=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogDebug("error while requesting suggestions")
		return req, err
	}
	req.Header.Set("Authority", "search.aol.com")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	req.Header.Set("Cookie", "rxx=6ssbahj4rr.24ij9p17&v=1; A1=d=AQABBCV-qV8CEFU0AEG41x6c9V16lAGo450FEgEAAgIjtF-AYNwQyiMA_SMAAAcIJH6pX2oi2owID4kt8DFU_nP-kTIt2jA_EwkBAAoBAg&S=AQAAAn-y58ghLjfYiOThDk7wX4g; A3=d=AQABBCV-qV8CEFU0AEG41x6c9V16lAGo450FEgEAAgIjtF-AYNwQyiMA_SMAAAcIJH6pX2oi2owID4kt8DFU_nP-kTIt2jA_EwkBAAoBAg&S=AQAAAn-y58ghLjfYiOThDk7wX4g; GUC=AQEAAgJftCNggEIf-ASf; A1S=d=AQABBCV-qV8CEFU0AEG41x6c9V16lAGo450FEgEAAgIjtF-AYNwQyiMA_SMAAAcIJH6pX2oi2owID4kt8DFU_nP-kTIt2jA_EwkBAAoBAg&S=AQAAAn-y58ghLjfYiOThDk7wX4g&j=GDPR; BX=8pmh2d9fqivh4&b=4&d=22E71ZltYFmHY_INqyck&s=f9&i=iS3wMVT.c_6RMi3aMD8T; sBS=dpr=2&vw=1035&vh=831; x_ms=cltid=af84dc3fc03c5b64802a5e655476bcb9")
	return req, nil
}

func parseAolResponse(response string) ([]string, error) {
	urls := []string{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return urls, err
	}

	doc.Find(".ac-algo.fz-l.ac-21th.lh-24").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			link, _ = url.QueryUnescape(link)
			if strings.Contains(link, "RU=") && strings.Contains(link, "/RK=0") {
				link = strings.Split(strings.Split(link, "RU=")[1], "/RK=0")[0]
				urls = append(urls, link)
			}
		}
	})

	return urls, nil
}
