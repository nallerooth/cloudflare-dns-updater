package iplookup

import (
	"io"
	"net/http"
	"regexp"

	"github.com/nallerooth/cloudflare-dns-updater/config"
)

func fetchPage(c *config.Config) ([]byte, error) {
	req, err := http.Get(c.IPLookupURL)
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

func GetExternalIPv4Address(c *config.Config) (string, error) {
	body, err := fetchPage(c)
	if err != nil {
		return "", err
	}

	ip, err := verifyIPv4Addr(body)
	if err != nil {
		return "", err
	}
	return ip, nil
}
