package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	. "github.com/getperf/gcagent/common"
	"github.com/getperf/gcagent/config"
	"github.com/getperf/gcagent/exporter"
	_ "github.com/getperf/gcagent/exporter/all"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Project struct {
	Home string
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

func initHomeDirectory(c *config.Config) error {
	createDirs := c.GetBaseDirs()
	for _, createDir := range createDirs {
		if err := os.MkdirAll(*createDir, 0777); err != nil {
			return fmt.Errorf("initialize agent directory : %s", err)
		}
	}
	return nil
}

func initAccountConfig(c *config.Config) error {
	log.Info("init account config ", c.AccountDir)
	return nil
}

func (project *Project) Create() error {
	home := project.Home
	log.Info("creating project directory ", home)
	if ok, _ := CheckDirectory(home); ok {
		return fmt.Errorf("project exist : %s", home)
	}
	c := config.NewConfig(home, config.NewConfigEnv())
	if err := initHomeDirectory(c); err != nil {
		return errors.Wrap(err, "failed to initialize")
	}
	contents := []byte(project.SampleConfig())
	if err := ioutil.WriteFile(c.ConfigPath, contents, 0666); err != nil {
		return errors.Wrap(err, "write config")
	}
	return nil
}

func (project *Project) Add(job string, si *config.Server) error {
	home := project.Home
	// 対象ジョブのエクスポーターを取得
	exporter, ok := exporter.Exporters[job]
	if !ok {
		return fmt.Errorf("invalid exporter job : %s. example : 'windowsconf'", job)
	}

	// 対象ジョブ用のコンフィグファイルテンプレートを取得
	text := exporter().SampleConfig()
	tpl, err := template.New("config").Parse(text)
	if err != nil {
		return errors.Wrap(err, "parse config template")
	}

	// コンフィグファイルを新規作成してオープン。既存のファイルがある場合は再作成します
	c := config.NewConfig(home, config.NewConfigEnv())
	nodePath := c.ServerConfig(job, si.Server)
	nodeFile, err := CreateAndOpenFile(nodePath)
	if err != nil {
		return errors.Wrap(err, "failed create config")
	}
	defer nodeFile.Close()

	// テンプレートからサーバコンフィグファイル生成
	if err := si.FillInInfo(); err != nil {
		return errors.Wrap(err, "check server info")
	}
	err = tpl.Execute(nodeFile, si)
	if err != nil {
		return errors.Wrap(err, "render config template")
	}
	log.Info("create node config ", nodePath)
	return nil
}
