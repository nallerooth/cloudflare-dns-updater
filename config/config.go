package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Cloudflare struct {
		APIToken string `json:"api_token"`
		Email    string `json:"email"`
		ZoneID   string `json:"zone_id"`
		ZoneName string `json:"zone_name"`
	}
	Notifications struct {
		Email struct {
			SMTP struct {
				Server   string
				Port     uint
				Username string
				Password string
			}
			Sender     string
			Recipients []string
		}
	}
	IPLookupURL  string `json:"ip_lookup_url"`
	SleepSeconds int    `json:"sleep_seconds"`
	VerboseMode  bool   `json:"verbose_mode"`
}

func LoadFromFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return LoadJSON(contents)
}

func LoadJSON(jsondata []byte) (*Config, error) {
	c := &Config{}
	err := json.Unmarshal(jsondata, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
