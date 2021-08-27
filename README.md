# Cloudflare DNS Updater
A simple tool for automatically updating Cloudflare DNS entries on external IP changes, using the Cloudflare API. This is useful when your ISP is unable to provide you with a static IP address.

It does this by talking to an external site which echoes the IP address the request was made from. There are a lot of such services available for free, but some of them respond with random IP addresses when queried without a "valid" browser User-Agent.

Please be a good netizen and find yourself a service that is okay with scripting, or host a script somewhere yourself.

## Building
This is really simple if you have a somewhat modern version of Go installed.
```bash
go build
```
This should result in a new binary named `cloudflare-dns-updater` in your current directory.

## Configuration
1. Create a new config file:
```
cp config.json.example config.json
```
2. Populate the configuration file according to your needs.
```json
{
    "cloudflare": {
        "api_token": "abc123abc123abc123abc123abc123abc123a",
        "email":    "your.email@domain.com",
        "zone_id":   "abc123abc123abc123abc123abc123ab",
        "zone_name": "your.domain.com"
    },
    "ip_lookup_url": "https://domain.com/whats-my-ip",
    "sleep_seconds": 60,
    "verbose_mode": false
}
```

### Notes on configuration
The details required for the Cloudflare API can be found in the [ Cloudflare API documentation ]( https://api.cloudflare.com/#getting-started-endpoints ).

| Key                    | Description                                                                                   |
|------------------------|-----------------------------------------------------------------------------------------------|
| `ip_lookup_url`        | The URL to query for your external IP. The first IPv4 address found on the page will be used. |
| `sleep_seconds`        | The sleep time between external IP checks. Only used if running in *watch* mode.              |
| `verbode_mode`         | Log more details to stdout. Can be enabled by the `-v` flag even if set to false.             |

## Flags
```
Usage of ./cloudflare-dns-updater:
  -c string
    	Path to config file (default "./config.json")
  -v	Verbose mode
  -w	Keep watching for external IP change
```
