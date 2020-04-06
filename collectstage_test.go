package gcwin

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewCollectStageBase(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	stage := NewCollectStageBase(home, "HW", time.Now())
	absDir := strings.TrimPrefix(stage.AbsDir(), home)
	t.Log("Dir : ", absDir)
	if len(absDir) == 0 {
		t.Error("new out log base is nil")
	}
	if TrimPathSeparator(absDir) != stage.RelDir() {
		t.Error("new out log base")
	}
}

func TestNewCollectStage(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	// config := NewConfig(home, NewConfigEnv())
	// config.InitAgent()
	// config.ParseConfigLine("STAT_INTERVAL.HW = 300")
	// config.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")

	date := "2020-01-20 11:00:00"
	now, _ := time.Parse("2006-01-02 15:04:05", date)

	// stage, err := config.NewCollectStage("HW", now)
	stage := NewCollectStageBase(home, "HW", now)
	// if err != nil {
	// 	t.Error("new out log", err)
	// }
	t.Log(stage.ZipFile("hoge"))
	if stage.ZipFile("hoge") != "arc_hoge__HW_20200120_110000.zip" {
		t.Error("new out log 2")
	}
}
