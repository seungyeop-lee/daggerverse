// Run remote server commands using SSH
//
// SSH in with a password or IdentityFile to execute commands on the remote server.
package main

import (
	"context"
	"dagger/ssh/internal/dagger"
	"errors"
	"fmt"
	"time"
)

// SSH dagger module
type Ssh struct{}

// Get the base container for the SSH module.
// Used when you need to inject a Service into a BaseContainer and run it.
//
// Example:
//
//	dag.SSH().Config("admin@sshd", SSHConfigOpts{
//		Port:    8022,
//		BaseCtr: dag.SSH().BaseContainer().WithServiceBinding("sshd", sshd),
//	})...
//
// Note: When a Service is passed directly as a parameter to an external Dagger function, it does not bind to the container created inside the Dagger function. (Confirmed in v0.13.3)
func (s *Ssh) BaseContainer() *dagger.Container {
	return dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "openssh-client", "sshpass"})
}

// Set configuration for SSH connections.
func (s *Ssh) Config(
	// destination to connect
	// ex) user@host
	destination string,
	// port to connect
	// +optional
	// +default=22
	port int,
	// base container
	// +optional
	baseCtr *dagger.Container,
) *SshConfig {
	if baseCtr == nil {
		baseCtr = s.BaseContainer()
	}

	return &SshConfig{
		Destination: destination,
		Port:        port,
		BaseCtr:     baseCtr,
	}
}

// SSH configuration
type SshConfig struct {
	// +private
	Destination string
	// +private
	Port int

	// +private
	BaseCtr *dagger.Container
}

// Set the password as the SSH connection credentials.
func (s *SshConfig) WithPassword(
	ctx context.Context,
	// password
	arg *dagger.Secret,
) (*SshCommander, error) {
	passwordText, err := arg.Plaintext(ctx)
	if err != nil {
		return nil, errors.New("invalid password secret")
	}

	return &SshCommander{
		BaseCtr:    s.BaseCtr,
		SshCommand: fmt.Sprintf(`sshpass -p %s ssh -o StrictHostKeyChecking=no -o LogLevel=error -p %d %s`, passwordText, s.Port, s.Destination),
	}, nil
}

// Set up identity file with SSH connection credentials.
//
// Note: Tested against RSA-formatted and OPENSSH-formatted private keys.
func (s *SshConfig) WithIdentityFile(
	// identity file
	arg *dagger.Secret,
) *SshCommander {
	keyPath := "/identity_key"

	return &SshCommander{
		BaseCtr:    s.BaseCtr.WithMountedSecret(keyPath, arg),
		SshCommand: fmt.Sprintf(`ssh -i %s -o StrictHostKeyChecking=no -o LogLevel=error -p %d %s`, keyPath, s.Port, s.Destination),
	}
}

// SSH command launcher
type SshCommander struct {
	// +private
	BaseCtr *dagger.Container
	// +private
	SshCommand string
}

// Returns a container that is ready to launch SSH command.
func (s *SshCommander) Container() *dagger.Container {
	return s.BaseCtr
}

// Run the command on the remote server.
func (s *SshCommander) Command(
	// command
	arg string,
) *dagger.Container {
	exec := s.BaseCtr.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{
			"bash",
			"-c",
			fmt.Sprintf(`%s "%s"`, s.SshCommand, arg),
		})

	return exec
}
