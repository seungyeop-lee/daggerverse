package main

import (
	"dagger/examples/internal/dagger"
)

type Examples struct{}

func (e *Examples) SSH_UsePassword(destination string, password string) *dagger.Container {
	return dag.SSH().
		Config(destination).
		WithPassword(dag.SetSecret("password", password)).
		Command(`echo "Hello, world!"`)
}

func (e *Examples) SSH_UsePasswordWithOption(sshd *dagger.Service, password string) *dagger.Container {
	return dag.SSH().
		Config("admin@sshd", dagger.SSHConfigOpts{
			Port:    8022,
			BaseCtr: dag.SSH().BaseContainer().WithServiceBinding("sshd", sshd),
		}).
		WithPassword(dag.SetSecret("password", password)).
		Command(`echo "Hello, world!"`)
}

func (e *Examples) SSH_UseKey(destination string, key *dagger.Secret) *dagger.Container {
	return dag.SSH().
		Config(destination).
		WithIdentityFile(key).
		Command(`echo "Hello, world!"`)
}

func (e *Examples) SSH_UseKeyWithOption(sshd *dagger.Service, key *dagger.Secret) *dagger.Container {
	return dag.SSH().
		Config("admin@sshd", dagger.SSHConfigOpts{
			Port:    8022,
			BaseCtr: dag.SSH().BaseContainer().WithServiceBinding("sshd", sshd),
		}).
		WithIdentityFile(key).
		Command(`echo "Hello, world!"`)
}
