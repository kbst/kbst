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
	cmd.Env = []string{"DOCKER_BUILDKIT=1"}
	cmd.Dir = path
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
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
		return err
	}
	return
}
