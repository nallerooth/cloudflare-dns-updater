package email

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nallerooth/cloudflare-dns-updater/config"
)

// Email section of config.json
var dummyConfig = []byte(`
{
	"notifications": {
		"email": {
			"smtp": {
				"server": "abc.email.com",
				"port": 587,
				"username": "useruser",
				"password": "passpass"
			},
			"sender": "sender@abc.email.com",
			"recipients": [
				"recipient.a@abc.email.com",
				"recipient.b@abc.email.com",
				"recipient.c@abc.email.com"
			]
		}
	}
}
`)

func TestEmail(t *testing.T) {
	c, err := config.LoadJSON(dummyConfig)
	if err != nil {
		t.Errorf("Loading of config failed: %s", err)
	}

	t.Run("getAddrFromConfig", func(t *testing.T) {
		actual, err := getAddrFromConfig(c)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		expected := "abc.email.com:587"

		if strings.Compare(expected, actual) != 0 {
			t.Errorf("Strings do not match. Expected: '%s', got: '%s'", expected, actual)
		}

		// Create empty conf
		badConf, _ := config.LoadJSON([]byte(`{}`))

		// Missing server address
		_, err = getAddrFromConfig(badConf)
		if err.Error() != "Missing SMTP server address" {
			t.Fail()
		}
		// Set server address
		badConf.Notifications.Email.SMTP.Server = "abc.email.com"

		// Still missing server port
		_, err = getAddrFromConfig(badConf)
		if err.Error() != "Missing SMTP server port" {
			t.Fail()
		}
	})

	t.Run("getSenderFromConfig", func(t *testing.T) {
		actual, err := getSenderFromConfig(c)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		expected := "sender@abc.email.com"

		if strings.Compare(expected, actual) != 0 {
			t.Errorf("Strings do not match. Expected: '%s', got: '%s'", expected, actual)
		}

		// Create empty conf
		badConf, _ := config.LoadJSON([]byte(`{}`))
		_, err = getSenderFromConfig(badConf)
		if err.Error() != "Missing sender email" {
			t.Fail()
		}
	})

	t.Run("getRecipientsFromConfig", func(t *testing.T) {
		actual, err := getRecipientsFromConfig(c)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		expected := []string{
			"recipient.a@abc.email.com",
			"recipient.b@abc.email.com",
			"recipient.c@abc.email.com",
		}

		if len(actual) != 3 {
			t.Errorf("Unexpected number of recipients. Wanted %d, got %d", 3, len(actual))
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Recipients do not match. Expected: '%v', got: '%v'", expected, actual)
		}

		// Create empty conf
		badConf, _ := config.LoadJSON([]byte(`{}`))
		_, err = getRecipientsFromConfig(badConf)
		if err.Error() != "Missing recipients" {
			t.Fail()
		}
	})
}
