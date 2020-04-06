package main

import (
	"log"

	"github.com/getperf/gcwin/cmd"
	_ "github.com/kardianos/minwinsvc"
)

func main() {
	log.SetFlags(0)
	cmd.Execute()
}
