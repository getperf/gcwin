package exporter

type Command struct {
	Level int    `toml:"level"`
	Id    string `toml:"id"`
	Text  string `toml:"text"`
}

type Env struct {
	Level     int
	DryRun    bool
	Datastore string
	LocalExec bool
	Messages  string
	ErrMsgs   string

	// バッチ用
	TemplateConfigs map[string]string
	ServerConfigs   map[string]string
}
