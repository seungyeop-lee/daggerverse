package main

import "dagger/examples/internal/dagger"

type Examples struct{}

func (e *Examples) PrivateGit_CloneByHttp(git *dagger.Service) *dagger.Directory {
	return dag.
		PrivateGit(dagger.PrivateGitOpts{
			BaseCtr: dag.PrivateGit().BaseContainer().WithServiceBinding("gitea", git),
		}).
		WithUserPassword("super", dag.SetSecret("password", "super")).
		WithRepoURL("http://gitea:3000/super/test.git").
		Clone().
		Directory()
}

func (e *Examples) PrivateGit_CloneBySsh(git *dagger.Service, key *dagger.Secret) *dagger.Directory {
	return dag.
		PrivateGit(dagger.PrivateGitOpts{
			BaseCtr: dag.PrivateGit().BaseContainer().WithServiceBinding("gitea", git),
		}).
		WithSSHKey(key).
		WithRepoURL("git@gitea:super/test.git").
		Clone().
		Directory()
}

func (e *Examples) PrivateGit_Push(git *dagger.Service, repo *dagger.Directory) *dagger.Container {
	return dag.
		PrivateGit(dagger.PrivateGitOpts{
			BaseCtr: dag.PrivateGit().BaseContainer().WithServiceBinding("gitea", git),
		}).
		WithUserPassword("super", dag.SetSecret("password", "super")).
		Repo(repo).
		Push()
}
