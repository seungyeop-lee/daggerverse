// Copy to/from remote server using SCP
//
// Performs copying of files and directories to and from a remote server over SCP using password or IdentityFile.
package main

import (
	"context"
	"errors"
	"path"
	"strconv"
)

// SCP dagger module
type Scp struct{}

// Set configuration for SCP connections.
func (s *Scp) Config(
	// destination to connect
	// ex) user@host
	destination string,
	// port to connect
	// +optional
	// +default=22
	port int,
) (*ScpConfig, error) {
	return &ScpConfig{
		Destination: destination,
		Port:        port,
		BaseCtr: dag.Container().
			From("ubuntu:22.04").
			WithExec([]string{"apt", "update"}).
			WithExec([]string{"apt", "install", "-y", "openssh-client", "sshpass"}),
	}, nil
}

// SCP configuration
type ScpConfig struct {
	// +private
	Destination string
	// +private
	Port int
	// +private
	BaseCtr *Container
}

// Set the password as the SCP connection credentials.
func (s *ScpConfig) WithPassword(
	ctx context.Context,
	// password
	arg *Secret,
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
			"-P", strconv.Itoa(s.Port),
		},
	}, nil
}

// Set up identity file with SCP connection credentials.
func (s *ScpConfig) WithIdentityFile(
	// identity file
	arg *Secret,
) (*ScpCommander, error) {
	keyPath := "/identity_key"

	return &ScpCommander{
		Destination: s.Destination,
		BaseCtr:     s.BaseCtr.WithMountedSecret(keyPath, arg),
		ScpBaseCommand: []string{
			"scp",
			"-i", keyPath,
			"-o", "StrictHostKeyChecking=no",
			"-P", strconv.Itoa(s.Port),
		},
	}, nil
}

// SCP command launcher
type ScpCommander struct {
	// +private
	Destination string
	// +private
	BaseCtr *Container
	// +private
	ScpBaseCommand []string
}

// Copy a file to a remote server.
func (s *ScpCommander) FileToRemote(
	ctx context.Context,
	// source file
	source *File,
	// destination path
	// (If not entered, '.' is used as the default)
	// +optional
	target string,
) (*Container, error) {
	if target == "" {
		target = "."
	}

	name, err := source.Name(ctx)
	if err != nil {
		return nil, err
	}

	return s.BaseCtr.
		WithFile(name, source).
		WithExec(append(s.ScpBaseCommand, name, s.Destination+":"+target)), nil
}

// Copy a file from a remote server.
func (s *ScpCommander) FileFromRemote(
	// source path
	source string,
) *File {
	_, file := path.Split(source)

	return s.BaseCtr.
		WithExec(append(s.ScpBaseCommand, s.Destination+":"+source, file)).
		File(file)
}

// Copy a directory to a remote server.
func (s *ScpCommander) DirectoryToRemote(
	// source directory
	source *Directory,
	// destination path
	// (If the path is an already existing directory, it will be copied to the '[path]/source-dir' location)
	target string,
) (*Container, error) {
	sourcePath := "/source-dir"

	return s.BaseCtr.
		WithDirectory(sourcePath, source).
		WithExec(append(s.ScpBaseCommand, "-r", sourcePath, s.Destination+":"+target)), nil
}

// Copy a directory from a remote server.
func (s *ScpCommander) DirectoryFromRemote(
	// source path
	source string,
) *Directory {
	targetPath := "/target-dir"

	return s.BaseCtr.
		WithExec(append(s.ScpBaseCommand, "-r", s.Destination+":"+source, targetPath)).
		Directory(targetPath)
}
