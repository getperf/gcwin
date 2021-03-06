package agent

import (
    "time"
)

type Schedule struct {
    DiskCapacity int `mapstructure:"disk_capacity"` // Disk free threshold(%)
    SaveHour     int `mapstructure:"save_hour"`     // Metric data retention(H)
    RecoveryHour int `mapstructure:"recovery_hour"` // Metric data retransmission(H)
    MaxErrorLog  int `mapstructure:"max_error_log"` // Max rows of error output

    Pid        int  /**< Scheduler process id */
    Status     int  /**< Process status */
    BackGround bool /**< Back Ground flag */

    LogLevel     int  `mapstructure:"log_level"`     // Log level
    DebugConsole bool `mapstructure:"debug_console"` // Console log enabled
    LogSize      int  `mapstructure:"log_size"`      // Log size
    LogRotation  int  `mapstructure:"log_rotation"`  // Number of log rotation
    LogLocalize  bool `mapstructure:"log_localize"`  // Flag of Japanese log

    HanodeEnable bool   `mapstructure:"hanode_enable"` // HA node check flag
    HanodeCmd    string `mapstructure:"hanode_cmd"`    // HA node check script

    PostEnable bool   `mapstructure:"post_enable"` // Post command enabled
    PostCmd    string `mapstructure:"post_cmd"`    // Post command

    RemhostEnable bool   `mapstructure:"remhost_enable"` // Remote transfer enabled
    UrlCM         string `mapstructure:"url_cm"`         // Web service URL (Configuration manager)
    UrlPM         string `mapstructure:"url_pm"`         // Web service URL (Performance manager)

    SiteKey     string `mapstructure:"site_key"` // Site key
    SoapTimeout int    `mapstructure:"soap_timeout"`

    ProxyEnable bool   `mapstructure:"proxy_enable"` // HTTP proxy enabled
    ProxyHost   string `mapstructure:"proxy_host"`   // Proxy host
    ProxyPort   int    `mapstructure:"proxy_port"`   // Proxy port

    ServiceUrl string `mapstructure:"service_url"`

    LastUpdate  time.Time /**< Last update of parameter file */
    ParseFailed bool      /**< Set true if parser failed */

    License *License
    Jobs    map[string]*Job // Job pids
}

func NewSchedule() *Schedule {
    var schedule Schedule
    schedule.DiskCapacity = DEFAULT_DISK_CAPACITY
    schedule.SaveHour = DEFAULT_SAVE_HOUR
    schedule.RecoveryHour = DEFAULT_RECOVERY_HOUR
    schedule.MaxErrorLog = DEFAULT_MAX_ERROR_LOG
    schedule.LogLevel = DEFAULT_LOG_LEVEL
    schedule.LogSize = DEFAULT_LOG_SIZE
    schedule.LogRotation = DEFAULT_LOG_ROTATION

    schedule.LogLocalize = true
    schedule.LastUpdate = time.Now()
    schedule.ParseFailed = false
    schedule.RemhostEnable = false
    schedule.License = NewLicense()
    schedule.Jobs = make(map[string]*Job)
    return &schedule
}

// func (schedule *Schedule) AddJob(job *Job) {
//     schedule.Jobs[job.Name] = job
// }

func (config *Config) GetJob(statName string) *Job {
    if schedule := config.Schedule; schedule != nil {
        return schedule.GetJob(statName)
    }
    return nil
}

func (schedule *Schedule) GetJob(statName string) *Job {
    if val, ok := schedule.Jobs[statName]; ok {
        return val
    }
    return nil
}

// func (schedule *Schedule) FindOrCreateJob(statName string) *Job {
//     if collector, ok := schedule.Jobs[statName]; ok {
//         return collector
//     } else {
//         id := len(schedule.Jobs) + 1
//         newCollector := NewJob(statName)
//         newJob.Id = id
//         schedule.AddJob(newJob)
//         return newJob
//     }
// }
