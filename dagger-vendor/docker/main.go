package main

import (
	"context"
	"dagger/docker/internal/dagger"
	"fmt"
	"strings"
)

// Docker represents the Docker module for Dagger.
type Docker struct{}

// Spawn an ephemeral Docker Engine in a container
func (d *Docker) Engine(
	// Docker Engine version
	//
	// +optional
	// +default="26.1"
	version string,

	// Persist the state of the engine in a cache volume
	//
	// +optional
	// +default=true
	persist bool,

	// Namespace for persisting the engine state. Use in combination with `persist`
	//
	// +optional
	namespace string,
) *dagger.Service {
	ctr := dag.Container().From(fmt.Sprintf("index.docker.io/docker:%s-dind", version))

	// disable the entrypoint
	ctr = ctr.WithoutEntrypoint()

	// expose the Docker Engine API
	ctr = ctr.WithExposedPort(2375)

	// persist the engine state
	if persist {
		var (
			name   = strings.TrimSuffix("docker-engine-state-"+version+"-"+namespace, "-")
			volume = dag.CacheVolume(name)
			opts   = dagger.ContainerWithMountedCacheOpts{Sharing: dagger.Locked}
		)

		ctr = ctr.WithMountedCache("/var/lib/docker", volume, opts)
	}

	return ctr.
		WithExec(
			[]string{
				"dockerd",
				"--host=tcp://0.0.0.0:2375",
				"--host=unix:///var/run/docker.sock",
				"--tls=false",
			},
			dagger.ContainerWithExecOpts{InsecureRootCapabilities: true},
		).
		AsService()
}

func (d *Docker) BindEngineAsService(
	ctx context.Context,
	target *dagger.Container,
	// Docker Engine version
	//
	// +optional
	// +default="26.1"
	version string,

	// Persist the state of the engine in a cache volume
	//
	// +optional
	// +default=true
	persist bool,

	// Namespace for persisting the engine state. Use in combination with `persist`
	//
	// +optional
	namespace string,
) (*dagger.Container, error) {
	// convert the container to a service.
	dockerService := d.Engine(version, persist, namespace)

	// get the endpoint of the service to set the DOCKER_HOST environment variable. The reason we're not using the
	// alias for docker is because the service alias is not available in the child containers of the container.
	endpoint, err := dockerService.Endpoint(ctx, dagger.ServiceEndpointOpts{Scheme: "tcp"})
	if err != nil {
		return nil, err
	}

	// bind the service to the container and set the DOCKER_HOST environment variable.
	return target.WithServiceBinding("docker", dockerService).WithEnvVariable("DOCKER_HOST", endpoint), nil
}
