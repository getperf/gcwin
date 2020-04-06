package gcwin

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type DateFormat int

const (
	DEFAULT DateFormat = iota
	YYYYMMDD
	HHMISS
	YYYYMMDD_HHMISS
	DIR
)

type ProcMode int

const (
	INIT ProcMode = iota
	WAIT
	RUN
	TIMEOUT
	END
	ERROR
)

// Default set fo gcwin.yaml
const (
	DEFAULT_DISK_CAPACITY   int = 10
	DEFAULT_SAVE_HOUR           = 3
	DEFAULT_RECOVERY_HOUR       = 3
	DEFAULT_MAX_ERROR_LOG       = 10
	DEFAULT_COMMAND_TIMEOUT     = 86400
	LIMIT_MAX_ERROR_LOG         = 10000
	DEFAULT_PROXY_PORT          = 8080
	DEFAULT_LOG_LEVEL           = 3
	DEFAULT_LOG_SIZE            = 100000
	DEFAULT_LOG_ROTATION        = 5
	POLLER_INTERVAL             = 60
)

type Config struct {
	Module        int       // Module id(S:Scheduler, C:Collector, W:Worker)
	ElapseTime    int       // Elapsed time(sec)
	StartTime     time.Time // Start time(UTC)
	Mode          ProcMode  // Status(INIT, WAIT, RUN, ...)
	ManagedPid    int       // Scheduler process id
	DaemonFlag    int       // Daemon flag
	Host          string    // Hostname(Convert to lowercase , except the domain part)
	ServiceName   string    // HA service name
	Pwd           string    // Current directory
	Home          string    // Home directory
	ParameterFile string    // Parameter file
	ProgramName   string    // Program name
	ProgramPath   string    // Program path

	OutDir        string // Metric collection directory
	WorkDir       string // Work directory
	WorkCommonDir string // Common work directory
	ArchiveDir    string // Archive directory
	ScriptDir     string // Script directory
	BinDir        string // Binary directory
	LogDir        string // Application log directory

	ExitFlag string // Exit flag file
	PidFile  string // PID file
	PidPath  string // PID absolute path

	// SslConfig *soap.SSLConfig // struct  GPFSSLConfig_t  *sslConfig; // SSL manager
	// logConfig     *LogConfig  // struct  GPFLogConfig_t  *logConfig; // Log manager
	Schedule *Schedule // struct  GPFSchedule_t   *schedule;  // Scheduler
}

type ConfigEnv struct {
	host        string
	programName string
	configFile  string
}

// NewCOnfigEnv は home を引数に設定情報を初期化します
func NewConfigEnv() *ConfigEnv {
	hostName, err := GetHostname()
	if err != nil {
		log.Errorf("get hostname for initialize config %s", err)
		hostName = "UnkownHost"
	}
	return NewConfigEnvBase(hostName, cmdName, "gcwin.ini")
}

func NewConfigEnvBase(host, program, config string) *ConfigEnv {
	configEnv := ConfigEnv{
		host:        host,
		programName: program,
		configFile:  config,
	}
	return &configEnv
}

func NewConfig(home string, configEnv *ConfigEnv) *Config {
	var config Config

	// ファイルパス初期化
	// sslDir := filepath.Join(home, "network")

	pidFile := fmt.Sprintf("_pid_%s", configEnv.programName)
	pid := os.Getpid()
	workDir := fmt.Sprintf("_%d", pid)

	// プロセス定義
	config.Module = 'S'           // モジュール識別子(S, C, W)
	config.ElapseTime = 0         // 経過時間(秒)
	config.StartTime = time.Now() // 起動時間(UTC)
	config.Mode = INIT            // 実行状態(init, run, end, ...)
	config.ManagedPid = 0         // 管理用プロセスID(ワーカで使用)
	config.DaemonFlag = 0         // デーモン化フラグ

	// プログラム名定義
	config.Host = configEnv.host               // ホスト名
	config.ServiceName = ""                    // HAサービス名
	config.Home = home                         // ホームディレクトリ
	config.ProgramName = configEnv.programName // プログラム名

	// ディレクトリ定義
	config.OutDir = filepath.Join(home, "datasotore")    // 採取データディレクトリ
	config.WorkDir = filepath.Join(home, "tmp", workDir) // ワークディレクトリ
	config.WorkCommonDir = filepath.Join(home, "tmp")    // 共有ワークディレクトリ
	config.ArchiveDir = filepath.Join(home, "arc")       // アーカイブ保存ディレクトリ
	config.ScriptDir = filepath.Join(home, "script")     // スクリプトディレクトリ
	config.BinDir = filepath.Join(home, "bin")           // バイナリディレクトリ
	config.LogDir = filepath.Join(home, "log")           // アプリログディレクトリ

	// 管理ファイル
	config.ExitFlag = "_exitFlag"                        // 終了フラグファイル
	config.PidFile = pidFile                             // PIDファイル
	config.PidPath = filepath.Join(home, "var", pidFile) // PIDファイル(絶対パス)

	// 構造体定義
	config.Schedule = NewSchedule() // スケジュール構造体

	return &config
}

func (config *Config) GetServiceOrHostName() string {
	host := config.Host
	if config.ServiceName != "" {
		host = config.ServiceName
	}
	return host
}

func (config *Config) CheckConfig() error {
	if config.Schedule == nil {
		return fmt.Errorf("not found schedule")
	}
	schedule := config.Schedule
	collectors := schedule.Collectors
	if len(collectors) == 0 {
		return fmt.Errorf("not found collector")
	}
	for statName, collector := range collectors {
		log.Info("check ", statName)
		if collector.StatTimeout == 0 {
			collector.StatTimeout = collector.StatInterval
		}
		if collector.StatMode == "" {
			collector.StatMode = "concurrent"
		}
	}
	return nil
}
