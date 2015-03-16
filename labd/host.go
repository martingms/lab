package main

import (
	"errors"
	"io"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

type Host struct {
	Name     string
	Host     string
	Username string
	KeyPath  string
	sshConf  *ssh.ClientConfig
}

func (h *Host) Run(cmd string, args ...string) (*Command, error) {
	if h.Name == "localhost" {
		return h.runLocal(cmd, args...)
	}

	if h.sshConf == nil {
		err := h.generateSSHConfig()
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("Remote commands not yet implemented")

	// TODO: Check if we have connection, if not dial, etc.
}

func (h *Host) runLocal(cmd string, args ...string) (*Command, error) {
	command := exec.Command(cmd, args...)
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = command.Start()
	if err != nil {
		return nil, err
	}

	return &Command{
		StdinPipe:  stdin,
		StdoutPipe: stdout,
		StderrPipe: stderr,
		localCmd:   command,
	}, nil
}

func (h *Host) generateSSHConfig() error {
	h.sshConf = &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			// TODO: Take public keys here.
			ssh.Password("password"),
		},
	}

	return nil
}

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
