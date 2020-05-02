package agent

import (
	"testing"

	"github.com/spf13/viper"
)

func TestReadSchedule(t *testing.T) {

	viper.SetConfigType("toml")
	viper.SetConfigName("gcagent.toml")
	viper.AddConfigPath("../testdata/ptune/")
	var schedule Schedule

	if err := viper.ReadInConfig(); err != nil {
		t.Error("Error reading config file, ", err)
	}
	err := viper.Unmarshal(&schedule)
	t.Log(schedule)
	if err != nil {
		t.Error("unable to decode into struct, ", err)
	}
	if schedule.DiskCapacity != 0 {
		t.Error("config disk capacity")
	}
	if schedule.Jobs["linuxconf"].Enable != true {
		t.Error("config linux collector setting")
	}
	// t.Log("port for this application is %d", configuration.Server.Port)

	// // viper.SetConfigType("toml")
	// // viper.SetConfigFile("testdata/ptune/gcagent.toml")
	// viper.SetConfigFile("./testdata/ptune/gcagent.toml")
	// // viper.AddConfigPath("testdata/ptune/")

	// viper.AutomaticEnv() // read in environment variables that match
	// if err := viper.ReadInConfig(); err == nil {
	// 	t.Error("Using config file:", viper.ConfigFileUsed())
	// }
	// t.Log("disk_capacity:", viper.Get("disk_capacity"))
	// t.Log(viper.ConfigFileUsed())
	// var schedule = Schedule{}
	// if err := viper.Unmarshal(&schedule); err != nil {
	// 	t.Error("viper unmarshal error")
	// }
	// t.Log("config: ", schedule)
}
