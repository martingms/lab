package main

import (
	"errors"

	"golang.org/x/crypto/ssh"
)

type Host struct {
	Name     string
	Host     string
	Username string
	KeyPath  string
	sshConf  *ssh.ClientConfig
}

func (h *Host) StartCmd(cmdStr string, args ...string) (*Command, error) {
	if h.Host == "localhost" {
		cmd, err := LocalCommand(cmdStr, args...)
        if err != nil {
            return nil, err
        }

        return cmd, nil
    }

	if h.sshConf == nil {
		err := h.generateSSHConfig()
		if err != nil {
			return nil, err
		}
	}

    return nil, errors.New("Remote commands not yet implemented")
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
