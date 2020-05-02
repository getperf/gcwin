package linuxconf

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Songmu/timeout"
	log "github.com/sirupsen/logrus"
)

type CommandInfo struct {
	CmdLine string
	OutPath string
	Timeout int

	Executed bool
	Pid      int
	ExitCode int
	Status   string
	OutBuf   string
}

// defaultTimeoutDuration is the duration after which a command execution will be timeout.
// timeoutKillAfter is option of `RunCommand()` set waiting limit to `kill -kill` after
// terminating the command.
var (
	defaultTimeoutDuration = 30 * time.Second
	timeoutKillAfter       = 1 * time.Second
)

var errTimedOut = errors.New("command timed out")

var cmdBase = []string{"sh", "-c"}

func init() {
	if runtime.GOOS == "windows" {
		cmdBase = []string{"cmd", "/U", "/c"}
	}
}

func MakeContext(timeout int) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		duration := time.Duration(timeout) * time.Second
		return context.WithTimeout(context.Background(), duration)
	} else {
		return context.WithCancel(context.Background())
	}
}

func (c *CommandInfo) ExecCommandRedirect() error {
	ctx, cancel := MakeContext(c.Timeout)
	defer cancel()
	return c.ExecCommandRedirectWithContext(ctx)
}

func (c *CommandInfo) ExecCommandRedirectWithContext(ctx context.Context) error {
	cmdArgs := append(cmdBase, filepath.FromSlash(c.CmdLine))

	log.Debug("exec command direct ", cmdArgs)
	args := append([]string{}, cmdArgs...)
	cmd := exec.Command(args[0], args[1:]...)

	if c.OutPath == "" {
		return fmt.Errorf("command output is nil")
	}
	// outDir := GetParentPath(c.OutPath, 1)
	// if err := os.MkdirAll(outDir, 0777); err != nil {
	// 	return fmt.Errorf("command output directory create failed")
	// }
	outFile, err := os.OpenFile(c.OutPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		return fmt.Errorf("command output %s %s", c.OutPath, err)
	}
	defer outFile.Close()
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  defaultTimeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	if c.Timeout != 0 {
		tio.Duration = time.Duration(c.Timeout) * time.Second
	}
	exitStatus, err := tio.RunContext(ctx)
	exitCode := -1
	if err == nil && (exitStatus.IsTimedOut() || exitStatus.Signaled) {
		log.Info("timeout : ", cmd.ProcessState.Pid())
		err = errTimedOut
		exitCode = exitStatus.GetChildExitCode()
	}
	exitCode = exitStatus.GetChildExitCode()
	if exitCode != 0 {
		log.Error("exit [", cmd.ProcessState.Pid(), "] rc=", exitCode)
	}
	c.Executed = true
	c.Pid = cmd.ProcessState.Pid()
	c.ExitCode = exitCode
	c.Status = cmd.ProcessState.String()
	return nil
}
