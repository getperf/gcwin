package cmd

import (
	"github.com/getperf/gcagent/agent"
	"github.com/getperf/gcagent/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var si = config.Server{}
var isTemplate bool

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
		if err := project.Add(job, &si, isTemplate); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	agent.SetLogForeground()
	// サブコマンドのフラグ定義
	addCmd.PersistentFlags().BoolVarP(&isTemplate, "template", "t", false, "Add template flag")
	addCmd.PersistentFlags().StringVarP(&si.Server, "server", "s", "", "Target server name")
	addCmd.PersistentFlags().StringVar(&si.Url, "url", "", "URL to access the target")
	addCmd.PersistentFlags().StringVarP(&si.Ip, "ip", "i", "", "IP Address")
	addCmd.PersistentFlags().StringVarP(&si.TemplateId, "id", "", "", "Template id")
	addCmd.PersistentFlags().StringVarP(&si.User, "user", "u", "", "specific user")
	addCmd.PersistentFlags().StringVarP(&si.Password, "password", "p", "", "specific password")

	rootCmd.AddCommand(addCmd)
}
