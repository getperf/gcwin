package golangconf

import (
	"fmt"

	"github.com/getperf/gcagent/exporter"
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

var sampleAccountConfig = `
# When collecting the inventory of GoLang platform, execute it locally.
# Therefore, no account setting is required
`

var sampleConfig = `
  ## Required GoLang Endpoint
  # server = "localhost"
  # export_level = 1
`

func (e *GoLang) Description() string {
	return "Gather GoLang inventorys"
}

func (e *GoLang) SampleScheduleConfig() string {
	return sampleScheduleConfig
}

func (e *GoLang) SampleAccountConfig() string {
	return sampleAccountConfig
}

func (e *GoLang) SampleConfig() string {
	return sampleConfig
}

func (e *GoLang) Setup() {
	fmt.Printf("export '%s' through GoLang platform\n", e.Server)
}

func (e *GoLang) Run(env *exporter.Env) error {
	fmt.Printf("run '%s' through GoLang platform\n", e.Server)
	return nil
}

func init() {
	exporter.Add("golangconf", func() exporter.Exporter {
		return &GoLang{
			Server: "localhost",
		}
	})
}
