package proxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"h12.io/socks"
)

//BuildBetterProxyClient builds client from proxy
func BuildBetterProxyClient(proxyString, protocol string, timeout int) (*http.Client, error) {
	proxyParts := strings.Split(proxyString, ":")
	partsCount := len(proxyParts)
	var transport *http.Transport
	if strings.HasPrefix(protocol, "http") {
		if partsCount >= 2 {
			parsed, err := url.Parse(fmt.Sprintf("http://%v", proxyString))
			if err != nil {
				return nil, err
			}
			transport = &http.Transport{Proxy: http.ProxyURL(parsed)}
		}
		if partsCount == 4 {
			auth := proxyParts[2] + ":" + proxyParts[3]
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			transport.ProxyConnectHeader = http.Header{}
			transport.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
		}
	} else if strings.HasPrefix(protocol, "socks") {
		socksString := proxyString + "?timeout=" + strconv.Itoa(timeout) + "s"
		if partsCount == 4 {
			socksString = proxyParts[2] + ":" + proxyParts[3] + "@" + proxyParts[0] + ":" + proxyParts[1] + "?timeout=" + strconv.Itoa(timeout) + "s"
		}
		dial := socks.Dial(fmt.Sprintf(protocol+"://%v", socksString))
		transport = &http.Transport{Dial: dial}
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return client, nil
}
