package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"nallerooth.com/config"
)

const baseURL = "https://api.cloudflare.com/client/v4/"

type APIError struct {
	Code    int
	Message string
}

func (ae APIError) String() string {
	return fmt.Sprintf("Cloudflare API error %d: %s", ae.Code, ae.Message)
}

type APIResponse struct {
	Success  bool
	Errors   []APIError
	Messages []string
}

type DNSEntry struct {
	APIResponse

	Result []struct {
		Content string
		Name    string
		Type    string
	}
}

func newRequest(method string, reqURL *url.URL, c *config.Config) *http.Request {
	req := &http.Request{
		Method: method,
		URL:    reqURL,
		Header: http.Header{},
	}
	req.Header.Add("X-Auth-Email", c.Email)
	req.Header.Add("X-Auth-Key", c.APIToken)
	req.Header.Add("ContentType", "application/json")

	return req
}

func GetDNSEntry(c *config.Config) (*DNSEntry, error) {
	path := fmt.Sprintf("zones/%s/dns_records?name=%s", c.ZoneID, c.ZoneName)

	apiURL, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, err
	}

	req := newRequest("GET", apiURL, c)
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	responseText, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	entry := &DNSEntry{}
	err = json.Unmarshal(responseText, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func UpdateDNSEntry(c *config.Config, ip string) (bool, error) {
	return false, nil
}
