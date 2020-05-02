package agent

import (
	"fmt"
	"path/filepath"
	"time"
)

type DatastoreMode int

const (
	TIME_ON DatastoreMode = iota
	TIME_OFF
)

type Datastore struct {
	StartTime time.Time
	Host      string
	StatName  string
	DateDir   string
	TimeDir   string
	OutDir    string
	Mode      DatastoreMode
}

func NewDatastoreBase(host, outDir, statName string, start time.Time, mode DatastoreMode) *Datastore {
	ds := &Datastore{
		StartTime: start,
		Host:      host,
		StatName:  statName,
		DateDir:   GetTimeString(YYYYMMDD, start),
		TimeDir:   GetTimeString(HHMISS, start),
		OutDir:    outDir,
		Mode:      mode,
	}
	return ds
}

// func (c *Config) NewDatastore(statName string, start time.Time) (*Datastore, error) {
// 	schedule := c.Schedule
// 	if schedule == nil {
// 		return nil, fmt.Errorf("new out log schedule not found")
// 	}
// 	job := schedule.Jobs[statName]
// 	if job != nil {
// 		interval := time.Duration(job.Interval)
// 		start = start.Truncate(time.Second * interval)
// 	}
// 	ds := NewDatastoreBase(c.Host, c.OutDir, statName, start)
// 	return ds, nil
// }

func (c *Config) NewDatastoreRemote(host, statName string, start time.Time, mode DatastoreMode) (*Datastore, error) {
	schedule := c.Schedule
	if schedule == nil {
		return nil, fmt.Errorf("new out log schedule not found")
	}
	job := schedule.Jobs[statName]
	if job != nil {
		interval := time.Duration(job.Interval)
		start = start.Truncate(time.Second * interval)
	}
	ds := NewDatastoreBase(host, c.OutDir, statName, start, mode)
	return ds, nil
}

func (c *Config) NewDatastore(statName string, start time.Time, mode DatastoreMode) (*Datastore, error) {
	host := c.GetServiceOrHostName()
	return c.NewDatastoreRemote(host, statName, start, mode)
}

func (ds *Datastore) Path() string {
	if ds.Mode == TIME_ON {
		return filepath.Join(ds.Host, ds.StatName, ds.DateDir, ds.TimeDir)
	} else {
		return filepath.Join(ds.Host, ds.StatName)
	}
}

func (ds *Datastore) PathAbs() string {
	return filepath.Join(ds.OutDir, ds.Path())
}

// func (ds *Datastore) RelDir() string {
// 	return filepath.Join(ds.StatName, ds.DateDir, ds.TimeDir)
// }

// func (ds *Datastore) AbsDir() string {
// 	return filepath.Join(ds.OutDir, ds.RelDir())
// }

func (ds *Datastore) ZipFile(host string) string {
	if ds.Mode == TIME_ON {
		return fmt.Sprintf("arc_%s__%s_%s_%s.zip",
			host, ds.StatName, ds.DateDir, ds.TimeDir,
		)
	} else {
		return fmt.Sprintf("arc_%s__%s.zip",
			host, ds.StatName,
		)
	}
}

func (ds *Datastore) OldZipFile(host string, hour int) string {
	start := ds.StartTime.Add(-1 * time.Hour * time.Duration(hour))
	dateDir := GetTimeString(YYYYMMDD, start)
	timeDir := GetTimeString(HHMISS, start)
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip", host, ds.StatName, dateDir, timeDir)
}
