package main

import (
	"net/url"
	"testing"
)

// TestReadConfig do test
func TestConfigRead(t *testing.T) {
	config, err := readConfig("api.yml")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(config)
	if config.UID == "" {
		t.Log("UID is empty")
		t.Fail()
	}
}
