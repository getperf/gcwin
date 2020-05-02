package vmwareconf

import (
	"fmt"

	"github.com/getperf/gcagent/exporter"
)

type VMWare struct {
	Server string `toml:"server"`
}

var sampleScheduleConfig = `
  [jobs.vmwareconf]
  enable = true
  interval = 300
  timeout = 340
`

var sampleAccountConfig = `
# VMWare account settings
# Enter account information for vCenter users
# Please specify the user id key in [[accounts.{User ID}]]
# 
# example:
#    [[accounts.vmmanager01]]
#
#    url = "192.168.10.100"    # vCenter URL
#    user = "test_user"
#    password = "P@ssword"
`

var sampleConfig = `
  ## Required VMWare Endpoint
  # server = "localhost"
  # export_level = 1
`

func (e *VMWare) Description() string {
	return "Gather VMWare inventorys"
}

func (e *VMWare) SampleScheduleConfig() string {
	return sampleScheduleConfig
}

func (e *VMWare) SampleAccountConfig() string {
	return sampleAccountConfig
}

func (e *VMWare) SampleConfig() string {
	return sampleConfig
}

func (e *VMWare) Setup() {
	fmt.Printf("export '%s' through VMWare platform\n", e.Server)
}

func (e *VMWare) Run(env *exporter.Env) error {
	fmt.Printf("run '%s' through VMWare platform\n", e.Server)
	return nil
}

func init() {
	exporter.Add("vmwareconf", func() exporter.Exporter {
		return &VMWare{
			Server: "localhost",
		}
	})
}
