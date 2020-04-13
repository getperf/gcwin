package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/getperf/gcagent/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var schedule agent.Schedule

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start agent service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		// 設定ファイルの内容を構造体にコピーする
		if err := viper.Unmarshal(&schedule); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		schedule.BackGround = false
		fmt.Printf("schedule: %v\n", &schedule)
		agent.Run(context.Background(), cfgFile, &schedule)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
