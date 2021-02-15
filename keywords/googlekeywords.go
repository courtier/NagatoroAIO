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

type googleResponseType []string

func doGoogleRequest(client *http.Client, ch chan []string, query string) {
	newSuggestions := []string{}
	url := buildGoogleRequest(query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogError("error while requesting suggestions")
		return
	}
	req.Header.Set("Authority", "suggestqueries.google.com")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Client-Data", "CIq2yQEIpbbJAQjEtskBCKmdygEIr8LKAQisx8oBCPbHygEI+MfKAQijzcoBCNzVygEI7ZjLAQi6m8sBCIqcywEIwpzLAQ==")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Dest", "script")
	req.Header.Set("Referer", "https://keywordshitter.com/")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	req.Header.Set("Cookie", "CONSENT=YES+DE.en-GB+202010; ANID=AHWqTUlMtlLN_kTgB9n5-7v4A-3zGGRk4tM1HhF3KbVmT0w4zczcgrQqZZeDie37; __Secure-3PSID=5QevD6Aw6vN4w6p85P6suH7ooaYIZA0jeYbjw3RNGvVfzKwmfP1c6Np8g03EFm3lO0t4_A.; __Secure-3PAPISID=HqC0Q5Kvs0OmG5OF/A0L5b1rB1WeOFHdvR; NID=206=JXvfPKqYdT8iJSERUC92imvp0a8eNF7iIbsIlaH2Y21KhCvS_dIZtMY6T81OWEOaZnFPDMBfxD8wwDm_xv7-b7HxSnd0yWlZmbio1XB6n52ZiFHTbLCRoZ1jnvgg121DWQBZbrBITqkv50JIPKLy-LfKLxARcvoBjTqAmXLaWyxByRdAx7YOjexWmTU3w7_1cAtIPcdLb57Y0FM9Tj4a3XWkAapgMZOyic7usjWE4U-DfyFu5idKim0fLe8zwzMBE356jznChLf6e41HcAyt-3YotNm-v1fSMu92_ySbaTE49tA; 1P_JAR=2021-01-05-21; __Secure-3PSIDCC=AJi4QfHrbs0QudidlYTrhzJ7jfAzyC_7dczPs9riJDPlCilIfyyyTge8H2fBsFMkbYJmr2pj3mw")
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

func buildGoogleRequest(query string) string {
	query = url.QueryEscape(query)
	url := "https://suggestqueries.google.com/complete/search?jsonp=true&client=chrome&q=" + query
	return url
}

func parseGoogleResponse(responseJSON string) (suggestions []string, err error) {
	responseJSON = strings.SplitN(responseJSON, ",", 2)[1]
	responseJSON = strings.SplitN(responseJSON, "],", 2)[0]
	responseJSON += "]"
	var response googleResponseType
	err = json.Unmarshal([]byte(responseJSON), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
