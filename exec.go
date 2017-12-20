package inttest_utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

type Application struct {
	cmd     *exec.Cmd
	name    string
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	verbose bool

	ExitCode int
}

func StartApplication(cmdstr string, env map[string]string, verbose bool) (*Application, error) {
	cmdParts := strings.Split(cmdstr, " ")
	cmd := exec.CommandContext(context.Background(), cmdParts[0], cmdParts[1:]...)
	if env != nil {
		var envStr []string
		for k, v := range env {
			envStr = append(envStr, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = envStr
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return &Application{cmd, cmdstr, stdout, stderr, verbose, -1}, nil
}

func (a *Application) Stop() error {
	err := a.cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		return err
	}
	ps, err := a.cmd.Process.Wait()
	if err != nil {
		return err
	}
	ws := ps.Sys().(syscall.WaitStatus)
	a.ExitCode = ws.ExitStatus()

	defer a.stdout.Close()
	defer a.stderr.Close()

	if !a.verbose {
		return nil
	}

	out, _ := ioutil.ReadAll(a.stdout)
	fmt.Printf("======= %s app stdout ======\n%s\n", a.name, out)
	out, _ = ioutil.ReadAll(a.stderr)
	fmt.Printf("======= %s app stderr ======\n%s\n", a.name, out)
	return nil
}
