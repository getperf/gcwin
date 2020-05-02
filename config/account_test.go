package config

import (
	"testing"
)

func TestAccountConfig(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	account := config.AccountConfig("windowsconf")
	t.Log(account)
	if account == "" {
		t.Error("get server config")
	}
}
