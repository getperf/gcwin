package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func (p *Project) initConfigFile(configPath string) error {
	if ok, _ := CheckFile(configPath); ok {
		log.Warnf("'%s' exist, Backup to '%s'", configPath, configPath+"_bak")
		if err := CopyFile(configPath, configPath+"_bak"); err != nil {
			return errors.Wrap(err, "backup config file")
		}
	}
	contents := []byte(p.SampleConfig())
	if err := ioutil.WriteFile(configPath, contents, 0664); err != nil {
		return errors.Wrap(err, "write config")
	}
	return nil
}

func (p *Project) initHomeDirectory(c *config.Config) error {
	createDirs := c.GetBaseDirs()
	for _, createDir := range createDirs {
		if err := os.MkdirAll(*createDir, 0777); err != nil {
			return fmt.Errorf("initialize agent directory : %s", err)
		}
	}
	if err := p.initConfigFile(c.ConfigPath); err != nil {
		return errors.Wrap(err, "initialize config")
	}
	return nil
}

func (p *Project) Create() error {
	home := p.Home
	c := config.NewConfig(home, config.NewConfigEnv())
	if ok, _ := CheckDirectory(home); ok {
		// 実行バイナリ保存ディレクトリ以外でディレクトリが存在する場合はエラーを返す
		homePath, _ := filepath.Abs(home)
		if homePath != c.BaseDir {
			return fmt.Errorf("'%s' exist, Please specify nonexistent directory.", home)
		}
	}
	if err := p.initHomeDirectory(c); err != nil {
		return errors.Wrap(err, "failed to initialize")
	}
	return nil
}

func (p *Project) Add(job string, si *config.Server, isTemplate bool) error {
	home := p.Home
	// 対象ジョブのエクスポーターを取得
	exp, ok := exporter.Exporters[job]
	if !ok {
		return fmt.Errorf("invalid exporter job : %s. example : 'windowsconf'", job)
	}

	// 対象ジョブ用のコンフィグファイルテンプレートを取得
	text := exp().Config(exporter.SERVER)
	if isTemplate {
		text = exp().Config(exporter.TEMPLATE)
		if si.TemplateId == "" {
			return fmt.Errorf("'--id' must specified to create template")
		}
	}
	tpl, err := template.New("config").Parse(text)
	if err != nil {
		return errors.Wrap(err, "creating config")
	}

	if err := si.FillInInfo(); err != nil {
		return errors.Wrap(err, "creating config")
	}

	// コンフィグファイルを新規作成してオープン。既存のファイルがある場合は再作成します
	c := config.NewConfig(home, config.NewConfigEnv())
	var nodePath string
	if isTemplate {
		nodePath = c.TemplateConfig(job, si.TemplateId)
	} else {
		nodePath = c.ServerConfig(job, si.Server)
	}
	nodeFile, err := CreateAndOpenFile(nodePath)
	if err != nil {
		return errors.Wrap(err, "creating config")
	}
	defer nodeFile.Close()

	// テンプレートからサーバコンフィグファイル生成
	err = tpl.Execute(nodeFile, si)
	if err != nil {
		return errors.Wrap(err, "creating config")
	}
	log.Info("config created : ", nodePath)
	return nil
}
