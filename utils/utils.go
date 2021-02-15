package utils

import (
	"math/rand"
	"time"
)

//SubtractSlices subtract smaller array from larger one
func SubtractSlices(firstArr, secondArr []string) []string {
	found := []string{}
	exists := make(map[string]int)

	for _, element := range firstArr {
		exists[element] = 1
	}

	for _, element := range secondArr {
		if exists[element] != 1 {
			found = append(found, element)
		}
	}
	return found
}

//SortArraysBySize sorts an array of arrays by size
func SortArraysBySize(arrays [][]string) [][]string {
	for i := 0; i < len(arrays)-1; i++ {
		for i := 0; i < len(arrays)-i-1; i++ {
			if len(arrays[i]) > len(arrays[i+1]) {
				arrays[i], arrays[i+1] = arrays[i+1], arrays[i]
			}
		}
	}
	return arrays
}

//InterfaceArrToStringArr convert intr. arr. to str. arr.
func InterfaceArrToStringArr(inter []interface{}) []string {
	str := make([]string, len(inter))
	for i, v := range inter {
		str[i] = v.(string)
	}
	return str
}

//NanoToMilliStamp convert nano timestamp to milliseconds
func NanoToMilliStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//GenRandomSentence gen sentence
func GenRandomSentence() string {
	sentences := []string{"count yo stacks", "purified water in my glass", "eyes low like im 5'8", "these raps show up on the richter",
		"broke boy, aint broke", "fill up that hollow stomach with my sorrow", "im just gonna much up a little bit man"}
	rand.Seed(time.Now().Unix())
	return sentences[rand.Intn(len(sentences))]
}

//RemoveDupesOffSlice dedupe slice, keep sort
func RemoveDupesOffSlice(arr []string) []string {
	exists := make(map[string]bool)
	newArr := []string{}
	for _, s := range arr {
		if exists[s] != true {
			newArr = append(newArr, s)
			exists[s] = true
		}
	}
	return newArr
}

//ContainsString check if str slice contains string
func ContainsString(arr []string, el string) bool {
	for _, s := range arr {
		if s == el {
			return true
		}
	}
	return false
}

//SplitSlice splits slice into n pieces
func SplitSlice(arr []string, n int) [][]string {
	if len(arr) < n {
		return [][]string{arr}
	}
	batchSize := len(arr) / n
	var batchAmount int
	if len(arr)%batchSize == 0 {
		batchAmount = len(arr) / batchSize
	} else {
		batchAmount = (len(arr) / batchSize) + 1
	}
	split := make([][]string, batchAmount)
	for i := 0; i < batchAmount; i++ {
		min := i * batchSize
		max := (i + 1) * batchSize
		if max > len(arr) {
			max = len(arr)
		}
		split[i] = arr[min:max]
	}
	return split
}

//IsProxyProtocolValid validate protocol
func IsProxyProtocolValid(protocol string) bool {
	if protocol != "http" && protocol != "https" && protocol != "socks4" && protocol != "socks4a" && protocol != "socks5" {
		return false
	}
	return true
}
