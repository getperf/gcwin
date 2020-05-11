package agent

import (
	// "context"
	// "fmt"
	// "os"
	// "os/signal"
	// "time"

	// "github.com/getperf/gcagent/config"
	// "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// license
// 	(*c)AuthLicense(expire)
// 	(*c)UnzipSSLConf()
// 	(*c)DownloadLicense()
// 	(*c)CheckLicense(expired)

func (agent *Agent) LoadLicense() error {
	log.Info("load license")
	return nil
}

func (agent *Agent) CheckLicense(expired int) error {
	log.Info("check license ", expired)
	return nil
}

func (agent *Agent) VerifyLicense() error {
	cfg := agent.cfg
	log.Infof("check license, post url : %v", cfg.ServiceUrl)
	return nil
}
