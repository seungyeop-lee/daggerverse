// Containers that can run [ACT](https://github.com/nektos/act)

package main

type Act struct{}

func (a *Act) Container() *Container {
	actContainer := dag.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "git", "bash", "sudo", "curl"}).
		WithExec([]string{"bash", "-c", `curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash`})

	withDocker := dag.Docker().
		BindAsService(actContainer)

	return withDocker
}
