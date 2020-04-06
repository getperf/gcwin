package gcwin

import (
	"fmt"
	"path/filepath"
	"time"
)

type CollectStage struct {
	StartTime time.Time
	StatName  string
	DateDir   string
	TimeDir   string
	OutDir    string
}

func NewCollectStageBase(outDir, statName string, start time.Time) *CollectStage {
	outLog := &CollectStage{
		StartTime: start,
		StatName:  statName,
		DateDir:   GetTimeString(YYYYMMDD, start),
		TimeDir:   GetTimeString(HHMISS, start),
		OutDir:    outDir,
	}
	return outLog
}

func (c *Config) NewCollectStage(statName string, start time.Time) (*CollectStage, error) {
	schedule := c.Schedule
	if schedule == nil {
		return nil, fmt.Errorf("new out log schedule not found")
	}
	collector := schedule.Collectors[statName]
	if collector == nil {
		return nil, fmt.Errorf("new out log collector not found %s", statName)
	}
	interval := time.Duration(collector.StatInterval)
	start = start.Truncate(time.Second * interval)
	outLog := NewCollectStageBase(c.OutDir, statName, start)
	return outLog, nil
}

func (c *Config) NewCollectStageCurrent(statName string) (*CollectStage, error) {
	return c.NewCollectStage(statName, time.Now())
}

func (o *CollectStage) RelDir() string {
	return filepath.Join(o.StatName, o.DateDir, o.TimeDir)
}

func (o *CollectStage) AbsDir() string {
	return filepath.Join(o.OutDir, o.RelDir())
}

func (o *CollectStage) ZipFile(host string) string {
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip",
		host, o.StatName, o.DateDir, o.TimeDir,
	)
}

func (o *CollectStage) InventoryZipFile(host string) string {
	return fmt.Sprintf("arc_%s__%s.zip", host, o.StatName)
}

func (o *CollectStage) OldZipFile(host string, hour int) string {
	start := o.StartTime.Add(-1 * time.Hour * time.Duration(hour))
	dateDir := GetTimeString(YYYYMMDD, start)
	timeDir := GetTimeString(HHMISS, start)
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip", host, o.StatName, dateDir, timeDir)
}
