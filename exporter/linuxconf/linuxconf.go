package linuxconf

import (
	"fmt"
	"path/filepath"

	"github.com/getperf/gcagent/exporter"
	log "github.com/sirupsen/logrus"
)

type Linux struct {
	Server   string `toml:"server"`
	IsRemote bool   `toml:"is_remote"`
	Ip       string `toml:"ip"`
	UserId   string `toml:"user_id"`
	User     string `toml:"specific_user`
	Password string `toml:"specific_password`
}

var sampleScheduleConfig = `
  [jobs.linuxconf]
  enable = true
  interval = 300
  timeout = 340
`

var sampleAccountConfig = `
# Linux account settings
# Enter account information for OS general users
# Please specify the userid key in [[accounts.{User Id}]]
# 
# example:
#    [[accounts.admin02]]
#    user = "someuser"
#    password = "P@ssword"
`

var sampleConfig = `
# Required Linux Endpoint
# 
server = "{{ .Server }}"
is_remote = {{ .IsRemote }}
ip = "{{ .Ip }}"
user_id = "{{ .UserId }}"
specific_user = "{{ .User }}"
specific_password = "{{ .Password }}"
`

var commands = []exporter.Command{
	{Level: 0, Id: "hostname", Text: "hostname -s"},
}

func (e *Linux) Description() string {
	return "Gather Linux inventorys"
}

func (e *Linux) SampleScheduleConfig() string {
	return sampleScheduleConfig
}

func (e *Linux) SampleAccountConfig() string {
	return sampleAccountConfig
}

func (e *Linux) SampleConfig() string {
	return sampleConfig
}

func (target *Linux) Setup() {
	if target.IsRemote {
		log.Info("create ssh session")
	}
}

func (e *Linux) Run(env *exporter.Env) error {
	fmt.Printf("run '%s' through linux platform\n", e.Server)
	fmt.Printf("env [%v]\n", env)
	for _, command := range commands {
		if command.Level > env.Level {
			break
		}
		// if env.
		log.Info("run command ", command.Id)
		c := CommandInfo{
			CmdLine: command.Text,
			OutPath: filepath.Join(env.Datastore, command.Id),
			Timeout: 30,
		}
		c.ExecCommandRedirect()
	}
	return nil
}

func init() {
	exporter.Add("linuxconf", func() exporter.Exporter {
		return &Linux{
			Server: "localhost",
		}
	})
}
