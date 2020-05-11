package agent

type Exporter interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Label returns a one-sentence Label on the Input
	Label() string

	Setup()
	Run()
}
