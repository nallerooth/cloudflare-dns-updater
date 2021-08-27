package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	cf "github.com/nallerooth/cloudflare-dns-updater/cloudflare"
	"github.com/nallerooth/cloudflare-dns-updater/config"
	"github.com/nallerooth/cloudflare-dns-updater/iplookup"
)

type options struct {
	configPath string
	verbose    bool
	watch      bool
}

func getLaunchFlags() options {
	o := options{}

	flag.StringVar(&o.configPath, "c", "./config.json", "Path to config file")
	flag.BoolVar(&o.watch, "w", false, "Keep watching for external IP change")
	flag.BoolVar(&o.verbose, "v", false, "Verbose mode")
	flag.Parse()

	return o
}

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

func run(conf *config.Config) error {
	ip, err := iplookup.GetExternalIPv4Address(conf)
	if err != nil {
		return fmt.Errorf("Unable to get external IP address: %s", err)
	}

	if conf.VerboseMode {
		log.Println("External IP Address", ip)
	}

	dns, err := cf.GetDNSEntry(conf)
	if err != nil {
		return fmt.Errorf("Unable to fetch DNS entry from Cloudflare: %s", err)
	}

	if compareIPAddrs(ip, dns) {
		if conf.VerboseMode {
			log.Println("Cloudflare DNS record matches external IP address")
		}
	} else {
		log.Printf("Cloudflare IP (%s) does not match external IP -> Updating Cloudflare DNS record to %s", dns.Content, ip)
		dns.Content = ip
		success, err := cf.UpdateDNSEntry(dns, conf)
		if err != nil {
			return fmt.Errorf("Unable to update Cloudflare DNS record: %s", err)
		}

		if success {
			log.Println("Cloudflare DNS record successfully updated")
		} else {
			log.Println("The DNS record failed to update, but no error was given by Cloudflare")
		}
	}

	return nil
}

func main() {
	flags := getLaunchFlags()

	//fmt.Printf("%+v\n", flags) // TODO: Remove

	conf, err := loadConfig(flags.configPath)
	if err != nil {
		log.Fatalln("Unable to load config file")
	}
	// Allow override of verbose mode set in config file
	conf.VerboseMode = conf.VerboseMode || flags.verbose

	//fmt.Printf("%+v\n", conf) // TODO: Remove

	for {
		if err = run(conf); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}

		if !flags.watch {
			break
		}

		time.Sleep(time.Second * time.Duration(conf.SleepSeconds))
	}
}
