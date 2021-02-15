package keywords

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
)

type startpageResponseType []string

func doStartpageRequest(client *http.Client, ch chan []string, query string) {
	url := buildStartpageRequest(query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogError("error while requesting suggestions")
		return
	}
	req.Header.Set("Authority", "www.startpage.com")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://www.startpage.com/")
	resp, err := client.Do(req)
	if err != nil {
		logger.LogError("error while processing response")
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError("error while reading body")
		return
	}
	responseJSON := string(respBody)
	newSuggestions, err := parseStartpageResponse(responseJSON)
	if err != nil {
		logger.LogError("error while parsing json")
		return
	}
	logger.LogDebug("scraped", strconv.Itoa(len(newSuggestions)), "keywords from keyword", query)
	ch <- newSuggestions
}

func buildStartpageRequest(query string) string {
	query = url.QueryEscape(query)
	url := "https://www.startpage.com/do/suggest?limit=15&lang=english&format=json&query=" + query
	return url
}

func parseStartpageResponse(responseJSON string) (suggestions []string, err error) {
	responseJSON = strings.SplitN(responseJSON, ",", 2)[1]
	responseJSON = responseJSON[:len(responseJSON)-1]
	var response startpageResponseType
	err = json.Unmarshal([]byte(responseJSON), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
