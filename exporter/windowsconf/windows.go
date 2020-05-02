package windowsconf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/getperf/gcagent/config"
	"github.com/getperf/gcagent/exporter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (e *Windows) configTest(config *config.Config) error {
	return nil
}

func (e *Windows) Setup() {
	fmt.Printf("export '%s' through Windows platform\n", e.Server)
}

func (e *Windows) writeScript(doc io.Writer, env *exporter.Env) error {
	tmpl, err := template.ParseFiles("powershell.tpl")
	if err != nil {
		return errors.Wrap(err, "failed read template")
	}
	var filteredCommands []*Command
	for _, command := range append(commands, e.Commands...) {
		if command.Level > env.Level {
			continue
		}
		log.Info("add test item ", command.Id)
		filteredCommands = append(filteredCommands, command)
	}
	if err := tmpl.Execute(doc, filteredCommands); err != nil {
		return errors.Wrap(err, "failed generate script")
	}
	return nil
}

func (e *Windows) CreateScript(env *exporter.Env) error {
	log.Info("create temporary log dir for test ", env.Datastore)
	e.ScriptPath = filepath.Join(env.Datastore, "get_windows_inventory.ps1")
	outFile, err := os.OpenFile(e.ScriptPath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed create script")
	}
	defer outFile.Close()
	e.writeScript(outFile, env)
	return nil
}

func (e *Windows) Run(env *exporter.Env) error {
	if err := e.CreateScript(env); err != nil {
		return errors.Wrap(err, "run windows inventory")
	}
	if runtime.GOOS != "windows" {
		return fmt.Errorf("windows powershell environment only")
	}
	cmdPowershell := []string{
		"powershell",
		e.ScriptPath,
		env.Datastore,
	}
	log.Info(cmdPowershell)
	cmd := exec.Command(cmdPowershell[0], cmdPowershell[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "create windows inventory command pipe")
	}
	startTime := time.Now()
	testId := ""
	intervalTime := startTime
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if testId != "" {
			log.Info(testId, ", Elapse: ", time.Since(intervalTime))
		}
		testId = scanner.Text()
		intervalTime = time.Now()
	}
	if testId != "" {
		log.Info(testId, ", Elapse: ", time.Since(intervalTime))
	}
	cmd.Wait()
	// tio := &timeout.Timeout{
	// 	Cmd:       cmd,
	// 	Duration:  defaultTimeoutDuration,
	// 	KillAfter: timeoutKillAfter,
	// }
	// exitstatus, stdout, stderr, err := tio.Run()

	log.Infof("finish windows inventory script, elapse [%s]", time.Since(startTime))
	// if err != nil {
	// 	return fmt.Errorf("test3 %s", err)
	// }
	// enc := japanese.ShiftJIS
	// bb, _, err := transform.Bytes(enc.NewDecoder(), []byte(stderr))

	// log.Info("RC:", exitstatus)
	// log.Info("STDOUT:", stdout)
	// log.Info("STDERR:", string(bb))

	return nil

}

func init() {
	exporter.Add("windowsconf", func() exporter.Exporter {
		return &Windows{
			Server: "localhost",
		}
	})
}
