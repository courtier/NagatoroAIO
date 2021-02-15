package proxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/courtier/NagatoroAIO/utils"
	"golang.org/x/net/proxy"
)

//BuildProxyClient builds client from proxy
func BuildProxyClient(proxyString, protocol string, timeout int) *http.Client {
	proxyParts := strings.Split(proxyString, ":")
	partsCount := len(proxyParts)
	var transport *http.Transport
	if protocol == "http" || protocol == "https" {
		var proxyURL *url.URL
		if partsCount == 4 {
			proxyAddress := proxyParts[0] + ":" + proxyParts[1]
			proxyURL, _ = url.Parse("http://" + proxyAddress)
			auth := proxyParts[2] + ":" + proxyParts[3]
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			authHeader := http.Header{}
			authHeader.Add("Proxy-Authorization", basicAuth)
			transport = &http.Transport{
				Proxy:              http.ProxyURL(proxyURL),
				ProxyConnectHeader: authHeader,
			}
		} else {
			proxyURL, _ = url.Parse("http://" + proxyString)
			transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	} else {
		fmt.Println("username:", proxyParts[2])
		auth := proxy.Auth{
			User:     proxyParts[2],
			Password: proxyParts[3],
		}
		//https://play.golang.org/p/l0iLtkD1DV
		dialer, err := proxy.SOCKS5("tcp", proxyParts[0]+":"+proxyParts[1], &auth, proxy.Direct)
		if err != nil {
			logger.LogError("can't connect to the proxy:", err.Error())
		}
		transport = &http.Transport{}
		transport.Dial = dialer.Dial
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return client
}

//MapStringToClient maps proxies to string
func MapStringToClient(first, second []string, protocol string) map[string]*http.Client {
	mapped := make(map[string]*http.Client, len(first))
	timeout := int(utils.Config.Get("timeout").(int64))
	if len(first) <= len(second) {
		for i, el := range first {
			proxy := second[i]
			client, err := BuildBetterProxyClient(proxy, protocol, timeout)
			if err != nil {
				logger.LogError("error while creating client for proxy")
			}
			mapped[el] = client
		}
		return mapped
	}
	batchSize := len(first) / len(second)
	if batchSize > len(second) {
		batchSize = len(second)
	}
	proxyClientMapped := make(map[string]*http.Client, len(second))
	for _, proxy := range second {
		client, err := BuildBetterProxyClient(proxy, protocol, timeout)
		if err != nil {
			logger.LogError("error while creating client for proxy")
		}
		proxyClientMapped[proxy] = client
	}
	totalCounter, smallCounter := 0, 0
	for x := 0; x < len(first); x++ {
		if smallCounter >= len(second) {
			smallCounter = 0
		}
		mapped[first[totalCounter]] = proxyClientMapped[second[smallCounter]]
		totalCounter++
		smallCounter++
	}
	return mapped
}
