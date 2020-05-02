package cmd

import (
	"fmt"
	"os"

	"github.com/getperf/gcagent/agent"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create collector agent project",
	Long:  `create project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		home := args[0]
		project := agent.NewProjectFromHome(home)
		if err := project.Create(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// fmt.Println(project.SampleConfig())
		// project.Add("linuxconf", "ostrich", "192.168.10.1", "admin")
		// // 設定ファイルの内容を構造体にコピーする
		// if err := viper.Unmarshal(&schedule); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// schedule.BackGround = false
		// fmt.Printf("schedule: %v\n", &schedule)
		// agent.Run(context.Background(), cfgFile, &schedule)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
