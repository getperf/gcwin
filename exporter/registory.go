package exporter

type ConfigType int

const (
	SCHEDULE ConfigType = iota
	TEMPLATE
	SERVER
)

type Exporter interface {
	// SampleConfig returns the default configuration of the Exporter
	Config(configType ConfigType) string

	// Label returns a one-sentence Label on the Input
	Label() string

	Setup(env *Env) error
	Run(env *Env) error
}

type Creator func() Exporter

var Exporters = map[string]Creator{}

func AddExporter(name string, creator Creator) {
	Exporters[name] = creator
}
