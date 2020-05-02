package config

import "time"

type Job struct {
	Id            int    `toml:"id"`   /**< primary key(sequence) */
	Name          string `toml:"name"` /**< Metric */
	Status        string
	Enable        bool `toml:"enable"`     /**< Enabled */
	Build         int  `toml:"build"`      /**< Build version */
	StdoutLog     bool `toml:"stdout_log"` /**< Standard output flag */
	Interval      int  `toml:"interval"`   /**< Interval(sec) */
	Timeout       int  `toml:"timeout"`    /**< Timeout(sec) */
	NextTimestamp time.Time
	Mode          string /**< Stataus mode */
}

func NewJob(statName string) *Job {
	job := Job{
		Name:          statName,
		NextTimestamp: time.Now(),
	}
	return &job
}
