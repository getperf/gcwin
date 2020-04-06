package gcwin

import (
	"context"
	"io"

	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
)

const cmdName = "gcwin"

func gops() error {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	return nil
}

func Run(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
	gops()

	return nil
}
