package cmd

import (
	"github.com/getperf/gcagent/agent"
	"github.com/getperf/gcagent/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var si = config.Server{}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add exporter job to the project",
	Example: "gcagent add windowsconf [flags]",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		job := args[0]
		project, err := agent.NewProject(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := project.Add(job, &si); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	agent.SetLogForeground()
	// サブコマンドのフラグ定義
	addCmd.PersistentFlags().StringVarP(&si.Server, "server", "s", "", "target server")
	addCmd.PersistentFlags().StringVar(&si.Url, "url", "", "service url")
	addCmd.PersistentFlags().StringVarP(&si.Ip, "ip", "i", "", "ip address")
	addCmd.PersistentFlags().StringVarP(&si.UserId, "userid", "u", "", "user id")
	addCmd.PersistentFlags().StringVar(&si.User, "su", "", "specific user")
	addCmd.PersistentFlags().StringVar(&si.Password, "sp", "", "specific password")

	rootCmd.AddCommand(addCmd)
}
