package linuxconf

import (
	"fmt"
	"path/filepath"

	. "github.com/getperf/gcagent/exporter"
	log "github.com/sirupsen/logrus"
)

type Linux struct {
	Server     string `toml:"server"`
	IsRemote   bool   `toml:"is_remote"`
	Ip         string `toml:"ip"`
	TemplateId string `toml:"user_id"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
}

var sampleScheduleConfig = `
  [jobs.linuxconf]
  enable = true
  local_exec = true
  interval = 300
  timeout = 340
`

var sampleTemplateConfig = `
# Linux template settings
# Enter template information for OS general users
# 
# example:
#    user = "someuser"
#    password = "P@ssword"

user = "{{ .User }}"
password = "{{ .Password }}"
`

var sampleConfig = `
# Required Linux Endpoint
# 
server = "{{ .Server }}"
is_remote = {{ .IsRemote }}
ip = "{{ .Ip }}"
template_id = "{{ .TemplateId }}"
user = "{{ .User }}"
password = "{{ .Password }}"
`

var commands = []Command{
	{Level: 0, Id: "hostname", Text: "hostname -s"},
}

func (e *Linux) Label() string {
	return "Linux"
}

func (e *Linux) Config(configType ConfigType) string {
	switch configType {
	case SCHEDULE:
		return sampleScheduleConfig
	case TEMPLATE:
		return sampleTemplateConfig
	case SERVER:
		return sampleConfig
	default:
		return ""
	}
}

func (target *Linux) Setup(env *Env) error {
	return nil
}

func (e *Linux) Run(env *Env) error {
	fmt.Printf("run '%s' through linux platform\n", e.Server)
	fmt.Printf("env [%v]\n", env)
	for _, command := range commands {
		if command.Level > env.Level {
			break
		}
		if env.DryRun {
			log.Infof("command[%s] : %s", command.Id, command.Text)
		} else {
			c := CommandInfo{
				CmdLine: command.Text,
				OutPath: filepath.Join(env.Datastore, command.Id),
				Timeout: 30,
			}
			c.ExecCommandRedirect()
		}
	}
	return nil
}

func init() {
	AddExporter("linuxconf", func() Exporter {
		return &Linux{
			Server: "localhost",
		}
	})
}
