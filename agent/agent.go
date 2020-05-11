package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/getperf/gcagent/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// main
// 	(*c)InitAgent()
// 	(*c)PrepareTask(*Job)
// 	(*c)PostTask(*Task)
// 	(*c)SaveReport(*Task, *Datastore)
// 	(*c)Run()
// 	(*c)RunCollector(*Job, count, chan)
// 	(*c)RunWithContext(ctx)

type RunMode int

const (
	EXEC_ONCE RunMode = iota
	EXEC_LOOP
)

type Agent struct {
	Home     string
	jobCount int
	cfg      *config.Config
}

func InitAgent(c *config.Config) (*Agent, error) {
	agent := Agent{
		Home: c.Home,
		cfg:  c,
	}
	// PID ファイル更新
	err := c.WriteProcInfo()
	return &agent, err
}

func (agent *Agent) RunTask(job *config.Job, tasks chan<- *Task) {
	log.Infof("Run collector : %s", job.Name)
	task, err := NewTask(agent.cfg, job.Name)
	if err != nil {
		log.Errorf("run task :%s", err)
	}
	time.Sleep(time.Second * 3)
	task.EndTime = time.Now()
	tasks <- task
}

func (agent *Agent) PostTask(task *Task) {
	// config := agent.cfg
	log.Info("Post task : ", task.JobName)
	// ds := task.Datastore
	// id := ds.Path()
	// // elapse := task.EndTime.Sub(task.StartTime)
	// log.Info("end task [", id, "]")
	// if err := config.SaveReport(task, ds); err != nil {
	// 	log.Error("report failed ", err)
	// }
	// log.Debug("check send enable ", config.Schedule.RemhostEnable)
	// if config.Schedule.RemhostEnable == true {
	// 	if err := config.ArchiveData(task, ds); err != nil {
	// 		log.Error("zip failed ", err)
	// 		return
	// 	}
	// 	if err := config.SendCollectorDataAll(ds); err != nil {
	// 		log.Error("send failed ", err)
	// 		return
	// 	}
	// }
	// if err := config.PurgeData(task, ds); err != nil {
	// 	log.Error("purge failed ", err)
	// 	return
	// }
}

func (agent *Agent) Run(ctx context.Context, runMode RunMode) error {
	ctx, cancel := context.WithCancel(ctx)
	if err := agent.VerifyLicense(); err != nil {
		return errors.Wrap(err, "collector process")
	}
	jobs := agent.cfg.Jobs
	tasks := make(chan *Task, len(jobs)*2)
	for statName, job := range jobs {
		if !job.Enable {
			continue
		}
		log.Infof("start task : %s", statName)
		agent.jobCount++
		go func(ctx context.Context, job *config.Job) {
			count := 0
			go agent.RunTask(job, tasks)
			count++
			if runMode == EXEC_ONCE {
				return
			}
			interval := time.Second * time.Duration(job.Interval)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					log.Info("task canceled")
					return
				case <-ticker.C:
					count++
					// log.Infof("continue task %s(%d)", statName, count)
					go agent.RunTask(job, tasks)
				}
			}
		}(ctx, job)
	}
	log.Info("start process monitor")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()
	count := 0

	taskCount := 0
monitorChannel:
	for {
		select {
		case task := <-tasks:
			agent.PostTask(task)
			if runMode == EXEC_ONCE {
				taskCount++
				log.Infof("finish task %d/%d", taskCount, agent.jobCount)
				if taskCount == agent.jobCount {
					break monitorChannel
				}
			}
		case <-ticker.C:
			status, err := agent.CheckExitFile()
			if err != nil {
				log.Error("check exit file ", err)
			} else if status == "STOP" {
				agent.Stop()
				cancel()
			}
			// log.Println("CHECK stop flag file", count)
			count++
		case <-quit:
			log.Debug("CATCH SIGNAL")
			agent.Stop()
			cancel()
			return fmt.Errorf("catch signal")

			// errors.Wrap(ctx.Err(), "agent monitor canceled")
			// os.Exit(0)
		case <-ctx.Done():
			agent.Stop()
			return errors.Wrap(ctx.Err(), "agent monitor canceled")
		}
	}
	return nil
}
