package vmwareconf

import (
	"fmt"

	. "github.com/getperf/gcagent/exporter"
)

type VMWare struct {
	Server     string `toml:"server"`
	TemplateId string `toml:"user_id"`
	Url        string `toml:"url"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
}

var sampleScheduleConfig = `
  [jobs.vmwareconf]
  enable = true
  interval = 300
  timeout = 340
`

var sampleTemplateConfig = `
# VMWare template settings
# Enter template information for vCenter users
# 
# example:
#
#    url = "192.168.10.100"    # vCenter URL
#    user = "test_user"
#    password = "P@ssword"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
`

var sampleConfig = `
## Required VMWare Endpoint
# server = "localhost"
# export_level = 1
server = "{{ .Server }}"
template_id = "{{ .TemplateId }}"
url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
`

func (e *VMWare) Label() string {
	return "vCenter"
}

func (e *VMWare) Config(configType ConfigType) string {
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

func (e *VMWare) Setup(env *Env) error {
	fmt.Printf("export '%s' through VMWare platform\n", e.Server)
	return nil
}

func (e *VMWare) Run(env *Env) error {
	fmt.Printf("run '%s' through VMWare platform\n", e.Server)
	return nil
}

func init() {
	AddExporter("vmwareconf", func() Exporter {
		return &VMWare{
			Server: "localhost",
		}
	})
}
