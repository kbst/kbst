package runner

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

var _ TerraformContainer = &localTerraformContainer{}

type TerraformContainer interface {
	Run(destroy bool) (err error)
}

type localTerraformContainer struct {
	destroy      bool
	hash         string
	module       tfconfig.Module
	path         string
	preflight    bool
	stateVolume  string
	socketVolume string
}

func NewLocalTerraformContainer(path string) (*localTerraformContainer, error) {
	ltc := localTerraformContainer{
		path:         path,
		destroy:      false,
		preflight:    false,
		socketVolume: "/var/run/docker.sock:/var/run/docker.sock",
	}

	hash, err := pathHash(path)
	if err != nil {
		return &ltc, fmt.Errorf("path error: %s", err)
	}
	ltc.hash = hash

	ltc.stateVolume = fmt.Sprintf(
		"kbst-loc-terraform-state-%s:%s",
		ltc.hash,
		"/infra/terraform.tfstate.d",
	)

	// parse the Terraform config
	module, diags := tfconfig.LoadModule(ltc.path)
	if diags.HasErrors() {
		return &ltc, fmt.Errorf("error parsing terraform config: %s", err)
	}
	ltc.module = *module

	return &ltc, nil
}

func (ltc *localTerraformContainer) Run(destroy bool) (err error) {
	ltc.destroy = destroy
	buildCmd := ltc.buildCmd()
	err = buildCmd.Run()
	fmt.Errorf("234324docker run error: %s", err)
	if err != nil {
		return fmt.Errorf("docker build error: %s", err)
	}

	if ltc.preflight == false {
		b := bytes.NewBuffer(nil)
		w := bufio.NewWriter(b)
		defer w.Flush()

		preflightCmd := ltc.preflightCmd()
		preflightCmd.Stderr = w
		preflightCmd.Stdout = w
		err = preflightCmd.Run()
		if err != nil {
			// temp workaround for upstream issue on MacOS
			// https://github.com/docker/for-mac/issues/4755
			// if the first preflight fails, we try the raw socket
			ltc.socketVolume = "/var/run/docker.sock.raw:/var/run/docker.sock"
			preflightCmd := ltc.preflightCmd()
			err = preflightCmd.Run()
			if err != nil {
				log.Fatalf("docker preflight error:\r\n%s", b)
			}
		}
		ltc.preflight = true
	}

	applyCmd := ltc.applyCmd()
	err = applyCmd.Run()
	if err != nil {
		return fmt.Errorf("docker run error: %s", err)
	}

	return
}

func (ltc *localTerraformContainer) buildCmd() (buildCmd exec.Cmd) {
	args := ltc.buildArgs()
	return dockerBuildCommand(ltc.path, args)
}

func (ltc *localTerraformContainer) preflightCmd() (preflightCmd exec.Cmd) {
	preflightArgs := ltc.preflightArgs()
	return dockerRunCommand(preflightArgs)
}

func (ltc *localTerraformContainer) applyCmd() (applyCmd exec.Cmd) {
	// run terraform apply/destroy script inside container
	applyArgs := ltc.applyArgs(ltc.module.ModuleCalls)

	applyCmd = dockerRunCommand(applyArgs)

	return applyCmd
}

func (ltc *localTerraformContainer) imageTag() (tag string) {
	return dockerImageTag(ltc.hash, "loc")
}

func (ltc *localTerraformContainer) rewriteModules(moduleCalls map[string]*tfconfig.ModuleCall) []string {
	sedArgs := []string{}
	for _, value := range moduleCalls {
		// prepare original and replacement sources
		// escape slashes
		o := strings.Replace(value.Source, "/", "\\/", -1)
		r := strings.Replace(value.Source, "/", "\\/", -1)

		// replace cluster with cluster-local module
		r = strings.Replace(r, "/cluster?", "/cluster-local?", 1)

		// concatenate and append the sed flag
		arg := fmt.Sprintf("-e s#%s#%s#g", o, r)
		sedArgs = append(sedArgs, arg)
	}
	sort.Strings(sedArgs)
	return sedArgs
}

func (ltc *localTerraformContainer) renderApplySh(sedArgs []string, destroy bool) string {
	tfCommand := "apply"
	if destroy {
		tfCommand = "destroy"
	}

	sh := fmt.Sprintf(`
	#!/bin/sh
	set -e
	
	# disable eventual remote state
	rm -f state.tf
	
	# replace cluster module sources with cluster-local implementation
	sed -i %s *.tf
	
	terraform init
	
	terraform workspace new loc || true
	terraform workspace select loc
	
	terraform %s --auto-approve
	`, strings.Join(sedArgs, " "), tfCommand)

	return sh
}

func (ltc *localTerraformContainer) buildArgs() (buildArgs []string) {
	// get current user id to set chown during docker build
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	// build the docker image for this apply run
	buildArgs = []string{
		"--file", "Dockerfile.loc",
		"--tag", ltc.imageTag(),
		"--build-arg", fmt.Sprintf("UID=%s", u.Uid),
		"--build-arg", fmt.Sprintf("GID=%s", u.Gid),
		"."}

	return buildArgs
}

func (ltc *localTerraformContainer) defaultRunArgs() (runArgs []string) {
	runArgs = []string{
		"--rm",
		"--privileged",
		"--volume", ltc.stateVolume,
		"--volume", ltc.socketVolume,
		"--net", "host",
		ltc.imageTag()}

	return runArgs
}

func (ltc *localTerraformContainer) preflightArgs() (preflightArgs []string) {
	// run docker info as a preflight check
	preflightArgs = append(ltc.defaultRunArgs(), []string{"docker", "info"}...)

	return preflightArgs
}

func (ltc *localTerraformContainer) applyArgs(moduleCalls map[string]*tfconfig.ModuleCall) (applyArgs []string) {
	// prepare list of all module sources that need to be rewritten
	sedArgs := ltc.rewriteModules(moduleCalls)

	// render the script to run
	applySh := ltc.renderApplySh(sedArgs, ltc.destroy)

	applyArgs = append(ltc.defaultRunArgs(), []string{"sh", "-c", applySh}...)

	return applyArgs
}
