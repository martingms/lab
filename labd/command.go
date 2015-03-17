// This is an attempt to abstract away the differences between
// crypto/ssh#Session and os/exec#Command

// TODO: When crypto/ssh#Session and os/exec#Command has same method signature,
// refactor this away, using an interface. (Or when io.ReadCloser passes as a
// io.Reader).
// PS: Might very well be a better way to do this, 'cause this is ugly...

package main

import (
	"io"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Command struct {
	StdinPipe  io.WriteCloser
	StdoutPipe io.Reader
	StderrPipe io.Reader
	localCmd   *exec.Cmd
	sshSession *ssh.Session
}

func StartLocalCommand(cmdStr string, args ...string) (*Command, error) {
	localCmd := exec.Command(cmdStr, args...)
	stdin, err := localCmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := localCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := localCmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	err = localCmd.Start()
	if err != nil {
		return nil, err
	}

	return &Command{
		StdinPipe:  stdin,
		StdoutPipe: stdout,
		StderrPipe: stderr,
		localCmd:   localCmd,
	}, nil
}

func StartSSHCommand(session *ssh.Session, cmdStr string, args ...string) (*Command, error) {
	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = session.Start(cmdStr + " " + strings.Join(args, " "))
	if err != nil {
		return nil, err
	}

	return &Command{
		StdinPipe:  stdin,
		StdoutPipe: stdout,
		StderrPipe: stderr,
		sshSession: session,
	}, nil
}

func (cmd *Command) Wait() error {
	if cmd.localCmd != nil {
		return cmd.localCmd.Wait()
	}

	return cmd.sshSession.Wait()
}
