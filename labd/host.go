package main

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Host struct {
	Name     string
	Host     string
	Username string
	KeyPath  string
	sshConf  *ssh.ClientConfig
	sshConn  *ssh.Client
}

func (h *Host) StartCmd(cmdStr string, args ...string) (*Command, error) {
	var err error
	if h.Host == "localhost" && h.Username == "" {
		cmd, err := StartLocalCommand(cmdStr, args...)
		if err != nil {
			return nil, err
		}

		return cmd, nil
	}

	if h.sshConf == nil {
		err = h.generateSSHConfig()
		if err != nil {
			return nil, err
		}
	}

	if h.sshConn == nil {
		h.sshConn, err = ssh.Dial("tcp", h.Host, h.sshConf)
		if err != nil {
			return nil, err
		}
	}

	session, err := h.sshConn.NewSession()
	if err != nil {
		return nil, err
	}

	cmd, err := StartSSHCommand(session, cmdStr, args...)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (h *Host) generateSSHConfig() error {
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return err
	}

	agent := agent.NewClient(sock)
	signers, err := agent.Signers()
	if err != nil {
		return err
	}

	auths := []ssh.AuthMethod{ssh.PublicKeys(signers...)}

	h.sshConf = &ssh.ClientConfig{
		User: h.Username,
		Auth: auths,
	}

	return nil
}
