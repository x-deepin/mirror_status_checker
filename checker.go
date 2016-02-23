package main

import "net/http"

func CheckURLExists(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
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
	ok := 0
	for _, g := range guards {
		url := mirrorUrl + "/" + g
		v := CheckURLExists(url)
		if v {
			ok = ok + 1
		}
	}
	return float64(ok) / float64(len(guards))
}

func HandleReportMirrorProgress(indexUrl string, mirrorUrl string) (float64, error) {
	guards, err := ParseIndex(indexUrl)
	if err != nil {
		return 0, err
	}
	return DetectMirrorProgress(mirrorUrl, guards), nil
}
