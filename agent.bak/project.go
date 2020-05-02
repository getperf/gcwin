package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/getperf/gcagent/exporter"
	_ "github.com/getperf/gcagent/exporter/all"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Project struct {
	Home string
}

type ServerInfo struct {
	Server   string
	IsRemote bool
	Url      string
	Ip       string
	UserId   string
	User     string
	Password string
}

func NewProjectFromHome(home string) *Project {
	project := &Project{
		Home: home,
	}
	return project
}

func NewProject(configPath string) (*Project, error) {
	home, err := GetParentAbsPath(configPath, 1)
	if err != nil {
		return nil, errors.Wrap(err, "config path search")
	}
	project := &Project{
		Home: home,
	}
	return project, nil
}

func (project *Project) Create() error {
	home := project.Home
	log.Info("creating project directory ", home)
	if ok, _ := CheckDirectory(home); ok {
		return fmt.Errorf("project exist : %s", home)
	}
	config := NewConfig(home, NewConfigEnv())
	if err := config.InitHome(); err != nil {
		return errors.Wrap(err, "failed to initialize")
	}
	configPath := config.GetConfigPath()
	contents := []byte(project.SampleConfig())
	if err := ioutil.WriteFile(configPath, contents, 0666); err != nil {
		return errors.Wrap(err, "write config")
	}
	return nil
}

func (si *ServerInfo) FillInInfo() error {
	if si.Server == "" {
		si.Server, _ = GetHostname()
	}
	if si.Ip == "" && si.Url == "" {
		si.IsRemote = false
	} else {
		si.IsRemote = true
	}
	if si.IsRemote == true && si.UserId == "" {
		return fmt.Errorf("user_id must specified")
	}
	return nil
}

func (project *Project) Add(job string, si *ServerInfo) error {
	home := project.Home
	exporter, ok := exporter.Exporters[job]
	if !ok {
		return fmt.Errorf("invalid exporter job : %s. Set example : 'windowsconf'", job)
	}
	if err := si.FillInInfo(); err != nil {
		return errors.Wrap(err, "check server info")
	}

	text := exporter().SampleConfig()
	tpl, err := template.New("config").Parse(text)
	if err != nil {
		return errors.Wrap(err, "parse config template")
	}
	config := NewConfig(home, NewConfigEnv())
	nodeDir := filepath.Join(config.NodeDir, job)
	if err := os.MkdirAll(nodeDir, 0777); err != nil {
		return errors.Wrap(err, "create node directory")
	}

	nodePath := filepath.Join(nodeDir, fmt.Sprintf("%s.toml", si.Server))
	nodeFile, err := os.OpenFile(nodePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed create script")
	}
	defer nodeFile.Close()
	err = tpl.Execute(nodeFile, si)
	if err != nil {
		return errors.Wrap(err, "render config template")
	}
	log.Info("create node config ", nodePath)
	return nil
}
