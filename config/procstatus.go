package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type ProcessInfo struct {
	Pid int
}

func NewProcessInfo() *ProcessInfo {
	return &ProcessInfo{}
}

func (c *Config) WriteProcInfo() error {
	p := ProcessInfo{
		Pid: c.Pid,
	}
	bytes, _ := json.Marshal(&p)
	return ioutil.WriteFile(c.PidPath, bytes, os.FileMode(0600))
}

func (c *Config) ReadProcInfo() (ProcessInfo, error) {
	p := ProcessInfo{}
	file, err := ioutil.ReadFile(c.PidPath)
	if err != nil {
		return p, errors.Wrap(err, "read proc info")
	}
	err = json.Unmarshal(file, &p)
	return p, err
}
