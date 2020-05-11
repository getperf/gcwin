package golangconf

import (
	"fmt"

	. "github.com/getperf/gcagent/exporter"
)

type GoLang struct {
	Server string `toml:"server"`
}

var sampleScheduleConfig = `
  [jobs.golangconf]
  enable = true
  interval = 300
  timeout = 340
`

var sampleTemplateConfig = `
# When collecting the inventory of GoLang platform, execute it locally.
# Therefore, no template setting is required
`

var sampleConfig = `
  ## Required GoLang Endpoint
  # server = "localhost"
  # export_level = 1
`

func (e *GoLang) Label() string {
	return "Go"
}

func (e *GoLang) Config(configType ConfigType) string {
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

func (e *GoLang) Setup(env *Env) error {
	fmt.Printf("export '%s' through GoLang platform\n", e.Server)
	return nil
}

func (e *GoLang) Run(env *Env) error {
	fmt.Printf("run '%s' through GoLang platform\n", e.Server)
	return nil
}

func init() {
	AddExporter("golangconf", func() Exporter {
		return &GoLang{
			Server: "localhost",
		}
	})
}
