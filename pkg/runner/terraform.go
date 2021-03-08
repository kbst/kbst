package runner

import (
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
	destroy bool
	hash    string
	module  tfconfig.Module
	path    string
}

func NewLocalTerraformContainer(path string) (*localTerraformContainer, error) {
	ltc := localTerraformContainer{
		path:    path,
		destroy: false,
	}

	hash, err := pathHash(path)
	if err != nil {
		return &ltc, fmt.Errorf("path error: %s", err)
	}
	ltc.hash = hash

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
	if err != nil {
		return fmt.Errorf("docker build error: %s", err)
	}

	runCmd := ltc.runCmd()
	err = runCmd.Run()
	if err != nil {
		return fmt.Errorf("docker run error: %s", err)
	}

	return
}

func (ltc *localTerraformContainer) buildCmd() (buildCmd exec.Cmd) {
	args := ltc.buildArgs()
	return dockerBuildCommand(ltc.path, args)
}

func (ltc *localTerraformContainer) runCmd() (runCmd exec.Cmd) {
	// run terraform apply/destroy script inside container
	runArgs := ltc.runArgs(ltc.module.ModuleCalls)

	runCmd = dockerRunCommand(runArgs)

	return runCmd
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

	tag := ltc.imageTag()

	// build the docker image for this apply run
	buildArgs = []string{
		"--file", "Dockerfile.loc",
		"--tag", tag,
		"--build-arg", fmt.Sprintf("UID=%s", u.Uid),
		"--build-arg", fmt.Sprintf("GID=%s", u.Gid),
		"."}
    log.Println(buildArgs)

	return buildArgs
}

func (ltc *localTerraformContainer) runArgs(moduleCalls map[string]*tfconfig.ModuleCall) (runArgs []string) {
	// prepare list of all module sources that need to be rewritten
	sedArgs := ltc.rewriteModules(moduleCalls)

	// render the script to run
	applySh := ltc.renderApplySh(sedArgs, ltc.destroy)

	// prepare volumes
	stateVolume := fmt.Sprintf(
		"kbst-loc-terraform-state-%s:%s",
		ltc.hash,
		"/infra/terraform.tfstate.d",
	)
	socketVolume := "/var/run/docker.sock.raw:/var/run/docker.sock"

	runArgs = []string{
		"--rm",
		"--privileged",
		"--volume", stateVolume,
		"--volume", socketVolume,
		"--net", "host",
		ltc.imageTag(),
		"sh", "-c", applySh}

    log.Println(runArgs)

	return runArgs
}
