package config

import (
	"testing"
)

func TestServerConfig(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	server := config.ServerConfig("windowsconf", "win2016")
	if server == "" {
		t.Error("get server config")
	}
}

func TestServerConfigs(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	servers, err := config.ServerConfigs("windowsconf")
	t.Log("windowsconf size : ", len(servers))
	if err != nil || len(servers) == 0 {
		t.Error("windows configs ", err)
	}
	servers, err = config.ServerConfigs("linuxconf")
	if err != nil || len(servers) == 0 {
		t.Error("linux configs ", err)
	}
}
