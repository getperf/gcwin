package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/getperf/gcagent/common"
	log "github.com/sirupsen/logrus"
)

var programName = "gcagent"

type Config struct {
	DaemonFlag    bool   // Daemon flag
	Pid           int    // Process ID
	Host          string // Hostname(Convert to lowercase , except the domain part)
	ServiceName   string // HA service name
	Pwd           string // Current directory
	Home          string // Home directory
	ConfigPath    string // Config file
	ProgramPath   string // Program path
	ExitFlagPath  string // Exit flag file
	PidPath       string // PID absolute path
	BaseDir       string // Executable file base directory
	BinDir        string // Project specific Executable file directory
	OutDir        string // Metric collection directory
	NodeDir       string // Server info directory
	TemplateDir   string // Server template info directory
	WorkDir       string // Work directory
	WorkCommonDir string // Common work directory
	ArchiveDir    string // Archive directory
	LogDir        string // Application log directory
	TslDir        string // TLS manage directory

	SaveHour     int    `toml:"save_hour"`     // Metric data retention(H)
	RecoveryHour int    `toml:"recovery_hour"` // Log retransmission time [H]
	RunLevel     int    `toml:"run_level"`     // Collection level. 0 means most
	LogLevel     string `toml:"log_level"`     // Log level
	LogRotation  int    `toml:"log_rotation"`  // Number of log rotation
	ProxyEnable  bool   `toml:"proxy_enable"`  // HTTP proxy enabled
	ProxyUrl     string `toml:"proxy_url"`     // Proxy host
	HttpTimeout  int    `toml:"http_timeout"`  // Service request timeout
	HttpRetry    int    `toml:"http_retry"`    // Service request retry count
	ServiceUrl   string `toml:"service_url"`   // Web service URL (Configuration manager)
	SiteKey      string `toml:"site_key"`      // Site key

	DryRun bool
	Jobs   map[string]*Job `toml:"jobs"` // Job pids
}

type ConfigEnv struct {
	host       string
	configPath string
}

// NewCOnfigEnv は home を引数に設定情報を初期化します
func NewConfigEnv() *ConfigEnv {
	hostName, err := GetHostname()
	if err != nil {
		log.Errorf("get hostname for initialize config %s", err)
		hostName = "UnkownHost"
	}
	return NewConfigEnvBase(hostName, "")
}

func NewConfigEnvBase(host, config string) *ConfigEnv {
	configEnv := ConfigEnv{
		host:       host,
		configPath: config,
	}
	return &configEnv
}

func NewConfig(home string, configEnv *ConfigEnv) *Config {
	var config Config

	baseDir := home
	// /tmp/go-build...ではないコンパイル済みバイナリからの実行かチェック
	exe, err := os.Executable()
	if err == nil && strings.Index(exe, "go-build") == -1 {
		baseDir = filepath.Dir(exe)
	} else {
		log.Warn("failed to get program name")
		exe = ""
	}
	log.Info("set base directory ", baseDir)

	// ファイルパス初期化
	varDir := filepath.Join(home, "var")
	tmpDir := filepath.Join(varDir, "tmp")
	pidFile := fmt.Sprintf("_pid_%s", programName)
	pid := os.Getpid()
	workDir := fmt.Sprintf("_%d", pid)
	configPath := configEnv.configPath
	if configPath == "" {
		configPath = filepath.Join(home, "gcagent.toml")
	}

	// プログラム名定義
	config.Host = configEnv.host   // ホスト名
	config.ServiceName = ""        // HAサービス名
	config.Home = home             // ホームディレクトリ
	config.Pid = pid               // プロセスID
	config.ConfigPath = configPath // パラメータファイル
	config.ProgramPath = exe       // プログラム名
	config.BaseDir = baseDir

	// ディレクトリ定義
	config.OutDir = filepath.Join(varDir, "data")           // 採取データディレクトリ
	config.NodeDir = filepath.Join(home, "node/")           // 採取対象定義ディレクトリ
	config.TemplateDir = filepath.Join(baseDir, "template") // テンプレート定義ディレクトリ
	config.BinDir = filepath.Join(home, "bin")              // プロジェクトバイナリ
	config.WorkDir = filepath.Join(tmpDir, workDir)         // ワークディレクトリ
	config.WorkCommonDir = filepath.Join(tmpDir)            // 共有ワークディレクトリ
	config.ArchiveDir = filepath.Join(varDir, "arc")        // アーカイブ保存ディレクトリ
	config.LogDir = filepath.Join(varDir, "log")            // アプリログディレクトリ

	// SSL証明書定義
	config.TslDir = filepath.Join(varDir, "network")

	// 設定ファイルパラメータ既定値定義
	config.SaveHour = 72       // Metric data retention[h]
	config.RecoveryHour = 3    // Log retransmission time [h]
	config.LogLevel = "warn"   // Log level
	config.LogRotation = 5     // Number of log rotation
	config.ProxyEnable = false // HTTP proxy enabled
	config.HttpTimeout = 300   // Service request timeout [s]
	config.HttpRetry = 3       // Service request retry count

	// 管理ファイル
	config.ExitFlagPath = filepath.Join(tmpDir, "_exitFlag") // 終了フラグファイル
	config.PidPath = filepath.Join(tmpDir, pidFile)          // PIDファイル(絶対パス)

	return &config
}

func (config *Config) GetServiceName() string {
	host := config.Host
	if config.ServiceName != "" {
		host = config.ServiceName
	}
	return host
}

func (config *Config) GetBaseDirs() []*string {
	return []*string{
		&config.OutDir,
		&config.NodeDir,
		&config.TemplateDir,
		&config.WorkDir,
		&config.WorkCommonDir,
		&config.ArchiveDir,
		&config.LogDir,
	}
}

func (config *Config) CheckConfig() error {
	jobs := config.Jobs
	if len(jobs) == 0 {
		return fmt.Errorf("not found job")
	}
	for statName, job := range jobs {
		log.Info("check ", statName)
		job.Name = statName
		if job.Timeout == 0 {
			job.Timeout = job.Interval
		}
		if job.Mode == "" {
			job.Mode = "concurrent"
		}
	}
	return nil
}
