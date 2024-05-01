// Run remote server commands using SSH
//
// SSH in with a password or IdentityFile to execute commands on the remote server.
package main

import (
	"context"
	"errors"
	"fmt"
)

// SSH dagger module
type Ssh struct{}

// Set configuration for SSH connections.
func (s *Ssh) Config(
	// destination to connect
	// ex) user@host
	destination string,
	// port to connect
	// +optional
	// +default=22
	port int,
) *SshConfig {
	return &SshConfig{
		Destination: destination,
		Port:        port,
		BaseCtr: dag.Container().
			From("ubuntu:22.04").
			WithExec([]string{"apt", "update"}).
			WithExec([]string{"apt", "install", "-y", "openssh-client", "sshpass"}),
	}
}

// SSH configuration
type SshConfig struct {
	// +private
	Destination string
	// +private
	Port int

	// +private
	BaseCtr *Container
}

// Set the password as the SSH connection credentials.
func (s *SshConfig) WithPassword(
	ctx context.Context,
	// password
	arg *Secret,
) (*SshCommander, error) {
	passwordText, err := arg.Plaintext(ctx)
	if err != nil {
		return nil, errors.New("invalid password secret")
	}

	return &SshCommander{
		BaseCtr:    s.BaseCtr,
		SshCommand: fmt.Sprintf(`sshpass -p %s ssh -o StrictHostKeyChecking=no -p %d %s`, passwordText, s.Port, s.Destination),
	}, nil
}

// Set up identity file with SSH connection credentials.
//
// Note: Recommend using RSA-formatted private key files. Cannot use OPENSSH-formatted private key files.
// https://github.com/dagger/dagger/issues/7220
func (s *SshConfig) WithIdentityFile(
	// identity file
	arg *Secret,
) (*SshCommander, error) {
	keyPath := "/identity_key"

	return &SshCommander{
		BaseCtr:    s.BaseCtr.WithMountedSecret(keyPath, arg),
		SshCommand: fmt.Sprintf(`ssh -i %s -o StrictHostKeyChecking=no -p %d %s`, keyPath, s.Port, s.Destination),
	}, nil
}

// SSH command launcher
type SshCommander struct {
	// +private
	BaseCtr *Container
	// +private
	SshCommand string
}

// Run the command on the remote server.
func (s *SshCommander) Command(
	// command
	arg string,
) *Container {
	exec := s.BaseCtr.WithExec([]string{
		"bash",
		"-c",
		fmt.Sprintf(`%s "%s"`, s.SshCommand, arg),
	})

	return exec
}
