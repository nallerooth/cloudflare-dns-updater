package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cf "nallerooth.com/cloudflare"
	"nallerooth.com/config"
	"nallerooth.com/iplookup"
)

func loadConfig(filename string) (*config.Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	c := &config.Config{}
	err = json.Unmarshal(contents, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func compareIPAddrs(ip string, dns *cf.DNSRecordDetails) bool {
	return strings.Compare(ip, dns.Content) == 0
}

func main() {
	conf, err := loadConfig("./config.json")
	if err != nil {
		log.Fatalln("Unable to load config file")
	}

	ip, err := iplookup.GetExternalIPv4Address(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get external IP address: %s", err)
	}

	log.Println("External IP Address", ip)

	dns, err := cf.GetDNSEntry(conf)
	if err != nil {
		log.Fatalln(err)
	}

	if compareIPAddrs(ip, dns) {
		log.Println("Correct IP address in DNS record. Going to sleep..")
	} else {
		log.Printf("Cloudflare IP (%s) does not match external IP -> Updating CF to %s", dns.Content, ip)
		dns.Content = ip
		success, err := cf.UpdateDNSEntry(dns, conf)
		if err != nil {
			log.Fatalln(err)
		}

		if success {
			log.Println("Cloudflare DNS record successfully updated")
		} else {
			log.Println("The DNS record failed to update, but no error was given by Cloudflare")
		}
	}
}
