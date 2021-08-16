package config

type Config struct {
	APIToken    string `json:"api_token"`
	Email       string `json:"email"`
	ZoneID      string `json:"zone_id"`
	ZoneName    string `json:"zone_name"`
	IPLookupURL string `json:"ip_lookup_url"`
}
