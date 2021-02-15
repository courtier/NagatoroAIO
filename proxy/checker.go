package proxy

import (
	"net/http"
)

var checkRequest *http.Request

func checkProxyRequest(ch chan string, proxy string) {
	client, err := BuildBetterProxyClient(proxy, protocol, timeout)
	if err != nil {
		ch <- "broken"
		client = nil
		return
	}
	resp, err := client.Do(checkRequest)
	if err != nil {
		ch <- "broken"
		client = nil
		resp = nil
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		ch <- proxy
		client = nil
		resp = nil
	} else {
		ch <- "broken"
		client = nil
		resp = nil
	}
}

func buildCheckRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", "http://httpstat.us/200", nil)
	if err != nil {
		return req, err
	}
	return req, nil
}
