package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/getperf/gcagent/common"
	"github.com/pkg/errors"
)

type Server struct {
	Server   string `toml:"server"`
	IsRemote bool   `toml:"is_remote"`
	Url      string `toml:"url"`
	Ip       string `toml:"ip"`
	UserId   string `toml:"user_id"`
	User     string `toml:"specific_user"`
	Password string `toml:"specific_password"`
}

func NewServer(server string) *Server {
	return &Server{
		Server: server,
	}
}

func (s *Server) FillInInfo() error {
	if s.Server == "" {
		s.Server, _ = GetHostname()
	}
	if s.Ip == "" && s.Url == "" {
		s.IsRemote = false
	} else {
		s.IsRemote = true
	}
	if s.IsRemote == true && s.UserId == "" {
		return fmt.Errorf("user_id must specified")
	}
	return nil
}

func (c *Config) ServerConfigs(job string) ([]string, error) {
	configs := make([]string, 0, 100)
	servers, err := ioutil.ReadDir(c.NodeDir)
	if err != nil {
		return configs, errors.Wrap(err, "get configs")
	}
	for _, server := range servers {
		filepath := c.ServerConfig(job, server.Name())
		if ok, _ := CheckFile(filepath); ok {
			configs = append(configs, filepath)
		}
	}
	return configs, nil
}

func (c *Config) ServerConfig(job, server string) string {
	filename := fmt.Sprintf("%s.toml", job)
	return filepath.Join(c.NodeDir, server, filename)
}
