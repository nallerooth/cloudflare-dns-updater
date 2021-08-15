package iplookup

import (
	"io"
	"net/http"
	"regexp"
)

const IPLookupURL = "https://nallerooth.com/ip.php"

func fetchPage() ([]byte, error) {
	req, err := http.Get(IPLookupURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func verifyIPv4Addr(body []byte) (string, error) {
	pattern := regexp.MustCompile(`(\d{1,3}\.){3}\d+`)
	ip := pattern.Find(body)

	return string(ip), nil
}

func GetExternalIPv4Address() (string, error) {
	body, err := fetchPage()
	if err != nil {
		return "", err
	}

	ip, err := verifyIPv4Addr(body)
	if err != nil {
		return "", err
	}
	return ip, nil
}
