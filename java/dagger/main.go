// Java module
//
// It is written with the assumption that it is used for Java applications that use Spring Boot and a Maven or Gradle wrapper.

package main

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

// Returns the conainer with the settings applied.
func (c *JavaConfig) Container() *Container {
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

	return ctr
}

// Run the command in the environment you set up.
func (c *JavaConfig) Run(
	// Directory to run the command in
	dir *Directory,
	// Command to run
	cmd []string,
) *Container {
	ctr := c.Container()

	return ctr.
		WithDirectory(DefaultWorkdir, dir).
		WithExec(cmd)
}
