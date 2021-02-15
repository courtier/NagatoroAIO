package keywords

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
)

type yahooResponseType []string

func doYahooRequest(client *http.Client, ch chan []string, query string) {
	newSuggestions := []string{}
	url := buildGoogleRequest(query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogError("error while requesting suggestions")
		return
	}
	req.Header.Set("Authority", "sugg.search.yahoo.net")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Dest", "script")
	req.Header.Set("Referer", "https://soovle.com/")
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
	newSuggestions, err = parseGoogleResponse(responseJSON)
	if err != nil {
		logger.LogError("error while parsing json")
		return
	}
	logger.LogDebug("scraped", strconv.Itoa(len(newSuggestions)), "keywords from keyword", query)
	ch <- newSuggestions
}

func buildYahooRequest(query string) string {
	query = url.QueryEscape(query)
	time := strconv.FormatInt(utils.NanoToMilliStamp(), 10)
	url := "https://sugg.search.yahoo.net/sg/?output=fxjsonp&nresults=20&_=" + time + "&command=" + query
	return url
}

func parseYahooResponse(responseJSON string) (suggestions []string, err error) {
	responseJSON = strings.SplitN(responseJSON, ",", 2)[1]
	responseJSON = strings.SplitN(responseJSON, "],", 2)[0]
	responseJSON += "]"
	fmt.Println(responseJSON)
	var response googleResponseType
	err = json.Unmarshal([]byte(responseJSON), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
