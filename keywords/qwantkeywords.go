package keywords

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/courtier/NagatoroAIO/logger"
)

type qwantSuggestionType struct {
	value       string
	suggestType int
}

type qwantResponseType struct {
	status string
	data   struct {
		items []qwantSuggestionType
	}
	special        []interface{}
	availableQwick []interface{}
}

func doQwantRequest(client *http.Client, ch chan []string, query string) {
	url := buildStartpageRequest(query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogError("error while requesting suggestions")
		return
	}
	req.Header.Set("Authority", "api.qwant.com")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.qwant.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.qwant.com/")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
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

func buildQwantRequest(query string) string {
	query = url.QueryEscape(query)
	url := "https://api.qwant.com/api/suggest?lang=en_en&q=" + query
	return url
}

func parseQwantResponse(responseJSON string) ([]string, error) {
	var response qwantResponseType
	err := json.Unmarshal([]byte(responseJSON), &response)
	if err != nil {
		return nil, err
	}
	suggestions := []string{}
	for _, suggestion := range response.data.items {
		suggestions = append(suggestions, suggestion.value)
	}
	return suggestions, nil
}
