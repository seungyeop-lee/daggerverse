// Java module
//
// It is written with the assumption that it is used for Java applications that use Spring Boot and a Maven or Gradle wrapper.

package main

import "errors"

const DefaultImage = "eclipse-temurin:21-jdk"
const DefaultWorkdir = "/app"

// Java dagger module
type Java struct {
}

// Initializes a JavaConfig with default settings.
func (j Java) Init() *JavaConfig {
	return &JavaConfig{
		Image:       DefaultImage,
		GradleCache: false,
		MavenCache:  false,
	}
}

// JavaConfig struct definition
// This struct holds the configuration for Java applications.
type JavaConfig struct {
	// +private
	Image string
	// +private
	GradleCache bool
	// +private
	MavenCache bool
	// +private
	Dir *Directory
}

// Sets a custom Docker image for the Java application.
func (c *JavaConfig) WithImage(image string) *JavaConfig {
	c.Image = image
	return c
}

// Enables the Gradle cache.
func (c *JavaConfig) WithGradleCache() *JavaConfig {
	c.GradleCache = true
	return c
}

// Enables the Maven cache.
func (c *JavaConfig) WithMavenCache() *JavaConfig {
	c.MavenCache = true
	return c
}

// Sets the directory to run the command in.
func (c *JavaConfig) WithDir(dir *Directory) *JavaConfig {
	c.Dir = dir
	return c
}

// Returns the conainer with the settings applied.
func (c *JavaConfig) Container() (*Container, error) {
	ctr := dag.Container().
		From(DefaultImage).
		WithWorkdir(DefaultWorkdir)

	if c.GradleCache {
		ctr = ctr.
			WithMountedCache(DefaultWorkdir+"/build", dag.CacheVolume("java-app-build-cache")).
			WithMountedCache(DefaultWorkdir+"/.gradle", dag.CacheVolume("java-app-gradle-cache")).
			WithMountedCache("/root/.gradle", dag.CacheVolume("java-root-gradle-cache"))
	}

	if c.MavenCache {
		ctr = ctr.
			WithMountedCache(DefaultWorkdir+"/target", dag.CacheVolume("java-app-target-cache")).
			WithMountedCache("/root/.m2", dag.CacheVolume("java-root-maven-cache"))
	}

	if c.Dir != nil {
		ctr = ctr.WithDirectory(DefaultWorkdir, c.Dir)
	} else {
		return nil, errors.New("dir is required")
	}

	return ctr, nil
}

// Run the command in the environment you set up.
func (c *JavaConfig) Run(
	// Command to run
	cmd []string,
) (*Container, error) {
	ctr, err := c.Container()
	if err != nil {
		return nil, err
	}

	return ctr.WithExec(cmd), nil
}
