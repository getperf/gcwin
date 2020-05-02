package agent

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfigEnv(t *testing.T) {
	configEnv := NewConfigEnv()
	t.Log(configEnv)
	if configEnv == nil {
		t.Error("new config")
	}
}

func TestNewConfig(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	t.Log(config)
	t.Log(config.Schedule)
	if config == nil {
		t.Error("new config")
	}
}

// func TestAddCollector(t *testing.T) {
// 	schedule := NewSchedule()
// 	colLinux := NewJob("Linux")
// 	schedule.AddJob(colLinux)
// 	if schedule.Jobs["Linux"].Name != "Linux" {
// 		t.Error("add collector")
// 	}
// }

// func TestAddJob(t *testing.T) {
// 	colLinux := NewCollector("Linux")
// 	job := NewJob(1, "uname -a")
// 	colLinux.AddJob(job)
// 	if colLinux.Jobs[0].Cmd != "uname -a" {
// 		t.Error("add job")
// 	}
// }
