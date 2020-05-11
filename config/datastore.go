package config

import (
	"path/filepath"
	"strings"
	"time"

	. "github.com/getperf/gcagent/common"
)

type ApplicationType int

const (
	INVENTORY ApplicationType = iota
	PERFMETRIC
	ETC
)

type JobStatus int

const (
	SUCCESS JobStatus = iota
	WARN
	ERROR
)

type Datastore struct {
	OutDir   string
	JobName  string
	Host     string
	AppType  ApplicationType
	Now      time.Time
	Interval int

	Status    JobStatus /**< Process status */
	StartTime time.Time /**< Start time(UTC) */
	EndTime   time.Time /**< End time(UTC) */
	Messages  string
}

func getAppFromJob(jobName string) ApplicationType {
	if strings.HasSuffix(jobName, "conf") {
		return INVENTORY
	} else if strings.HasSuffix(jobName, "perf") {
		return PERFMETRIC
	} else {
		return ETC
	}
}

func NewDatastoreBase(outDir, jobName, host string, now time.Time, interval int) *Datastore {
	ds := Datastore{
		OutDir:    outDir,
		JobName:   jobName,
		Host:      host,
		AppType:   getAppFromJob(jobName),
		Now:       now,
		Interval:  interval,
		StartTime: time.Now(),
	}
	return &ds
}

func (c *Config) NewDatastore(jobName, host string, now time.Time) *Datastore {
	interval := 24 * 3600
	job, ok := c.Jobs[jobName]
	if ok {
		interval = job.Interval
	}
	return NewDatastoreBase(c.OutDir, jobName, host, now, interval)
}

func (ds *Datastore) Path() string {
	return filepath.Join(ds.OutDir, ds.RelDir())
}

func (ds *Datastore) RelDir() string {
	start := ds.Now
	interval := time.Duration(ds.Interval)
	start = start.Truncate(time.Second * interval)
	dateDir := GetTimeString(YYYYMMDD, start)
	timeDir := GetTimeString(HHMISS, start)
	return filepath.Join(ds.Host, ds.JobName, dateDir, timeDir)
}

func (ds *Datastore) TargetDir() string {
	switch ds.AppType {
	case INVENTORY:
		return filepath.Join(ds.Host, ds.JobName)
	case PERFMETRIC:
		return ds.RelDir()
	default:
		return ds.RelDir()
	}
}

func (ds *Datastore) SetJobInfo(status JobStatus, message string, endTime time.Time) {
	ds.Status = status
	ds.Messages = message
	ds.EndTime = endTime
}
