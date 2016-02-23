package main

import "net/http"
import "time"

func CheckURLExists(url string) bool {
	n := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		DebugOut("ERR:%v\n", err)
		return false
	}
	DebugOut("CheckURLExists: %v %v %v\n", url, resp.Status, time.Now().Sub(n))
	defer resp.Body.Close()

	switch resp.StatusCode / 100 {
	case 4, 5:
		return false
	case 3, 2, 1:
		return true
	}
	return false
}

func init() {
}

func ParseIndex(indexUrl string) ([]string, error) {
	resp, err := http.Get(indexUrl)
	if err != nil {
		return nil, err
	}
	return DecodeIndex(resp.Body)
}

func DetectMirrorProgress(mirrorUrl string, guards []string) float64 {
	result := make(chan float64)
	queue := make(chan bool, *DetecThreadN)

	go func() {
		ok, count := 0, 0
		for v := range queue {
			if v {
				ok = ok + 1
			}
			count = count + 1
			if count == len(guards) {
				break
			}
		}
		result <- float64(ok) / float64(len(guards))
	}()
	for _, g := range guards {
		go func(target string) {
			queue <- CheckURLExists(target)
		}(mirrorUrl + "/" + g)
	}
	return <-result
}

func HandleReportMirrorProgress(indexUrl string, mirrorUrl string) (float64, error) {
	guards, err := ParseIndex(indexUrl)
	if err != nil {
		return 0, err
	}
	return DetectMirrorProgress(mirrorUrl, guards), nil
}
