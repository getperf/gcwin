package cmd

import (
	"testing"
)

func TestBootSetting(t *testing.T) {
	set, err := makeBootSettings("/tmp/hoge.toml")
	t.Log(err)
	t.Log(set)
}
