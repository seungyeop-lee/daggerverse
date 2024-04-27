package main

import (
	"context"
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

func (g *PrivateGit) Repo(
	dir *Directory,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoDir: dir,
	}
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

type PrivateGitSsh struct {
	// +private
	BaseCtr *Container
}

func (g *PrivateGitSsh) WithRepoUrl(
	sshUrl string,
) *PrivateGitRepoUrl {
	return &PrivateGitRepoUrl{
		BaseCtr: g.BaseCtr,
		RepoUrl: sshUrl,
	}
}

func (g *PrivateGitSsh) Repo(
	dir *Directory,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoDir: dir,
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

func (g *PrivateGitHttp) WithRepoUrl(
	webUrl string,
) *PrivateGitRepoUrl {
	return &PrivateGitRepoUrl{
		BaseCtr: g.BaseCtr,
		RepoUrl: strings.ReplaceAll(webUrl, "://", fmt.Sprintf("://%s:%s@", g.Username, g.Password)),
	}
}

type PrivateGitRepoUrl struct {
	// +private
	BaseCtr *Container
	// +private
	RepoUrl string
}

func (g *PrivateGitRepoUrl) Clone(
	ctx context.Context,
) (*PrivateGitRepo, error) {
	repoDir, err := g.BaseCtr.
		WithExec([]string{"git", "clone", g.RepoUrl, "."}).
		Directory(WorkDir).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoDir: repoDir,
	}, nil
}

type PrivateGitRepo struct {
	// +private
	BaseCtr *Container
	// +private
	RepoDir *Directory
}

func (g *PrivateGitRepo) Directory() *Directory {
	return g.RepoDir
}

func (g *PrivateGitRepo) SetConfig(
	userName string,
	userEmail string,
) *PrivateGitRepo {
	g.BaseCtr = g.BaseCtr.
		WithExec([]string{"git", "config", "--global", "user.name", userName}).
		WithExec([]string{"git", "config", "--global", "user.email", userEmail})
	return g
}

func (g *PrivateGitRepo) Push() *Container {
	return g.BaseCtr.
		WithDirectory(WorkDir, g.RepoDir).
		WithExec([]string{"git", "push"})
}

func (g *PrivateGitRepo) Pull() *Directory {
	return g.BaseCtr.
		WithDirectory(WorkDir, g.RepoDir).
		WithExec([]string{"git", "pull"}).
		Directory(WorkDir)
}
