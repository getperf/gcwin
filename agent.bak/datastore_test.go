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
	ds := NewDatastoreBase("hoge", home, "HW", time.Now(), TIME_ON)
	absDir := strings.TrimPrefix(ds.PathAbs(), home)
	if len(absDir) == 0 {
		t.Error("new out log base is nil")
	}
	if TrimPathSeparator(absDir) != ds.Path() {
		t.Error("new out log base")
	}
}

func TestNewDatastoreBaseTimeOff(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	ds := NewDatastoreBase("hoge", home, "HW", time.Now(), TIME_OFF)
	if ds.Path() != "hoge/HW" {
		t.Error("time off datastore")
	}
}

func TestNewDatastoreRemote(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()

	_, err := config.NewDatastoreRemote("hoge", "HW", time.Now(), TIME_ON)
	if err != nil {
		t.Error("new out log", err)
	}
}

func TestNewDatastore(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	// config.ParseConfigLine("STAT_INTERVAL.HW = 300")
	// config.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")

	date := "2020-01-20 11:00:00"
	now, _ := time.Parse("2006-01-02 15:04:05", date)

	ds, err := config.NewDatastore("HW", now, TIME_ON)
	if err != nil {
		t.Error("new out log", err)
	}
	if ds.ZipFile("hoge") != "arc_hoge__HW_20200120_110000.zip" {
		t.Error("new out log 2")
	}
}

// func TestInventoryZipFile(t *testing.T) {
// 	home, _ := ioutil.TempDir("", "ptune")
// 	defer os.RemoveAll(home)
// 	ds := NewDatastoreBase("hoge", home, "HW", time.Now(), TIME_ON)
// 	zip := ds.InventoryZipFile("hoge")
// 	t.Log(zip)
// 	if len(zip) == 0 {
// 		t.Error("new inventory zip is nil")
// 	}
// }
