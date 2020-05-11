package agent

import (
	log "github.com/sirupsen/logrus"
)

// daemon
func (agent *Agent) CheckExitFile() (string, error) {
	status := ""
	// lines, _ := config.ReadWorkFileHead(config.ExitFlag, 1)
	// if len(lines) > 0 {
	// 	status = lines[0]
	// }
	return status, nil
}

func (agent *Agent) Stop() error {
	log.Info("stop process")
	// config.RemoveWorkFile(config.ExitFlag)
	// config.RemoveWorkFile("_running_flg")
	// os.RemoveAll(config.WorkDir)
	// os.Exit(0)
	return nil

}
