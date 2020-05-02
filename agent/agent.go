package agent

import (
	"context"

	"github.com/getperf/gcagent/config"
	log "github.com/sirupsen/logrus"
)

func Run(ctx context.Context, config *config.Config) {
	log.Info(config)
}
