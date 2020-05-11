package cmd

import (
	"os"
	"path/filepath"

	. "github.com/getperf/gcagent/common"
	"github.com/getperf/gcagent/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const ConfigName = "gcagent.toml"

type BootParameters struct {
	configPath string
	home       string
}

func (params *BootParameters) make(configFile string) error {
	// -c gcagent.toml で設定ファイルの指定がある場合、ファイルの有無
	// をチェック。ファイルが存在しなければエラーを返す。存在する場合は
	// ベースディレクトリをホームにセット
	if configFile != "" {
		if ok, err := CheckFile(configFile); !ok {
			return errors.Wrap(err, "make boot settings")
		}
		configPath, _ := filepath.Abs(configFile)
		params.configPath = configPath
		params.home = filepath.Dir(configPath)

		// -c gcagent.toml 指定がない場合、カレントディレクトリを
		// ホームにして設定ファイルの有無をチェック。存在する場合は、
		// カレントディレクトリをホームにセット
	} else {
		currentDir, _ := os.Getwd()
		configFile := filepath.Join(currentDir, ConfigName)
		if ok, _ := CheckFile(configFile); ok {
			params.configPath = configFile
			params.home = currentDir
		}
	}
	return nil
}

func (params *BootParameters) NewConfig() *config.Config {
	hostName, err := GetHostname()
	if err != nil {
		log.Errorf("get hostname for initialize config %s", err)
		hostName = "UnkownHost"
	}
	configEnv := config.NewConfigEnvBase(hostName, params.configPath)
	return config.NewConfig(params.home, configEnv)
}
