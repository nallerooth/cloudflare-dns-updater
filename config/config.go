package config

type Config struct {
	Cloudflare struct {
		APIToken string `json:"api_token"`
		Email    string `json:"email"`
		ZoneID   string `json:"zone_id"`
		ZoneName string `json:"zone_name"`
	}
	IPLookupURL  string `json:"ip_lookup_url"`
	SleepSeconds int    `json:"sleep_seconds"`
	VerboseMode  bool   `json:"verbose_mode"`
}
