package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/getperf/gcagent/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run agent command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called boot params : ", bootParameters)
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
		agentInstance.Run(context.Background(), agent.EXEC_ONCE)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
