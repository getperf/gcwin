package agent

import (
	"fmt"
	"path/filepath"
	"time"
)

type Datastore struct {
	StartTime time.Time
	StatName  string
	DateDir   string
	TimeDir   string
	OutDir    string
}

func NewDatastoreBase(outDir, statName string, start time.Time) *Datastore {
	ds := &Datastore{
		StartTime: start,
		StatName:  statName,
		DateDir:   GetTimeString(YYYYMMDD, start),
		TimeDir:   GetTimeString(HHMISS, start),
		OutDir:    outDir,
	}
	return ds
}

func (c *Config) NewDatastore(statName string, start time.Time) (*Datastore, error) {
	schedule := c.Schedule
	if schedule == nil {
		return nil, fmt.Errorf("new out log schedule not found")
	}
	collector := schedule.Collectors[statName]
	if collector == nil {
		return nil, fmt.Errorf("collector not found %s", statName)
	}
	interval := time.Duration(collector.StatInterval)
	start = start.Truncate(time.Second * interval)
	fmt.Printf("interval : %v, start :%v\n", interval, start)
	ds := NewDatastoreBase(c.OutDir, statName, start)
	return ds, nil
}

func (c *Config) NewDatastoreCurrent(statName string) (*Datastore, error) {
	return c.NewDatastore(statName, time.Now())
}

func (ds *Datastore) RelDir() string {
	return filepath.Join(ds.StatName, ds.DateDir, ds.TimeDir)
}

func (ds *Datastore) AbsDir() string {
	return filepath.Join(ds.OutDir, ds.RelDir())
}

func (ds *Datastore) ZipFile(host string) string {
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip",
		host, ds.StatName, ds.DateDir, ds.TimeDir,
	)
}

func (ds *Datastore) InventoryZipFile(host string) string {
	return fmt.Sprintf("inventory_%s__%s.zip",
		host, ds.StatName,
	)
}

func (ds *Datastore) OldZipFile(host string, hour int) string {
	start := ds.StartTime.Add(-1 * time.Hour * time.Duration(hour))
	dateDir := GetTimeString(YYYYMMDD, start)
	timeDir := GetTimeString(HHMISS, start)
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip", host, ds.StatName, dateDir, timeDir)
}
