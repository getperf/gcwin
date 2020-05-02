package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/getperf/gcagent/agent"
	"github.com/getperf/gcagent/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var conf config.Config

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start agent service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		// 設定ファイルの内容を構造体にコピーする
		if err := viper.Unmarshal(&conf); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// ToDo:
		// env := ConfigEnv{host:ホスト, configPath:パス}　作成
		// config := NewConfig(ホーム、env)
		// viper.Unmarshal(&config)  で設定ファイル読み込み
		// agent.Run(context.Background(), &config) でエージェント実行
		// 旧ソース
		// fmt.Printf("config: %v\n", &config)
		agent.Run(context.Background(), &conf)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
