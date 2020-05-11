package config

import (
	"time"
	// . "github.com/getperf/gcagent/common"
)

// type Application int

// const (
// 	INVENTORY Application = iota
// 	PERFMETRIC
// 	ETC
// )

type Job struct {
	Id            int    `toml:"id"`   /**< primary key(sequence) */
	Name          string `toml:"name"` /**< Metric */
	Status        string
	Enable        bool `toml:"enable"`     /**< Enabled */
	Build         int  `toml:"build"`      /**< Build version */
	LocalExec     bool `toml:"local_exec"` /**< Enable local execution */
	Interval      int  `toml:"interval"`   /**< Interval(sec) */
	Timeout       int  `toml:"timeout"`    /**< Timeout(sec) */
	NextTimestamp time.Time
	Mode          string /**< Stataus mode */
}

func NewJob(jobName string) *Job {
	job := Job{
		Name:          jobName,
		NextTimestamp: time.Now(),
	}
	return &job
}

// func GetAppFromJob(jobName string) Application {
// 	if strings.HasSuffix(jobName, "conf") {
// 		return INVENTORY
// 	} else if strings.HasSuffix(jobName, "perf") {
// 		return PERFMETRIC
// 	} else {
// 		return ETC
// 	}
// }

// func (c *Config) Datastore(jobName string, now time.Time) string {
// 	dateDir := GetTimeString(YYYYMMDD, now)
// 	timeDir := GetTimeString(HHMISS, now)
// 	return filepath.Join(c.OutDir, c.Host, jobName, dateDir, timeDir)
// }

func (c *Config) LocalExec(jobName string) bool {
	job, ok := c.Jobs[jobName]
	if ok && job.LocalExec {
		return true
	} else {
		return false
	}
}
