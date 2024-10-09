// Copy to/from remote server using SCP
//
// Performs copying of files and directories to and from a remote server over SCP using password or IdentityFile.
package main

import (
	"context"
	"dagger/scp/internal/dagger"
	"errors"
	"path"
	"strconv"
	"time"
)

// SCP dagger module
type Scp struct{}

// Get the base container for the SCP module.
// Used when you need to inject a Service into a BaseContainer and run it.
//
// Example:
//
//	dag.Scp().Config("admin@sshd", ScpConfigOpts{
//		Port:    8022,
//		BaseCtr: dag.Scp().BaseContainer().WithServiceBinding("sshd", sshd),
//	})...
//
// Note: When a Service is passed directly as a parameter to an external Dagger function, it does not bind to the container created inside the Dagger function. (Confirmed in v0.13.3)
func (s *Scp) BaseContainer() *dagger.Container {
	return dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "openssh-client", "sshpass"})
}

// Set configuration for SCP connections.
func (s *Scp) Config(
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
) *ScpConfig {
	if baseCtr == nil {
		baseCtr = s.BaseContainer()
	}

	return &ScpConfig{
		Destination: destination,
		Port:        port,
		BaseCtr:     baseCtr,
	}
}

// SCP configuration
type ScpConfig struct {
	// +private
	Destination string
	// +private
	Port int
	// +private
	BaseCtr *dagger.Container
}

// Set the password as the SCP connection credentials.
func (s *ScpConfig) WithPassword(
	ctx context.Context,
	// password
	arg *dagger.Secret,
) (*ScpCommander, error) {
	passwordText, err := arg.Plaintext(ctx)
	if err != nil {
		return nil, errors.New("invalid password secret")
	}

	return &ScpCommander{
		Destination: s.Destination,
		BaseCtr:     s.BaseCtr,
		ScpBaseCommand: []string{
			"sshpass",
			"-p", passwordText,
			"scp",
			"-o", "StrictHostKeyChecking=no",
			"-o", "LogLevel=error",
			"-P", strconv.Itoa(s.Port),
		},
	}, nil
}

// Set up identity file with SCP connection credentials.
//
// Note: Tested against RSA-formatted and OPENSSH-formatted private keys.
func (s *ScpConfig) WithIdentityFile(
	// identity file
	arg *dagger.Secret,
) (*ScpCommander, error) {
	keyPath := "/identity_key"

	return &ScpCommander{
		Destination: s.Destination,
		BaseCtr:     s.BaseCtr.WithMountedSecret(keyPath, arg),
		ScpBaseCommand: []string{
			"scp",
			"-i", keyPath,
			"-o", "StrictHostKeyChecking=no",
			"-o", "LogLevel=error",
			"-P", strconv.Itoa(s.Port),
		},
	}, nil
}

// SCP command launcher
type ScpCommander struct {
	// +private
	Destination string
	// +private
	BaseCtr *dagger.Container
	// +private
	ScpBaseCommand []string
}

// Returns a container that is ready to launch SCP command.
func (s *ScpCommander) Container() *dagger.Container {
	return s.BaseCtr
}

// Copy a file to a remote server.
func (s *ScpCommander) FileToRemote(
	ctx context.Context,
	// source file
	source *dagger.File,
	// destination path
	// (If not entered, '.' is used as the default)
	// +optional
	target string,
) (*dagger.Container, error) {
	if target == "" {
		target = "."
	}

	name, err := source.Name(ctx)
	if err != nil {
		return nil, err
	}

	return s.BaseCtr.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithFile(name, source).
		WithExec(append(s.ScpBaseCommand, name, s.Destination+":"+target)), nil
}

// Copy a file from a remote server.
func (s *ScpCommander) FileFromRemote(
	// source path
	source string,
) *dagger.File {
	_, file := path.Split(source)

	return s.BaseCtr.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec(append(s.ScpBaseCommand, s.Destination+":"+source, file)).
		File(file)
}

// Copy a directory to a remote server.
func (s *ScpCommander) DirectoryToRemote(
	// source directory
	source *dagger.Directory,
	// destination path
	// (If the path is an already existing directory, it will be copied to the '[path]/source-dir' location)
	target string,
) (*dagger.Container, error) {
	sourcePath := "/source-dir"

	return s.BaseCtr.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithDirectory(sourcePath, source).
		WithExec(append(s.ScpBaseCommand, "-r", sourcePath, s.Destination+":"+target)), nil
}

// Copy a directory from a remote server.
func (s *ScpCommander) DirectoryFromRemote(
	// source path
	source string,
) *dagger.Directory {
	targetPath := "/target-dir"

	return s.BaseCtr.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec(append(s.ScpBaseCommand, "-r", s.Destination+":"+source, targetPath)).
		Directory(targetPath)
}
