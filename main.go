package main

import (
	"log"

	"github.com/getperf/gcagent/cmd"

	_ "github.com/kardianos/minwinsvc"
)

func main() {
	// agent.Run2()
	log.SetFlags(0)
	cmd.Execute()
	// err := agent.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	// if err != nil && err != flag.ErrHelp {
	// 	log.Println(err)
	// 	exitCode := 1
	// 	if ecoder, ok := err.(interface{ ExitCode() int }); ok {
	// 		exitCode = ecoder.ExitCode()
	// 	}
	// 	os.Exit(exitCode)
	// }
}
