package windowsconf

type Windows struct {
	Server   string     `toml:"server"`
	Commands []*Command `toml:"commands"`

	ScriptPath string
}

var sampleScheduleConfig = `
  [jobs.windowsconf]
  enable = true
  interval = 86400
  timeout = 340
`

var sampleAccountConfig = `
# When collecting the inventory of Windows platform, execute it locally.
# Therefore, no account setting is required
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

func (e *Windows) Description() string {
	return "Gather Windows inventorys"
}

func (e *Windows) SampleScheduleConfig() string {
	return sampleScheduleConfig
}

func (e *Windows) SampleAccountConfig() string {
	return sampleAccountConfig
}

func (e *Windows) SampleConfig() string {
	return sampleConfig
}
