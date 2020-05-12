package agent

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	. "github.com/getperf/gcagent/common"
	"github.com/getperf/gcagent/config"
	"github.com/getperf/gcagent/exporter"
	_ "github.com/getperf/gcagent/exporter/all"
	log "github.com/sirupsen/logrus"
)

type Hoge struct {
	Server     string `toml:"server"`
	IP         string `toml:"ip"`
	TemplateId string `toml:"template_id"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
}

func (e *Hoge) Label() string { return "Hoge" }

func (e *Hoge) Config(configType exporter.ConfigType) string {
	return fmt.Sprintf("; '%s' Hoge config #%v", e.Server, configType)
}

func (e *Hoge) Setup(env *exporter.Env) error {
	fmt.Printf("start '%s' through Hoge\n", e.Server)
	return nil
}

func (e *Hoge) Run(env *exporter.Env) error {
	log.Info("run ", e.Server)
	return fmt.Errorf("connect '%s' by '%s', run '%s' through Hoge",
		e.IP, e.User, e.Server)
}

func init() {
	server := "localhost"
	if localHost, err := GetHostname(); err == nil {
		server = localHost
	}
	exporter.AddExporter("hogeconf", func() exporter.Exporter {
		return &Hoge{
			Server: server,
		}
	})
}

func createHogeConfig() *config.Config {
	cfg := config.NewConfig("../testdata/task_test_home/", config.NewConfigEnv())
	cfg.Host = "ostrich"
	toml.DecodeFile("../testdata/task_test_home/gcagent.toml", &cfg)
	cfg.CheckConfig()
	return cfg
}

func createTask(jobName string) *Task {
	cfg := createHogeConfig()
	task, _ := NewTask(cfg, jobName)
	return task
}

// Create task.

func TestNewTask(t *testing.T) {
	task := createTask("hogeconf")
	if task.LocalExec != true {
		t.Error("new task check local exec")
	}
	server := task.Exp()
	if server.Label() != "Hoge" {
		t.Error("new task call exporter label")
	}
	ds := task.Datastores["ostrich"]
	t.Log(ds.Path())
	if ds.Path() == "" || ds.RelDir() == "" || ds.TargetDir() == "" {
		t.Error("new task datastore")
	}
}

// Read nonexistent job.

func TestNGNewTask(t *testing.T) {
	cfg := createHogeConfig()
	task, err := NewTask(cfg, "hogehogeconf")
	if task != nil || err == nil {
		t.Error("new task unkown job")
	}
}

// Create test env.

func TestMakeExporterEnv(t *testing.T) {
	task := createTask("hogeconf")
	env, err := task.MakeExporterEnv()
	if env.Level != 0 || env.DryRun != false ||
		env.Datastore != "" || env.LocalExec != false ||
		env.Messages != "" || env.ErrMsgs != "" {
		t.Error("make exporter env")
	}
	if err != nil {
		t.Error("make exporter env")
	}
}

// Reading the server config file.

func TestDecodeServer(t *testing.T) {
	task := createTask("hogeconf")
	server := task.Exp()
	toml.Decode(`server = "ostrich2"`, server)
	config := server.Config(exporter.SERVER)
	if config != "; 'ostrich2' Hoge config #2" {
		t.Errorf("export : %s", config)
	}
	server2 := task.Exp()
	toml.Decode(`server = "hogehoge2"`, server2)
	if server2.Config(exporter.SERVER) != "; 'hogehoge2' Hoge config #2" {
		t.Errorf("export : %s", config)
	}
}

// Create local server structure with hostname 'ostrich' without the config file.

func TestNoConfigLocalServer(t *testing.T) {
	cfg := createHogeConfig()
	cfg.Host = "hogehoge"
	task, _ := NewTask(cfg, "hogeconf")
	env, _ := task.MakeExporterEnv()
	server, err := task.LocalServer(env)
	t.Log(env)
	t.Log(server)
	if env.Datastore == "" || env.LocalExec != true || err != nil {
		t.Errorf("local server config : %s", err)
	}
	if server.Label() != "Hoge" {
		t.Error("create local server")
	}
}

// Create local server structure with hostname 'ostrich'.
// Only set the environment LocalExec is true, and reading the config file as is.

func TestLocalServer(t *testing.T) {
	task := createTask("hogeconf")
	env, _ := task.MakeExporterEnv()
	server, err := task.LocalServer(env)
	t.Log(server)
	if env.Datastore == "" || env.LocalExec != true || err != nil {
		t.Error("local server config")
	}
	if server.Label() != "Hoge" {
		t.Error("create local server")
	}
}

// Create the server structure by reading the server config.
//
// testdata/task_test_home/node/ostrich/hogeconf.toml
//   # Required Linux Endpoint
//   #
//   server = "ostrich"
//   is_remote = true
//   ip = "192.168.10.1"
//   template_id = "admin01"
//   user = "manager"
//   password = "P@ssword"

func TestRemoteServer(t *testing.T) {
	task := createTask("hogeconf")
	env, _ := task.MakeExporterEnv()
	server, err := task.RemoteServer(env, "ostrich")
	if env.Datastore == "" || env.LocalExec != false || err != nil {
		t.Error("local server config")
	}
	err = server.Run(env)
	msg := "connect '192.168.10.1' by 'manager', run 'ostrich' through Hoge"
	if err.Error() != msg {
		t.Errorf("message result : %s, expected : %s", err, msg)
	}
}

// Create the server structure by reading the following two config files.
//
// testdata/task_test_home/node/cent7/hogeconf.toml
//
//   # Required Linux Endpoint
//   #
//   server = "cent7"
//   is_remote = true
//   ip = "192.168.0.5"
//   template_id = "admin01"
//
// testdata/task_test_home/template/admin01/linuxconf.toml
//
//   user = "psadmin"
//   password = "psadmin"

func TestUsingTemplateRemoteServer(t *testing.T) {

	task := createTask("hogeconf")
	env, _ := task.MakeExporterEnv()
	server, err := task.RemoteServer(env, "cent7")
	if env.Datastore == "" || env.LocalExec != false || err != nil {
		t.Error("local server config")
	}
	err = server.Run(env)
	// msg := "connect '192.168.0.5' by '', run 'cent7' through Hoge"
	msg := "connect '192.168.0.5' by 'psadmin', run 'cent7' through Hoge"
	if err.Error() != msg {
		t.Errorf("message result : %s, expected : %s", err, msg)
	}
}
