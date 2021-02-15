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

func parseGoogleRequest(client *http.Client, ch chan string, dork string, page int) {
	req, err := buildGoogleParseRequest(dork, page)
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
	newUrls, err := parseGoogleResponse(response)
	if err != nil {
		logger.LogDebug("error while parsing")
		return
	}
	logger.LogDebug("parsed", strconv.Itoa(len(newUrls)), "urls from dork", dork)
	for _, url := range newUrls {
		ch <- url
	}
}

func buildGoogleParseRequest(query string, page int) (*http.Request, error) {
	query = url.QueryEscape(query)
	page *= 10
	url := "https://www.google.com/search?q=" + query + "&sxsrf=ALeKk023uyiFdGcqGXTvQsZC_VvfnDBiBQ:1610550969036&ei=uQ7_X8C_AY2-sAfPsKmQAg&start=" + strconv.Itoa(page) + "&sa=N&ved=2ahUKEwiAgJvzmZnuAhUNH-wKHU9YCiIQ8tMDegQIGBA2&biw=1035&bih=845"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogDebug("error while requesting suggestions")
		return req, err
	}
	req.Header.Set("Authority", "www.google.com")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("X-Client-Data", "CIq2yQEIpbbJAQjEtskBCKmdygEIx8LKAQisx8oBCPjHygEIo83KAQjc1coBCO2YywEIk5rLAQi6m8sBCIqcywEIqZ3LAQiqncsBGPq4ygE=")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	req.Header.Set("Cookie", "CGIC=IocBdGV4dC9odG1sLGFwcGxpY2F0aW9uL3hodG1sK3htbCxhcHBsaWNhdGlvbi94bWw7cT0wLjksaW1hZ2UvYXZpZixpbWFnZS93ZWJwLGltYWdlL2FwbmcsKi8qO3E9MC44LGFwcGxpY2F0aW9uL3NpZ25lZC1leGNoYW5nZTt2PWIzO3E9MC45; CONSENT=YES+DE.en-GB+202010; SID=5QevD6Aw6vN4w6p85P6suH7ooaYIZA0jeYbjw3RNGvVfzKwm6Z4DPV_BPlLYkNaEgOKVZA.; __Secure-3PSID=5QevD6Aw6vN4w6p85P6suH7ooaYIZA0jeYbjw3RNGvVfzKwmfP1c6Np8g03EFm3lO0t4_A.; HSID=AfJF0ub5ICKfqFpay; SSID=AE1Fc_KKJfcfYncOm; APISID=a0advPvizdhwO9VN/AUV1tOuEdyEUpqmcT; SAPISID=HqC0Q5Kvs0OmG5OF/A0L5b1rB1WeOFHdvR; __Secure-3PAPISID=HqC0Q5Kvs0OmG5OF/A0L5b1rB1WeOFHdvR; OTZ=5795445_52_52_123900_48_436380; SEARCH_SAMESITE=CgQIy5EB; ANID=AHWqTUkf_MkpEeSDt3iqBmfOuNJgdlk9wqYFfj5QIgwwhA1n93NLgH84INqRdpwo; NID=207=YiG9nTxAVoAwz3_8_lxv14uEXUOWoBkrwSe6jHRFGGsp8zVY6mNo3cxXw8lZIqCmf9yvkstnDes9Ox5GDtC7BrYer4BuERc7PvN_HgyyBjoCZaKCzCDW5h_BfDaym0bbg26saOV-DPV4-zhxtOYEwN6kpNNI7swsDMHrPnbwPMpioVsEaZ_JsB-8lnMuFBN5GfkfC78Hf1MZWmG8iiZSFCWq8Q1Qwt6baIjOu5Gn62Box_cPvngMiIqJ3rLVYC-3xjCuqZZfn68BBrFr0ufde6LH9RysWT_Yi1TmAtCOFeA4QNU; 1P_JAR=2021-01-13-15; DV=U6p2y0p4-LlIIBy4P1MyROuwYxjFb9f65hDnJsyTVwUAAIAGwtOLQBGLqwEAAHC4M1I0Yx1kcQAAAA; SIDCC=AJi4QfGhuSdy6DvnJQlHYF2pcW1MjFSfpR4T3R3VMPmAbDbV42Wpu_nNrvrFWD_wXtY5AV-ed8c; __Secure-3PSIDCC=AJi4QfGt6YMSMabucdWqIlxMrbpBwQ_1r6DUPyzdpnMuWE955P-1cNDJxYy85lVW1wKv4H39pFk")
	return req, nil
}

func parseGoogleResponse(response string) ([]string, error) {
	urls := []string{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return urls, err
	}

	doc.Find(".yuRUbf a").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			link, _ = url.QueryUnescape(link)
			urls = append(urls, link)
		}
	})

	return urls, nil
}
