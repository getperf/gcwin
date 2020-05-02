package agent

type Exporter interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

	Setup()
	Run()
}
