package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	// . "github.com/getperf/gcagent/common"
	"github.com/getperf/gcagent/config"
)

func createSampleProject(base string) (*Project, error) {
	home := filepath.Join(base, "home")
	project := NewProjectFromHome(home)
	err := project.Create()
	return project, err
}

func TestProjectConfig(t *testing.T) {
	project := &Project{}
	text := project.SampleConfig()
	t.Log(text)
	if text == "" {
		t.Error("sample config")
	}
}

func TestProjectCreate(t *testing.T) {
	base, _ := ioutil.TempDir("", "project")
	defer os.RemoveAll(base)
	project, err := createSampleProject(base)
	if err != nil {
		t.Error("project create")
	}
	configPath := filepath.Join(project.Home, "gcagent.toml")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		t.Error("config file create")
	}
}

func TestInitAccountConfig(t *testing.T) {
	base, _ := ioutil.TempDir("", "project")
	defer os.RemoveAll(base)
	project, _ := createSampleProject(base)
	c := config.NewConfig(project.Home, config.NewConfigEnv())
	err := project.initAccountConfig(c)
	if err != nil {
		t.Error("init account config ", err)
	}
}

func TestAddServerToProject(t *testing.T) {
	base, _ := ioutil.TempDir("", "project")
	// defer os.RemoveAll(base)
	project, _ := createSampleProject(base)
	server := config.NewServer("hoge")
	if err := project.Add("windowsconf", server); err != nil {
		t.Error("add windows server : ", err)
	}
	if err := project.Add("hoge", server); err == nil {
		t.Error("add hoge server")
	}
}
