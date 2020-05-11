package windowsconf

import (
	. "github.com/getperf/gcagent/exporter"
)

type Windows struct {
	Server   string     `toml:"server"`
	Commands []*Command `toml:"commands"`

	ScriptPath string
}

var sampleScheduleConfig = `
  [jobs.windowsconf]
  enable = true
  local_exec = true
  interval = 86400
  timeout = 340
`

var sampleTemplateConfig = `
# When collecting the inventory of Windows platform, execute it locally.
# Therefore, no template setting is required
`

var sampleConfig = `
# Required Endpoint Windows server name.
server = "{{ .Server }}"

[[commands]]
# Describe the additional command list. Added to the default command list for
# Windows Inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# id = "osver"   # unique key
# level = 0      # command level [0-2]
# text = "Get-ItemProperty \"HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\" | FL"
`

func (e *Windows) Label() string {
	return "Windows"
}

func (e *Windows) Config(configType ConfigType) string {
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
