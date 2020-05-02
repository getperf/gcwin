package exporter

type Command struct {
	Level int
	Id    string
	Text  string
}

type Env struct {
	Level     int
	DryRun    bool
	Datastore string
	NodeDir   string
}
