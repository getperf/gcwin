package config

import (
	"fmt"
	"path/filepath"
)

type Template struct {
	UserId   string `toml:"user_id"`
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

func NewTemplate(id, user, password string) *Template {
	return &Template{
		UserId:   id,
		User:     user,
		Password: password,
	}
}

func (c *Config) TemplateConfig(job string) string {
	filename := fmt.Sprintf("%s.toml", job)
	return filepath.Join(c.TemplateDir, filename)
}
