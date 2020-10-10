package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathHashChars(t *testing.T) {
	path := "/tmp/test/path/one"
	p := PathHash(path)
	assert.Len(t, p, 7, nil)
}

func TestPathHashDestinct(t *testing.T) {
	path1 := "/tmp/test/path/one"
	p1 := PathHash(path1)

	path2 := "/tmp/test/path/two"
	p2 := PathHash(path2)

	assert.NotEqual(t, p1, p2, nil)
}

func TestDockerImageTag(t *testing.T) {
	path := "/tmp/test/path/one"
	p := PathHash(path)
	dt := DockerImageTag(path, "")

	assert.Equal(t, fmt.Sprintf("kbst:%s", p), dt, nil)
}

func TestDockerImageTagSuffix(t *testing.T) {
	path := "/tmp/test/path/one"
	p := PathHash(path)
	suffix := "test"
	dt := DockerImageTag(path, suffix)

	assert.Equal(t, fmt.Sprintf("kbst:%s-%s", p, suffix), dt, nil)
}

func TestDockerBuildCommand(t *testing.T) {
	path := "testpath"
	args := []string{"arg1", "arg2"}
	expArgs := append([]string{"docker", "build"}, args...)
	expEnv := append(os.Environ(), "DOCKER_BUILDKIT=1")

	cmd := DockerBuildCommand(path, args)

	assert.Equal(t, expArgs, cmd.Args, nil)
	assert.ElementsMatch(t, cmd.Env, expEnv, nil)
	assert.Equal(t, path, cmd.Dir, nil)
	assert.Equal(t, os.Stderr, cmd.Stderr, nil)
	assert.Equal(t, os.Stdout, cmd.Stdout, nil)
}

func TestDockerRunCommand(t *testing.T) {
	args := []string{"arg1", "arg2"}
	expArgs := append([]string{"docker", "run"}, args...)
	expEnv := append(os.Environ())

	cmd := DockerRunCommand(args)

	assert.Equal(t, expArgs, cmd.Args, nil)
	assert.ElementsMatch(t, cmd.Env, expEnv, nil)
	assert.Equal(t, os.Stderr, cmd.Stderr, nil)
	assert.Equal(t, os.Stdout, cmd.Stdout, nil)
}
