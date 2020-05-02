package exporter

type Exporter interface {
	// SampleScheduleConfig returns the default schedule configuration
	SampleScheduleConfig() string

	// SampleAccountConfig returns the default account configuration of servers
	SampleAccountConfig() string

	// SampleConfig returns the default configuration of the Exporter
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

	Setup()
	Run(env *Env) error
}

type Creator func() Exporter

var Exporters = map[string]Creator{}

func Add(name string, creator Creator) {
	Exporters[name] = creator
}
