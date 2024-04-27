package main

import (
	"dagger/git/internal/dagger"
	"fmt"
	"strings"
	"time"
)

const WorkDir = "/tmp/repo/"

type PrivateGit struct {
	// +private
	BaseCtr *Container
}

func New() *PrivateGit {
	return &PrivateGit{
		BaseCtr: dag.Container().
			From("ubuntu:22.04").
			WithWorkdir(WorkDir).
			WithExec([]string{"apt", "update"}).
			WithExec([]string{"apt", "install", "-y", "git"}).
			//https://archive.docs.dagger.io/0.9/cookbook/#invalidate-cache
			WithEnvVariable("CACHE_BUSTER", time.Now().String()).
			WithExec([]string{"git", "config", "--global", "--add", "--bool", "push.autoSetupRemote", "true"}),
	}
}

func (g *PrivateGit) Push(
	dir *Directory,
) *Container {
	return push(dir, g.BaseCtr)
}

func (g *PrivateGit) Pull(
	dir *Directory,
) *Directory {
	return pull(dir, g.BaseCtr)
}

func (g *PrivateGit) WithSshKey(
	sshKey *File,
) *PrivateGitSsh {
	return &PrivateGitSsh{
		BaseCtr: g.BaseCtr.
			WithFile("/tmp/.ssh/id", sshKey, ContainerWithFileOpts{Permissions: 0400}).
			WithEnvVariable("GIT_SSH_COMMAND", "ssh -i /tmp/.ssh/id -o StrictHostKeyChecking=accept-new"),
	}
}

type PrivateGitSsh struct {
	// +private
	BaseCtr *Container
}

func (g *PrivateGitSsh) WithRepository(
	sshUrl string,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoUrl: sshUrl,
	}
}

func (g *PrivateGitSsh) Push(
	dir *Directory,
) *Container {
	return push(dir, g.BaseCtr)
}

func (g *PrivateGitSsh) Pull(
	dir *Directory,
) *Directory {
	return pull(dir, g.BaseCtr)
}

func (g *PrivateGit) WithUserPassword(
	username string,
	password string,
) *PrivateGitHttp {
	return &PrivateGitHttp{
		BaseCtr:  g.BaseCtr,
		Username: username,
		Password: password,
	}
}

type PrivateGitHttp struct {
	// +private
	BaseCtr *Container
	// +private
	Username string
	// +private
	Password string
}

func (g *PrivateGitHttp) WithRepository(
	webUrl string,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoUrl: strings.ReplaceAll(webUrl, "://", fmt.Sprintf("://%s:%s@", g.Username, g.Password)),
	}
}

type PrivateGitRepo struct {
	// +private
	BaseCtr *Container
	// +private
	RepoUrl string
}

func (g *PrivateGitRepo) Clone() *Directory {
	return g.BaseCtr.
		WithExec([]string{"git", "clone", g.RepoUrl, "."}).
		Directory(WorkDir)
}

func (g *PrivateGitRepo) Config(
	userName string,
	userEmail string,
) *PrivateGitRepoWork {
	return &PrivateGitRepoWork{
		BaseCtr: g.BaseCtr.
			WithExec([]string{"git", "config", "--global", "user.name", userName}).
			WithExec([]string{"git", "config", "--global", "user.email", userEmail}),
	}
}

type PrivateGitRepoWork struct {
	// +private
	BaseCtr *Container
}

func push(dir *Directory, c *Container) *dagger.Container {
	return c.
		WithDirectory(WorkDir, dir).
		WithExec([]string{"git", "push"})
}

func pull(dir *Directory, c *Container) *dagger.Directory {
	return c.
		WithDirectory(WorkDir, dir).
		WithExec([]string{"git", "pull"}).
		Directory(WorkDir)
}
