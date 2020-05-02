package agent

import "time"

type Job struct {
	Id            int    `mapstructure:"id"`   /**< primary key(sequence) */
	Name          string `mapstructure:"name"` /**< Metric */
	Status        string
	Enable        bool `mapstructure:"enable"`     /**< Enabled */
	Build         int  `mapstructure:"build"`      /**< Build version */
	StdoutLog     bool `mapstructure:"stdout_log"` /**< Standard output flag */
	Interval      int  `mapstructure:"interval"`   /**< Interval(sec) */
	Timeout       int  `mapstructure:"timeout"`    /**< Timeout(sec) */
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

// func (s *Schedule) AddJob(job *Job) {
// 	s.Jobs[job.Name] = job
// 	// append(s.Jobs, job)
// }

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
