package translator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/courtier/NagatoroAIO/logger"
)

type bingTranslationResponseType []struct {
	DetectedLanguage struct {
		Language string  `json:"language"`
		Score    float64 `json:"score"`
	} `json:"detectedLanguage"`
	Translations []struct {
		Text    string `json:"text"`
		To      string `json:"to"`
		SentLen struct {
			SrcSentLen   []int `json:"srcSentLen"`
			TransSentLen []int `json:"transSentLen"`
		} `json:"sentLen"`
	} `json:"translations"`
}

func doBingTranslateRequest(client *http.Client, ch chan string, query, lang string) {
	url := buildBingURL(query, lang)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logger.LogError("error while requesting translation")
		return
	}
	req.Header.Set("Authority", "www.bing.com")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Arch", "\"x86\"")
	req.Header.Set("Sec-Ch-Ua-Full-Version", "\"87.0.4280.88\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"10_13_6\"")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Ch-Ua-Model", "")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Mac OS X\"")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.bing.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.bing.com/translator")
	req.Header.Set("Accept-Language", "en-AU,en;q=0.9,de-DE;q=0.8,de;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	req.Header.Set("Cookie", "MUID=1B7E0EE4676664FC1CCD001E667F6569; SRCHD=AF=NOFORM; SRCHUID=V=2&GUID=7390BE7564144DFE863C7076BE53E6A8&dmnchg=1; MUIDB=1B7E0EE4676664FC1CCD001E667F6569; BCP=AD=0&AL=0&SM=0; _RwBf=mtu=0&g=0&o=2&p=&c=&t=0&s=0001-01-01T00:00:00.0000000+00:00&ts=2021-01-05T12:05:08.3193637+00:00&ssg=0&cid=; _SS=SID=20720A56BE9D66DE243B05E0BF776751&R=5&RB=0&GB=0&RG=200&RP=0; SRCHUSR=DOB=20200728&T=1610110593000; SRCHHPGUSR=CW=1035&CH=822&DPR=2&UTC=60&DM=0&HV=1610110596&WTS=63745707393&BRW=HTP&BRH=M; _EDGE_S=SID=08D48CB7D7516BD40FB2830ED6BB6A73; ipv6=hit=1610114196656&t=6; _TTSS_IN=hist=WyJlbiIsImF1dG8tZGV0ZWN0Il0=; _tarLang=default=de; _TTSS_OUT=hist=WyJlbiIsImRhIiwiZGUiXQ==")
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
	translation, err := parseBingTranslationResponse(responseJSON)
	if err != nil {
		logger.LogError("error while parsing json")
		return
	}
	logger.LogDebug("translated", query, "to", translation)
	ch <- translation
}

func buildBingURL(query, lang string) string {
	url := "https://www.bing.com/ttranslatev3?isVertical=1&&IG=8F3339C5EBEF469DB8E020763037D113&IID=translator.5022.1&fromLang=en&text=" + query + "&to=" + lang
	return url
}

func parseBingTranslationResponse(responseJSON string) (string, error) {
	fmt.Println(responseJSON)
	var response bingTranslationResponseType
	translation := ""
	err := json.Unmarshal([]byte(responseJSON), &response)
	if err != nil {
		return translation, err
	}
	translation = response[0].Translations[0].Text
	return translation, nil
}
