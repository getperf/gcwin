package cmd

import (
	"testing"
)

func TestBootSetting(t *testing.T) {
	var params BootParameters

	err := params.make("")
	if err != nil || params.home != "" {
		t.Error("make boot param null")
	}

	err = params.make("../testdata/ptune/gcagent.toml")
	if err != nil || params.home == "" {
		t.Error("make boot param")
	}
}
