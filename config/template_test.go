package config

import (
	"testing"
)

func TestTemplateConfig(t *testing.T) {
	config := NewConfig("../testdata/ptune/", NewConfigEnv())
	template := config.TemplateConfig("windowsconf")
	t.Log(template)
	if template == "" {
		t.Error("get server config")
	}
}
