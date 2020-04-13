package agent

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewDatastoreBase(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	ds := NewDatastoreBase(home, "HW", time.Now())
	absDir := strings.TrimPrefix(ds.AbsDir(), home)
	if len(absDir) == 0 {
		t.Error("new out log base is nil")
	}
	if TrimPathSeparator(absDir) != ds.RelDir() {
		t.Error("new out log base")
	}
}

func TestNewDatastore(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	config.ParseConfigLine("STAT_INTERVAL.HW = 300")
	config.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")

	date := "2020-01-20 11:00:00"
	now, _ := time.Parse("2006-01-02 15:04:05", date)

	ds, err := config.NewDatastore("HW", now)
	if err != nil {
		t.Error("new out log", err)
	}
	if ds.ZipFile("hoge") != "arc_hoge__HW_20200120_110000.zip" {
		t.Error("new out log 2")
	}
}

func TestInventoryZipFile(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	ds := NewDatastoreBase(home, "HW", time.Now())
	zip := ds.InventoryZipFile("hoge")
	t.Log(zip)
	if len(zip) == 0 {
		t.Error("new inventory zip is nil")
	}
}
