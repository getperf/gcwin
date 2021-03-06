package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/getperf/gcagent/common"
	"github.com/pkg/errors"
)

type Server struct {
	Server     string `toml:"server"`
	IsRemote   bool   `toml:"is_remote"`
	Url        string `toml:"url"`
	Ip         string `toml:"ip"`
	TemplateId string `toml:"template_id"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
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
	if s.IsRemote == true && s.TemplateId == "" && s.User == "" {
		return fmt.Errorf("template_id or user must specified")
	}
	return nil
}

func (c *Config) ServerConfigs(job string) (map[string]string, error) {
	configs := make(map[string]string, 100)
	servers, err := ioutil.ReadDir(c.NodeDir)
	// log.Info("NODEDIR : ", c.NodeDir)
	// log.Info("READDIR : ", err)
	if err != nil {
		return configs, errors.Wrap(err, "get configs")
	}
	for _, server := range servers {
		serverName := server.Name()
		filepath := c.ServerConfig(job, serverName)
		if ok, _ := CheckFile(filepath); ok {
			configs[serverName] = filepath
		}
	}
	return configs, nil
}

func (c *Config) LocalHostConfig(job string) string {
	return c.ServerConfig(job, c.Host)
}

func (c *Config) ServerConfig(job, server string) string {
	filename := fmt.Sprintf("%s.toml", job)
	return filepath.Join(c.NodeDir, server, filename)
}
