package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cf "nallerooth.com/cloudflare"
	"nallerooth.com/iplookup"
)

func loadConfig(filename string) (*cf.Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	c := &cf.Config{}
	err = json.Unmarshal(contents, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func compareIPAddrs(ip string, entry *cf.DNSEntry) bool {
	for _, record := range entry.Result {
		if strings.Compare(ip, record.Content) == 0 {
			return true
		}
	}

	return false
}

func main() {
	cfConfig, err := loadConfig("./config/config.json")
	if err != nil {
		log.Fatalln("Unable to load config file")
	}

	ip, err := iplookup.GetExternalIPv4Address()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get external IP address: %s", err)
	}

	fmt.Println("External IP Address", ip)

	dnsEntries, err := cf.GetDNSEntry(cfConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}

	fmt.Printf("%+v\n", dnsEntries)

	if compareIPAddrs(ip, dnsEntries) {
		log.Println("Correct IP address in DNS record. Going to sleep..")
	} else {
		log.Println("Cloudflare IP does not match external IP -> Updating CF to", ip)
		cf.UpdateDNSEntry(cfConfig, ip)
	}
}
