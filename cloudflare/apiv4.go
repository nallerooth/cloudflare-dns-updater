package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nallerooth/cloudflare-dns-updater/config"
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
	Count    int
	Errors   []APIError
	Messages []string
	Success  bool
}

type DNSRecordDetails struct {
	ID      string
	Content string
	Name    string
	Proxied bool
	TTL     int
	Type    string
}

type ZoneResponse struct {
	APIResponse
	Result []DNSRecordDetails
}

func newRequest(method string, reqURL *url.URL, data []byte, c *config.Config) (*http.Request, error) {
	req, err := http.NewRequest(method, reqURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	//req := &http.Request{
	//Method: method,
	//URL:    reqURL,
	//Header: http.Header{},
	//}
	req.Header.Add("X-Auth-Email", c.Cloudflare.Email)
	req.Header.Add("X-Auth-Key", c.Cloudflare.APIToken)
	req.Header.Add("ContentType", "application/json")

	return req, nil
}

func GetDNSEntry(c *config.Config) (*DNSRecordDetails, error) {
	path := fmt.Sprintf("zones/%s/dns_records?name=%s", c.Cloudflare.ZoneID, c.Cloudflare.ZoneName)

	apiURL, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, err
	}

	req, err := newRequest("GET", apiURL, nil, c)
	if err != nil {
		return nil, err
	}

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

	response := &ZoneResponse{}
	err = json.Unmarshal(responseText, response)
	if err != nil {
		return nil, err
	}

	// Request was successful but result was not
	if response.Success == false {
		errs := ""
		for _, e := range response.Errors {
			errs = errs + e.String()
		}
		return nil, errors.New(errs)
	}

	if len(response.Result) == 0 {
		return nil, errors.New("No DNS records in Cloudflare API response")
	}

	return &response.Result[0], nil
}

func UpdateDNSEntry(dns *DNSRecordDetails, c *config.Config) (bool, error) {
	path := fmt.Sprintf("zones/%s/dns_records/%s", c.Cloudflare.ZoneID, dns.ID)
	apiURL, err := url.Parse(baseURL + path)
	if err != nil {
		return false, err
	}

	payload, err := json.Marshal(dns)
	if err != nil {
		return false, err
	}

	req, err := newRequest("PUT", apiURL, payload, c)
	if err != nil {
		return false, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	responseText, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	response := &APIResponse{}
	err = json.Unmarshal(responseText, response)
	if err != nil {
		return false, err
	}

	if len(response.Errors) > 0 {
		errs := make([]string, 0, 3)
		for _, e := range response.Errors {
			errs = append(errs, e.String())
		}
		return false, errors.New(strings.Join(errs, " :: "))
	}

	return response.Success, nil
}
