package main

import (
	"log"

	"github.com/getperf/gcagent/cmd"
	_ "github.com/kardianos/minwinsvc"
)

func main() {
	log.SetFlags(0)
	cmd.Execute()
}
