package linuxconf

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/getperf/gcagent/exporter"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func TestLinuxDocs(t *testing.T) {
	linux := &Linux{}
	linux.Label()
	linux.Config(exporter.TEMPLATE)
}

func TestLinuxInit(t *testing.T) {
	exporter, ok := exporter.Exporters["linuxconf"]
	if !ok {
		t.Fatal("exporter not defined")
	}
	linux := exporter().(*Linux)
	t.Log(linux)
	assert.Equal(t, "localhost", linux.Server)
}

func TestLinuxConfig(t *testing.T) {
	exp := exporter.Exporters["linuxconf"]()
	linux := &Linux{}
	t.Log(exp.Config(exporter.TEMPLATE))
	t.Log(linux)
	// err := toml.Unmarshal([]byte(exporter.SampleConfig()), linux)
	// require.NoError(t, err)
	// t.Log(linux)
	t.Log(exp.(*Linux))
}

func TestLinuxUnkown(t *testing.T) {
	_, ok := exporter.Exporters["hoge"]
	assert.Equal(t, false, ok)
}

func TestLinuxLocal(t *testing.T) {
	exp := exporter.Exporters["linuxconf"]()
	s := `{"server" = "hogehoge", "is_remote" = false}`
	json.Unmarshal([]byte(s), exp)
	datastore, _ := ioutil.TempDir("", "datasotre")
	env := exporter.Env{Level: 0, DryRun: true, Datastore: datastore}
	t.Log(datastore)
	exp.Run(&env)
}

func TestLinuxRemote(t *testing.T) {
	exp := exporter.Exporters["linuxconf"]()
	s := `server = "ostrich2"
is_remote = true
ip = "127.0.0.1"
user_id = "guiest01"
user = "psadmin"
password = "psadmin"`
	toml.Unmarshal([]byte(s), exp)
	t.Log("JSON:", exp)
	datastore, _ := ioutil.TempDir("", "datasotre")
	env := exporter.Env{Level: 0, DryRun: true, Datastore: datastore}
	exp.Run(&env)
}
