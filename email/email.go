package email

import (
	"errors"
	"fmt"

	"github.com/nallerooth/cloudflare-dns-updater/config"
)

func getAddrFromConfig(c *config.Config) (string, error) {
	addr := c.Notifications.Email.SMTP.Server
	if len(addr) == 0 {
		return "", errors.New("Missing SMTP server address")
	}

	port := c.Notifications.Email.SMTP.Port
	if port == 0 {
		return "", errors.New("Missing SMTP server port")
	}

	return fmt.Sprintf("%s:%d", addr, port), nil
}

func getSenderFromConfig(c *config.Config) (string, error) {
	sender := c.Notifications.Email.Sender
	if len(sender) == 0 {
		return "", errors.New("Missing sender email")
	}

	return sender, nil
}

func getRecipientsFromConfig(c *config.Config) ([]string, error) {
	recipients := c.Notifications.Email.Recipients
	if len(recipients) == 0 || recipients == nil {
		return nil, errors.New("Missing recipients")
	}

	return recipients, nil
}
