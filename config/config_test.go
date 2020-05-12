package config

import (
	"testing"
	// . "github.com/getperf/gcagent/common"
	"github.com/BurntSushi/toml"
)

func TestNewConfigEnv(t *testing.T) {
	configEnv := NewConfigEnv()
	t.Log(configEnv)
	if configEnv == nil {
		t.Error("new config")
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	t.Log(config.TemplateDir)
	if config.Home != "../testdata/ptune/" {
		t.Error("check home")
	}
	if config.LogLevel != "warn" {
		t.Error("check log level")
	}
	for _, createDir := range config.GetBaseDirs() {
		t.Log(*createDir)
	}
}

func TestLoadConfig(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	cfgFile := "../testdata/ptune/gcagent.toml"
	_, err := toml.DecodeFile(cfgFile, &config)
	if err != nil {
		t.Error("read in config ", err)
	}
	config.CheckConfig()
	t.Log(config.Jobs["linuxconf"])
	if config.Jobs["linuxconf"].Enable != true {
		t.Error("read job parameter ", err)
	}
}

func TestLoadNG(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	cfgFile := "../testdata/ptune/gcagent_bad.toml"
	_, err := toml.DecodeFile(cfgFile, &config)
	if err == nil {
		t.Error("read bad config")
	}
}

func TestCheckConfig(t *testing.T) {
	config := NewConfig("/", NewConfigEnv())
	toml.DecodeFile("../testdata/ptune/gcagent.toml", &config)
	err := config.CheckConfig()
	if err != nil {
		t.Error("check config")
	}
}
