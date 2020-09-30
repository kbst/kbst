package util

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
)

func PathHash(path string) string {
	h := sha512.New()

	h.Write([]byte(path))

	return hex.EncodeToString(h.Sum(nil))[0:7]
}

func DockerImageTag(path string, suffix string) string {
	hash := PathHash(path)

	tag := fmt.Sprintf("kbst:%s", hash)
	if suffix != "" {
		tag = fmt.Sprintf("%s-%s", tag, suffix)
	}

	return tag
}

func DockerBuild(path string, args []string) (err error) {
	buildArgs := append([]string{"build"}, args...)
	cmd := exec.Command("docker", buildArgs...)
	cmd.Env = getEnv()
	cmd.Dir = path
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker build error: %s", err)
	}

	return
}

func DockerRun(args []string) (err error) {
	runArgs := append([]string{"run"}, args...)
	cmd := exec.Command("docker", runArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker run error: %s", err)
	}

	return
}

func getEnv() (env []string) {
	env = []string{"DOCKER_BUILDKIT=1"}
	envHome, okHome := os.LookupEnv("HOME")
	// Mac and Win require HOME to be set due to upstream bug
	// https://github.com/kbst/kbst/issues/7
	if okHome {
		env = append(env, fmt.Sprintf("HOME=%s", envHome))
	}

	// Win requires PATH for docker-credentials-desktop.exe
	envPath, okPath := os.LookupEnv("PATH")
	if okPath {
		env = append(env, fmt.Sprintf("PATH=%s", envPath))
	}

	return env
}
