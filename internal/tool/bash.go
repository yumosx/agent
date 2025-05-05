package tool

import (
	"bytes"
	"errors"
	"os/exec"
	"time"
)

type BashTool struct {
	process  *exec.Cmd
	sentinel string
	timeout  time.Duration
}

func NewBashSession(timeout time.Duration) *BashTool {
	return &BashTool{timeout: timeout, sentinel: "<<exit>>"}
}

func (bash *BashTool) Start() error {
	bash.process = exec.Command("/bin/bash")
	bash.process.Stdin = &bytes.Buffer{}
	bash.process.Stdout = &bytes.Buffer{}
	bash.process.Stderr = &bytes.Buffer{}

	err := bash.process.Start()
	if err != nil {
		return err
	}
	return nil
}

func (bash *BashTool) Stop() error {
	if bash.process.ProcessState != nil && bash.process.ProcessState.Exited() {
		return nil
	}

	return bash.process.Process.Kill()
}

func (bash *BashTool) Run(cmd string) (string, string, error) {
	_, err := bash.process.Stdin.(*bytes.Buffer).WriteString(cmd + "; echo" + bash.sentinel + "'\n'")
	if err != nil {
		return "", "", err
	}
	done := make(chan error, 1)
	go func() {
		_, err = bash.process.Process.Wait()
		done <- err
	}()

	select {
	case <-time.After(bash.timeout):
		return "", "", errors.New("bash commend timeout")
	case err = <-done:
		if err != nil {
			return "", "", err
		}
	}

	output := bash.process.Stdout.(*bytes.Buffer).String()
	errOut := bash.process.Stderr.(*bytes.Buffer).String()

	if idx := bytes.Index([]byte(output), []byte(bash.sentinel)); idx != -1 {
		output = output[:idx]
	}

	return output, errOut, nil
}
