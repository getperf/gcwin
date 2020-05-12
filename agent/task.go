package agent

// タスクは指定したジョブに対して以下の処理を実行します。

// * リモートバッチ採取
//		エクスポーターのSetup()を呼び出し
// * ローカル採取,リモート採取
//		エクスポーターのRun()を呼び出し

// エクスポーターは各ジョブの個別シナリオを実行します。
// タスクはエクスポータの親処理として、各エクスポータの呼び出しをする
// 共通処理を実行します。処理フローは以下の通りです。

// 1. Setup : 初期処理
// 		指定したジョブの関連情報を env 構造体に登録します
// 			gcagent.toml の基本設定。採取レベルなど
// 			template/{template}/{job}.toml のテンプレート情報
// 			node/{server}/{job}.toml のサーバー情報
// 		エクスポーターの Setup(env) を呼び出します
// 			リモートで複数のサーバ情報のバッチ取得処理を実行します
// 				例えば、 vCenter から複数の VM 情報をリモート採取するなど
// 			Setup(env)は全サーバの採取を実行し、各サーバのデータストアに
// 			結果を保存します

// 2. Run(ローカル)
// 		gcagent.toml 基本設定の LocalExec が true の場合に実行します
// 		エクスポーター構造体に自ホストのインスタンスを生成します
// 		自サーバのノード定義ファイルがある場合は設定を読込みます
// 		インスタンス変数 IsRemote を false に設定
// 		エクスポーターの Run(env)を実行します

// 3. Run(リモート)
//      リモート採取する複数のサーバ設定ファイルを順に処理します
//      エクスポーター構造体にリモートサーバのインスタンスを生成します
//      サーバのノード定義ファイルを読込みます
//      インスタンス変数 IsRemote をfalse に設定し
//      エクスポータの Run(env) を実行します
//      エクスポータは IsRemote がtrueの場合はセッションの接続をします

// 4. CleanUp 終了処理
//      一連の処理の後始末

import (
	"fmt"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/getperf/gcagent/config"
	"github.com/getperf/gcagent/exporter"
	_ "github.com/getperf/gcagent/exporter/all"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ProcMode int

const (
	INIT ProcMode = iota
	WAIT
	RUN
	TIMEOUT
	END
	ERROR
)

type Task struct {
	Cfg        *config.Config
	JobName    string
	LocalExec  bool
	Exp        exporter.Creator
	Status     ProcMode  /**< Process status */
	StartTime  time.Time /**< Start time(UTC) */
	EndTime    time.Time /**< End time(UTC) */
	Datastores map[string]*config.Datastore
	// JobResults    map[string]config.JobResult
}

func NewTask(cfg *config.Config, jobName string) (*Task, error) {
	exp, ok := exporter.Exporters[jobName]
	if !ok {
		return nil, fmt.Errorf("exporter not defined '%s'", jobName)
	}
	now := time.Now()
	datastores := make(map[string]*config.Datastore, 100)
	// ローカルサーバのデータストア作成
	datastores[cfg.Host] = cfg.NewDatastore(jobName, cfg.Host, now)

	serverConfigs, err := cfg.ServerConfigs(jobName)
	if err != nil {
		return nil, errors.Wrap(err, "task creating")
	}
	for server := range serverConfigs {
		datastores[server] = cfg.NewDatastore(jobName, server, now)
	}

	task := Task{
		Cfg:        cfg,
		JobName:    jobName,
		Exp:        exp,
		LocalExec:  cfg.LocalExec(jobName),
		Status:     INIT,
		StartTime:  now,
		Datastores: datastores,
	}
	return &task, nil
}

func (task *Task) MakeExporterEnv() (*exporter.Env, error) {
	cfg := task.Cfg
	jobName := task.JobName
	serverConfigs, err := cfg.ServerConfigs(jobName)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}
	templateConfigs, err := cfg.TemplateConfigs(jobName)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}
	// 共通の設定の Level,DryRun 以外は既定値をセット
	env := exporter.Env{
		Level:  cfg.RunLevel,
		DryRun: cfg.DryRun,

		// バッチ用の定義ファイルパスをセット。バッチ以外は未使用
		TemplateConfigs: templateConfigs,
		ServerConfigs:   serverConfigs,
	}
	return &env, nil
}

func (task *Task) MakeServer(env *exporter.Env, host string, localExec bool) (exporter.Exporter, error) {
	server := task.Exp()
	ds, ok := task.Datastores[host]
	if !ok {
		return server, fmt.Errorf("datastore not found in creating server")
	}
	env.Datastore = ds.Path()
	env.LocalExec = localExec

	// ノード定義からサーバ情報取得
	serverConfig, ok := env.ServerConfigs[host]
	if ok {
		_, err := toml.DecodeFile(serverConfig, server)
		if err != nil {
			return server, errors.Wrap(err, "making server")
		}
	}
	// テンプレートからサーバ情報取得
	src := reflect.ValueOf(server).Elem()
	templateIdVal := src.FieldByName("TemplateId")
	templateId := ""
	if templateIdVal.CanSet() {
		templateId = templateIdVal.Interface().(string)
		log.Info("Use template : ", templateId)
	}
	if templateId != "" {
		templateConfig, ok := env.TemplateConfigs[templateId]
		if !ok {
			return server, fmt.Errorf("template not found : %s", templateId)
		}
		templateServer := task.Exp()
		_, err := toml.DecodeFile(templateConfig, templateServer)
		if err != nil {
			return server, errors.Wrap(err, "making server")
		}
		if err := mergo.Merge(server, templateServer); err != nil {
			return server, errors.Wrap(err, "making server")
		}
	}
	return server, nil
}

func (t *Task) LocalServer(e *exporter.Env) (exporter.Exporter, error) {
	host := t.Cfg.Host
	return t.MakeServer(e, host, true)
}

func (t *Task) RemoteServer(e *exporter.Env, host string) (exporter.Exporter, error) {
	return t.MakeServer(e, host, false)
}

func (task *Task) Run() error {
	exp := task.Exp
	log.Info("Desc ", exp().Label())
	env, err := task.MakeExporterEnv()
	if err != nil {
		return errors.Wrap(err, "create exporter env in task runnning")
	}
	exp().Setup(env)
	if task.LocalExec {
		server, err := task.LocalServer(env)
		if err != nil {
			return errors.Wrap(err, "prepare local server")
		}
		if err := server.Run(env); err != nil {
			return errors.Wrap(err, "run local")
		}
	}
	for serverName, serverConfig := range env.ServerConfigs {
		log.Info("server config : ", serverConfig)
		server, err := task.RemoteServer(env, serverName)
		if err != nil {
			return errors.Wrap(err, "prepare remove server")
		}
		if err := server.Run(env); err != nil {
			log.Warnf("run remote failed %s", err)
		}
	}
	task.EndTime = time.Now()
	return nil
}
