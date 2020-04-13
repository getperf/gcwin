package agent

import "time"

type Collector struct {
	Id            int    `mapstructure:"id"`        /**< primary key(sequence) */
	StatName      string `mapstructure:"stat_name"` /**< Metric */
	Status        string
	StatEnable    bool `mapstructure:"stat_enable"`     /**< Enabled */
	Build         int  `mapstructure:"build"`           /**< Build version */
	StatStdoutLog bool `mapstructure:"stat_stdout_log"` /**< Standard output flag */
	StatInterval  int  `mapstructure:"stat_interval"`   /**< Interval(sec) */
	StatTimeout   int  `mapstructure:"stat_timeout"`    /**< Timeout(sec) */
	NextTimestamp time.Time
	StatMode      string /**< Stataus mode */

	Jobs []*Job /**< First job */
}

func NewCollector(statName string) *Collector {
	collector := Collector{
		StatName:      statName,
		NextTimestamp: time.Now(),
	}
	return &collector
}

func (collector *Collector) AddJob(job *Job) {
	collector.Jobs = append(collector.Jobs, job)
}

// func (collector *Collector)FindOrCreateJob(cmd string) *Job {
//     if job, ok := collector.Jobs[cmd]; ok {
//         log.Info("DBG found : ", *job)
//         return job
//     } else {
//         id := len(collector.Jobs) + 1
//         newJob := NewJob(id, cmd)
//         collector.AddJob(newJob)
//         log.Info("DBG new : ", newJob)
//         return newJob
//     }
// }
