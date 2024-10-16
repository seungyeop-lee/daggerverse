// Containers that can run [ACT](https://github.com/nektos/act)

package main

import (
	"dagger/act/internal/dagger"
)

type Act struct{}

func (a *Act) Container() *dagger.Container {
	baseContainer := dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "git", "bash", "sudo", "curl"}).
		WithExec([]string{"bash", "-c", `curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash`})

	return dag.Docker().BindEngineAsService(baseContainer)
}
