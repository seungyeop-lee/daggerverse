// Helper for private git operations
//
// A module to help you easily perform clone, push, and pull operations on your private git.

package main

import (
	"context"
	"dagger/git/internal/dagger"
	"fmt"
	"strings"
	"time"
)

// +private
const WorkDir = "/tmp/repo/"

// PrivateGit dagger module
type PrivateGit struct {
	// +private
	BaseCtr *dagger.Container
}

func New(
	// base container
	// +optional
	baseCtr *dagger.Container,
) *PrivateGit {
	git := &PrivateGit{}
	if baseCtr != nil {
		git.BaseCtr = baseCtr
	} else {
		git.BaseCtr = git.BaseContainer()
	}
	return git
}

// Get the base container for the PrivateGit module.
// Used when you need to inject a Service into a BaseContainer and run it.
//
// Example:
//
//	dag.PrivateGit(PrivateGitOpts{
//		BaseCtr: dag.PrivateGit().BaseContainer().WithServiceBinding("gitea", git),
//	})...
//
// Note: When a Service is passed directly as a parameter to an external Dagger function, it does not bind to the container created inside the Dagger function. (Confirmed in v0.13.3)
func (g *PrivateGit) BaseContainer() *dagger.Container {
	return dag.Container().
		From("ubuntu:22.04").
		WithWorkdir(WorkDir).
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "git"}).
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"git", "config", "--global", "--add", "--bool", "push.autoSetupRemote", "true"})
}

// Set up an existing repository folder.
func (g *PrivateGit) Repo(
	dir *dagger.Directory,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoDir: dir,
	}
}

// Set the ssh key.
//
// Note: Tested against RSA-formatted and OPENSSH-formatted private keys.
func (g *PrivateGit) WithSshKey(
	// ssk key file
	sshKey *dagger.Secret,
) *PrivateGitSsh {
	keyPath := "/identity_key"

	return &PrivateGitSsh{
		BaseCtr: g.BaseCtr.
			WithMountedSecret(keyPath, sshKey).
			WithEnvVariable("GIT_SSH_COMMAND", fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no -o LogLevel=error", keyPath)),
	}
}

// Set up user and password information.
func (g *PrivateGit) WithUserPassword(
	username string,
	password *dagger.Secret,
) *PrivateGitHttp {
	return &PrivateGitHttp{
		BaseCtr:  g.BaseCtr,
		Username: username,
		Password: password,
	}
}

// PrivateGit with SSH settings
type PrivateGitSsh struct {
	// +private
	BaseCtr *dagger.Container
}

// Set the SSH URL of the target repository.
func (g *PrivateGitSsh) WithRepoUrl(
	sshUrl string,
) *PrivateGitRepoUrl {
	return &PrivateGitRepoUrl{
		BaseCtr: g.BaseCtr,
		RepoUrl: sshUrl,
	}
}

// Set up an existing repository folder.
func (g *PrivateGitSsh) Repo(
	dir *dagger.Directory,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoDir: dir,
	}
}

// PrivateGit with user and password information added
type PrivateGitHttp struct {
	// +private
	BaseCtr *dagger.Container
	// +private
	Username string
	// +private
	Password *dagger.Secret
}

// Set the Web URL of the target repository.
func (g *PrivateGitHttp) WithRepoUrl(
	ctx context.Context,
	webUrl string,
) (*PrivateGitRepoUrl, error) {
	passwordPlain, err := g.Password.Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	return &PrivateGitRepoUrl{
		BaseCtr: g.BaseCtr,
		RepoUrl: strings.ReplaceAll(webUrl, "://", fmt.Sprintf("://%s:%s@", g.Username, passwordPlain)),
	}, nil
}

// Set up an existing repository folder.
func (g *PrivateGitHttp) Repo(
	dir *dagger.Directory,
) *PrivateGitRepo {
	return &PrivateGitRepo{
		BaseCtr: g.BaseCtr,
		RepoDir: dir,
	}
}

// PrivateGit with target Repositorydml URL information added
type PrivateGitRepoUrl struct {
	// +private
	BaseCtr *dagger.Container
	// +private
	RepoUrl string
}

// Clone the Git repository.
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

// Working with private Git repositories
type PrivateGitRepo struct {
	// +private
	BaseCtr *dagger.Container
	// +private
	RepoDir *dagger.Directory
}

// Returns the container with RepoDir.
func (g *PrivateGitRepo) Container() *dagger.Container {
	return g.BaseCtr.
		WithDirectory(WorkDir, g.RepoDir)
}

// Returns the repository.
func (g *PrivateGitRepo) Directory() *dagger.Directory {
	return g.RepoDir
}

// Set the user's name and email.
func (g *PrivateGitRepo) SetConfig(
	userName string,
	userEmail string,
) *PrivateGitRepo {
	g.BaseCtr = g.BaseCtr.
		WithExec([]string{"git", "config", "--global", "user.name", userName}).
		WithExec([]string{"git", "config", "--global", "user.email", userEmail})
	return g
}

// Push the repository.
func (g *PrivateGitRepo) Push() *dagger.Container {
	return g.BaseCtr.
		WithDirectory(WorkDir, g.RepoDir).
		WithExec([]string{"git", "push"})
}

// Pull the repository.
func (g *PrivateGitRepo) Pull() *dagger.Directory {
	return g.BaseCtr.
		WithDirectory(WorkDir, g.RepoDir).
		WithExec([]string{"git", "pull"}).
		Directory(WorkDir)
}
