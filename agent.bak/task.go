package agent

// Task は Collecter クラスで渡されたコマンド実行セットを TaskJob
// クラスに渡して複数のコマンド実行をスケジュールします。実行モードが
// concurrent の場合は並列に、serial の場合はシーケンシャルにTaskJob
//  を実行します。
// 全ての TaskJob 終了後、各 TaskJob 構造体から実行結果を収集してレ
// ポートを作成します。

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getperf/gcagent/exporter"
	_ "github.com/getperf/gcagent/exporter/all"
	log "github.com/sirupsen/logrus"
)

type Task struct {
	Job       *Job
	Datastore *Datastore
	Pid       int      /**< Job process id */
	Status    ProcMode /**< Process status */
	Timeout   int

	StartTime time.Time /**< Start time(UTC) */
	EndTime   time.Time /**< End time(UTC) */

	// DateDir   string /**< Date directory */
	// TimeDir   string /**< Time directory */
	Odir      string /**< Output directory */
	ScriptDir string /**< Script directory */
	NodeDir   string
	// TaskJobs  []*TaskJob
}

func NewTask(job *Job, odir string) *Task {
	task := &Task{}
	task.Job = job
	task.StartTime = time.Now()
	task.Timeout = job.Timeout
	task.Odir = odir
	// task.ScriptDir = scriptDir
	// for i, job := range job.Jobs {
	// 	taskJob := NewTaskJob(i, job, odir, scriptDir)
	// 	taskJob.Timeout = job.Timeout
	// 	task.TaskJobs = append(task.TaskJobs, taskJob)
	// }
	return task
}

func (task *Task) Run() error {
	ctx, cancel := MakeContext(task.Timeout)
	defer cancel()
	return task.RunWithContext(ctx)
}

func (task *Task) RunWithContext(ctx context.Context) error {
	var err error
	job := task.Job
	prev, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("run task %s", err)
	}
	defer os.Chdir(prev)
	os.Chdir(task.ScriptDir)
	log.Info("chdir ", task.ScriptDir)
	if _, err := os.Stat(task.Odir); !os.IsNotExist(err) {
		if err := os.RemoveAll(task.Odir); err != nil {
			return fmt.Errorf("run task %s", err)
		}
	}
	if err := os.MkdirAll(task.Odir, 0777); err != nil {
		return fmt.Errorf("run task %s", err)
	}
	log.Debug("RunWithContext ", job.Mode)
	log.Debug(job.Name)
	exp, ok := exporter.Exporters[job.Name]
	if !ok {
		log.Error("exporter not defined")
	} else {
		log.Info("Desc ", exp().Description())
		env := &exporter.Env{
			Level:     0,
			DryRun:    false,
			Datastore: task.Odir,
			NodeDir:   task.NodeDir,
		}
		if err := exp().Run(env); err != nil {
			log.Error(err)
		}
	}

	// switch job.Mode {
	// case "serial":
	// 	for _, taskJob := range task.TaskJobs {
	// 		_, err = taskJob.RunWithContext(ctx)
	// 	}
	// case "concurrent":
	// 	begin := make(chan interface{})
	// 	var wg sync.WaitGroup
	// 	for i, taskJob := range task.TaskJobs {
	// 		wg.Add(1)
	// 		// 何れも id2の taskJobが実行されてしまう問題あり
	// 		// serial では正常動作、 go func の使い方調査 ※ v1.4で修正
	// 		// go test -v -run  TestConcurrentTask
	// 		go func(i int, job *TaskJob) {
	// 			defer wg.Done()
	// 			<-begin
	// 			_, err = job.RunWithContext(ctx)
	// 		}(i, taskJob)
	// 	}
	// 	close(begin)
	// 	wg.Wait()
	// }
	task.EndTime = time.Now()
	return err
}
