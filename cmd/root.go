package cmd

import (
	"fmt"
	"os"

	// "path/filepath"

	// . "github.com/getperf/gcagent/common"
	// homedir "github.com/mitchellh/go-homedir"
	// "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// -c gcagent.toml 指定がある場合

// 	絶対パスに変換 cofigPath セット。パスのファイル有無チェック。
// 	ない場合はエラーで終了。
// 	ある場合は、ベースディレクトリを取得して home にセット。

// -c gcagent.toml 指定がない場合

// 	カレントディレクトリ取得。
// 	カレントディレクトリの gcagent.toml ファイル有無チェック。
// 	ある場合は、 home 、 configPath にセット。
// 	ない場合は、home、configPath を ""(未指定)にセット。

var cfgFile string
var bootParameters BootParameters

var rootCmd = &cobra.Command{
	Use:   "gcagent",
	Short: "Inventory collector agent.",
	Long:  `Getconfig inventory collector agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if err := bootParameters.make(cfgFile); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// bootSettings, _ := makeBootSettings(cfgFile)
		fmt.Println("Root command, config : ", bootParameters)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if err := bootParameters.make(cfgFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// currentDir, _ := os.Getwd()
	// fmt.Printf("Current Dir:%s\n", currentDir)
	viper.SetConfigFile(bootParameters.configPath)
	fmt.Println("Init config : ", bootParameters)
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Printf("GETHOMEDIR:%s\n", home)

	// 	// Search config in home directory with name ".gcagent" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".gcagent")
	// }

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file ", viper.ConfigFileUsed())
	}

	// bootEnv, err := makeBootSettings(cfgFile)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Root command, config2 : ", bootSettings)
}
