package main

import "dagger/examples/internal/dagger"

type Examples struct{}

func (e *Examples) Java_BuildByMaven(dir *dagger.Directory) *dagger.File {
	return dag.Java().
		Init().
		WithMavenCache().
		WithDir(dir).
		Run([]string{"./mvnw", "package"}).
		File("jar/dagger-maven.jar")
}

func (e *Examples) Java_BuildByGradle(dir *dagger.Directory) *dagger.File {
	return dag.Java().
		Init().
		WithGradleCache().
		WithDir(dir).
		Run([]string{"./gradlew", "bootJar"}).
		File("jar/dagger-gradle.jar")
}
