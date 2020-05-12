package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/getperf/gcagent/common"
	"github.com/pkg/errors"
)

// type Template struct {
// 	UserId   string `toml:"user_id"`
// 	Url      string `toml:"url"`
// 	User     string `toml:"user"`
// 	Password string `toml:"password"`
// }

// func NewTemplate(id, user, password string) *Template {
// 	return &Template{
// 		UserId:   id,
// 		User:     user,
// 		Password: password,
// 	}
// }

func (c *Config) TemplateConfigs(job string) (map[string]string, error) {
	configs := make(map[string]string, 100)
	servers, err := ioutil.ReadDir(c.TemplateDir)
	if err != nil {
		return configs, errors.Wrap(err, "get configs")
	}
	for _, server := range servers {
		templateId := server.Name()
		filepath := c.TemplateConfig(job, templateId)
		if ok, _ := CheckFile(filepath); ok {
			configs[templateId] = filepath
		}
	}
	return configs, nil
}

func (c *Config) TemplateConfig(job string, templateId string) string {
	filename := fmt.Sprintf("%s.toml", job)
	return filepath.Join(c.TemplateDir, templateId, filename)
}
