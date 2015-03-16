package main

import (
    "os/exec"
    "io"
    "errors"
)

type Command struct {
	StdinPipe  io.WriteCloser
	StdoutPipe io.ReadCloser
	StderrPipe io.ReadCloser
	localCmd   *exec.Cmd
	//remoteCmd
}

func (cmd *Command) Wait() error {
	if cmd.localCmd != nil {
		return cmd.localCmd.Wait()
	}

	return errors.New("Remote commands not yet implemented")
}
