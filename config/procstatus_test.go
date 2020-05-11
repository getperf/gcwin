package config

import (
	"io/ioutil"
	"os"

	"testing"
)

func createHomeAndConfig() *Config {
	tmpdir, _ := ioutil.TempDir("", "ptune")
	config := NewConfig(tmpdir, NewConfigEnv())
	os.MkdirAll(config.WorkCommonDir, 0777)
	return config
}

func TestWriteProcInfo(t *testing.T) {
	c := createHomeAndConfig()
	defer os.Remove(c.Home)

	c.Pid = 12345
	c.WriteProcInfo()
	if stat, err := c.ReadProcInfo(); stat.Pid != 12345 || err != nil {
		t.Error("write pid")
	}
	os.Remove(c.PidPath)
	if stat, err := c.ReadProcInfo(); stat.Pid != 0 || err == nil {
		t.Error("write pid")
	}
}
