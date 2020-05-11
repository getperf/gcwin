package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/getperf/gcagent/agent"
	// "github.com/getperf/gcagent/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start agent service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called boot params : ", bootParameters)
		config := bootParameters.NewConfig()
		// 設定ファイルの内容を構造体にコピーする
		if err := viper.Unmarshal(config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		agentInstance, err := agent.InitAgent(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("config: %v\n", &config)
		// TODO ログ出力初期化、プロセスのデーモン化
		agentInstance.Run(context.Background(), agent.EXEC_LOOP)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
