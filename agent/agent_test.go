package agent

import (
	"context"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	// . "github.com/getperf/gcagent/common"
	"github.com/getperf/gcagent/config"
)

func createConfig() *config.Config {
	config := config.NewConfig("../testdata/ptune/", config.NewConfigEnv())
	toml.DecodeFile("../testdata/ptune/gcagent.toml", &config)
	config.Jobs["windowsconf"].Interval = 1
	config.CheckConfig()
	return config
}

func TestRunOnce(t *testing.T) {
	// config := createConfig()
	config := createConfig()
	agent, _ := InitAgent(config)
	// ctx, _ := context.WithCancel()
	err := agent.Run(context.Background(), EXEC_ONCE)
	t.Log(err)
	if err != nil {
		t.Error("run single exec")
	}
}

func TestRunLoop(t *testing.T) {
	config := createConfig()
	agent, _ := InitAgent(config)
	ctx, cancel := context.WithCancel(context.Background())

	go agent.Run(ctx, EXEC_LOOP)

	time.Sleep(time.Second * 3)
	cancel()
	time.Sleep(time.Second * 1)
	// ToDo : キャンセル後にタスクレポートから正常実行の検証を追加
	// if err != nil {
	// 	t.Error("run single exec")
	// }
}
