package config

import (
	"fmt"
	"path/filepath"
)

type Account struct {
	UserId   string `toml:"user_id"`
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

func NewAccount(id, user, password string) *Account {
	return &Account{
		UserId:   id,
		User:     user,
		Password: password,
	}
}

func (c *Config) AccountConfig(job string) string {
	filename := fmt.Sprintf("%s.toml", job)
	return filepath.Join(c.AccountDir, filename)
}
