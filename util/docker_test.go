package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerImageTag(t *testing.T) {
	h := "testhash123"
	dt := DockerImageTag(h, "")

	assert.Equal(t, fmt.Sprintf("kbst:%s", h), dt, nil)
}

func TestDockerImageTagSuffix(t *testing.T) {
	h := "testhash123"
	s := "test"
	dt := DockerImageTag(h, s)

	assert.Equal(t, fmt.Sprintf("kbst:%s-%s", h, s), dt, nil)
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
