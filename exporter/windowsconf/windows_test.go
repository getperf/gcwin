package windowsconf

import (
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/getperf/gcagent/exporter"
	"gopkg.in/yaml.v2"
)

var testNodeDir = "../../testdata/ptune/node"
var tomlpath = testNodeDir + "/win2016/windowsconf.toml"

// func init() {
// 	tomlpath = filepath.Join(testNodeDir, "win2016", "windowsconf.toml")
// }

func createTestEnv() *exporter.Env {
	datastore, _ := ioutil.TempDir("", "datasotre")
	env := &exporter.Env{
		Level:     0,
		DryRun:    false,
		Datastore: datastore,
		LocalExec: true,
	}
	return env
}

func TestWindowsNormal(t *testing.T) {
	exp := exporter.Exporters["windowsconf"]()
	env := createTestEnv()
	defer os.Remove(env.Datastore)
	if runtime.GOOS == "windows" {
		if err := exp.Run(env); err != nil {
			t.Error(err)
		}
	}
}

func TestWindowsToml(t *testing.T) {
	commands2 := Commands{Commands: commands}
	d, err := yaml.Marshal(commands2)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(d))
}

func TestWindowsConfig(t *testing.T) {
	var windows Windows
	// tomlpath := filepath.Join(testNodeDir, "win2016", "windowsconf.toml")
	_, err := toml.DecodeFile(tomlpath, &windows)
	if err != nil {
		t.Error(err)
	}
	t.Log(windows.Commands[0])
}

func TestWindowsInventoryCode(t *testing.T) {
	var windows Windows
	env := createTestEnv()
	defer os.Remove(env.Datastore)
	// tomlpath := filepath.Join(testNodeDir, "win2016.toml")
	_, err := toml.DecodeFile(tomlpath, &windows)
	if err != nil {
		t.Error(err)
	}
	stdout := new(bytes.Buffer)
	windows.writeScript(stdout, env)
	t.Log("Result: ", stdout.String())
}
