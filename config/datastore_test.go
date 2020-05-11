package config

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewDatastore(t *testing.T) {
	cfg := NewConfig("../testdata/ptune/", NewConfigEnv())
	now := time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)

	ds := cfg.NewDatastore("windowsconf", "hoge", now)

	t.Log(ds.AppType)
	if ds.AppType != INVENTORY {
		t.Error("datastore")
	}

	keyword := filepath.Join("20180102", "000000")
	if !strings.HasSuffix(ds.Path(), keyword) {
		t.Error("datastore")
	}

	relDir := filepath.Join("hoge", "windowsconf", "20180102", "000000")
	if ds.RelDir() != relDir {
		t.Error("datastore")
	}

	targetDir := filepath.Join("hoge", "windowsconf")
	if ds.TargetDir() != targetDir {
		t.Error("datastore")
	}
}

func TestNewDatastoreJob(t *testing.T) {
	cfg := NewConfig("../testdata/ptune/", NewConfigEnv())
	now := time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)

	ds := cfg.NewDatastore("windowsconf", "hoge", now)
	ds.SetJobInfo(ERROR, "error mesage", time.Now())
	if ds.Status != ERROR {
		t.Error("datastore job")
	}
}
