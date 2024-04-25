package main

import (
	"context"
	"errors"
	"fmt"
)

type Ssh struct {
	// +private
	Destination string
	// +private
	Port int

	// +private
	BaseCtr *Container
	// +private
	SshCommand string
}

func (s *Ssh) Config(
	// destination to connect
	destination string,
	// port to connect
	// +optional
	// +default=22
	port int,
) (*Ssh, error) {
	s.Destination = destination
	s.Port = port

	s.BaseCtr = dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "openssh-client", "sshpass"})

	return s, nil
}

func (s *Ssh) WithPassword(
	ctx context.Context,
	// password
	arg *Secret,
) (*Ssh, error) {
	passwordText, err := arg.Plaintext(ctx)
	if err != nil {
		return nil, errors.New("invalid password secret")
	}
	s.SshCommand = fmt.Sprintf(`sshpass -p %s ssh -o StrictHostKeyChecking=no -p %d %s`, passwordText, s.Port, s.Destination)

	return s, nil
}

func (s *Ssh) WithIdentityFile(
	// identity file
	arg *Secret,
) (*Ssh, error) {
	keyPath := "/identity_key"
	s.BaseCtr = s.BaseCtr.WithMountedSecret(keyPath, arg)
	s.SshCommand = fmt.Sprintf(`ssh -i %s -o StrictHostKeyChecking=no -p %d %s`, keyPath, s.Port, s.Destination)

	return s, nil
}

func (s *Ssh) Command(
	arg string,
) *Container {
	exec := s.BaseCtr.WithExec([]string{
		"bash",
		"-c",
		fmt.Sprintf(`%s "%s"`, s.SshCommand, arg),
	})

	return exec
}
