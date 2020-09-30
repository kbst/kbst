package util

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	cmd.Env = []string{"DOCKER_BUILDKIT=0"}
	if runtime.GOOS == "darwin" {
		// temp workaround can be removed once upstream fix is released
		// https://github.com/kbst/kbst/issues/7
		envHome, exists := os.LookupEnv("HOME")
		if exists {
			cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", envHome))
		}
	}
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
