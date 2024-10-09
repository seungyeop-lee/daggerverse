package main

import "dagger/examples/internal/dagger"

type Examples struct{}

func (e *Examples) Scp_CopyToRemote(destination string, key *dagger.Secret, file *dagger.File) *dagger.Container {
	return dag.Scp().
		Config(destination).
		WithIdentityFile(key).
		FileToRemote(file)
}

func (e *Examples) Scp_CopyToRemoteWithOption(sshd *dagger.Service, key *dagger.Secret, file *dagger.File, target string) *dagger.Container {
	return dag.Scp().
		Config("admin@sshd", dagger.ScpConfigOpts{
			Port:    8022,
			BaseCtr: dag.Scp().BaseContainer().WithServiceBinding("sshd", sshd),
		}).
		WithIdentityFile(key).
		FileToRemote(file, dagger.ScpCommanderFileToRemoteOpts{
			Target: target,
		})
}

func (e *Examples) Scp_CopyFromRemote(destination string, key *dagger.Secret, path string) *dagger.File {
	return dag.Scp().
		Config(destination).
		WithIdentityFile(key).
		FileFromRemote(path)
}

func (e *Examples) Scp_CopyFromRemoteWithOption(sshd *dagger.Service, key *dagger.Secret, path string) *dagger.File {
	return dag.Scp().
		Config("admin@sshd", dagger.ScpConfigOpts{
			Port:    8022,
			BaseCtr: dag.Scp().BaseContainer().WithServiceBinding("sshd", sshd),
		}).
		WithIdentityFile(key).
		FileFromRemote(path)
}

func (e *Examples) Scp_UsePassword(destination string, password string, file *dagger.File) *dagger.Container {
	return dag.Scp().
		Config(destination).
		WithPassword(dag.SetSecret("password", password)).
		FileToRemote(file)
}

func (e *Examples) Scp_CopyDirectoryToRemote(destination string, key *dagger.Secret, dir *dagger.Directory, target string) *dagger.Container {
	return dag.Scp().
		Config(destination).
		WithIdentityFile(key).
		DirectoryToRemote(dir, target)
}
