package agent

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
)

const cmdName = "gcagent"

func gops() error {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	return nil
}

// Run the gcagent
func Run(ctx context.Context, configPath string, schedule *Schedule) {
	gops()
	hostName, err := GetHostname()
	if err != nil {
		log.Fatal("get hostname ", err)
	}
	if configPath == "" {
		home, err := GetParentAbsPath(os.Args[0], 2)
		if err != nil {
			log.Fatal("get gcagent path ", err)
		}
		configPath = filepath.Join(home, "gcagent.toml")
	}
	configEnv := NewConfigEnvBase(hostName, cmdName, configPath)
	home, err := GetParentAbsPath(configPath, 1)
	if err != nil {
		log.Fatal("get home ", err)
	}
	config := NewConfig(home, configEnv)
	if schedule.BackGround {
		if err := SetLog(home); err != nil {
			log.Fatal("set log ", err)
		}
	}
	log.Info("agent start ", VersionMessage())
	// config.ParseConfigFile(config.ParameterFile)
	config.Schedule = schedule
	if err := config.CheckConfig(); err != nil {
		log.Fatal("check parameter ", err)
	}
	if err := SetLogLevel(config.Schedule.LogLevel); err != nil {
		log.Fatal("set log level ", err)
	}
	err = config.RunWithContext(ctx)
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		exitCode := 1
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode()
		}
		os.Exit(exitCode)
	}
}

func VersionMessage() string {
	return fmt.Sprintf("%s v%s (rev:%s)\n", cmdName, version, revision)
}
