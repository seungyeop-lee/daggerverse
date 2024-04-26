package main

import (
	"context"
	"errors"
	"path"
	"strconv"
)

type Scp struct {
	// +private
	Destination string
	// +private
	Port int

	// +private
	BaseCtr *Container
	// +private
	ScpBaseCommand []string
}

func (s *Scp) Config(
	// destination to connect
	destination string,
	// port to connect
	// +optional
	// +default=22
	port int,
) (*Scp, error) {
	s.Destination = destination
	s.Port = port

	s.BaseCtr = dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "openssh-client", "sshpass"})

	return s, nil
}

func (s *Scp) WithPassword(
	ctx context.Context,
	// password
	arg *Secret,
) (*Scp, error) {
	passwordText, err := arg.Plaintext(ctx)
	if err != nil {
		return nil, errors.New("invalid password secret")
	}
	s.ScpBaseCommand = []string{
		"sshpass",
		"-p", passwordText,
		"scp",
		"-o", "StrictHostKeyChecking=no",
		"-P", strconv.Itoa(s.Port),
	}

	return s, nil
}

func (s *Scp) WithIdentityFile(
	// identity file
	arg *Secret,
) (*Scp, error) {
	keyPath := "/identity_key"
	s.BaseCtr = s.BaseCtr.WithMountedSecret(keyPath, arg)
	s.ScpBaseCommand = []string{
		"scp",
		"-i", keyPath,
		"-o", "StrictHostKeyChecking=no",
		"-P", strconv.Itoa(s.Port),
	}

	return s, nil
}

func (s *Scp) FileToRemote(
	ctx context.Context,
	source *File,
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

func (s *Scp) FileFromRemote(
	source string,
) *File {
	_, file := path.Split(source)

	return s.BaseCtr.
		WithExec(append(s.ScpBaseCommand, s.Destination+":"+source, file)).
		File(file)
}

func (s *Scp) DirectoryToRemote(
	source *Directory,
	target string,
) (*Container, error) {
	sourcePath := "/source-dir"

	return s.BaseCtr.
		WithDirectory(sourcePath, source).
		WithExec(append(s.ScpBaseCommand, "-r", sourcePath, s.Destination+":"+target)), nil
}

func (s *Scp) DirectoryFromRemote(
	source string,
) *Directory {
	targetPath := "/target-dir"

	return s.BaseCtr.
		WithExec(append(s.ScpBaseCommand, "-r", s.Destination+":"+source, targetPath)).
		Directory(targetPath)
}
