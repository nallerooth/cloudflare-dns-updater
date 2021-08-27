package iplookup

import (
	"strings"
	"testing"
)

func TestIPLookup(t *testing.T) {
	t.Run("verifyIPv4Addr", func(t *testing.T) {
		input := []string{
			"127.0.0.1",
			"Your IP address is 127.0.0.1",
			"Your IP address is 127.0.0.1. This is some padding.",
		}

		expected := "127.0.0.1"

		for _, str := range input {
			actual, err := verifyIPv4Addr([]byte(str))
			if err != nil {
				t.Error(err)
			}
			if strings.Compare(expected, actual) != 0 {
				t.Errorf("Strings do not match. Expected '%s', got '%s'", expected, actual)
			}
		}

		bad := []byte("123.123. asdf asdf")
		_, err := verifyIPv4Addr(bad)
		if err == nil {
			t.Error("Expected error but got none")
		}
	})
}
