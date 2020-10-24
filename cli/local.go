package cli

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kbst/kbst/util"
)

type Local struct {
	Runner  TerraformContainer
	Watcher Watcher
}

func (l *Local) Apply(path string, skipWatch bool) (err error) {
	// provision the development environment
	err = l.Runner.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("provisioning local environment error: %s", err))
	}

	if skipWatch {
		return
	}

	// start watching for repository changes
	run := l.Watcher.Start(path)

	for {
		log.Println("#### Watching for changes")
		<-run
		err = l.Runner.Run()
		if err != nil {
			return errors.New(fmt.Sprintf("updating local environment error: %s", err))
		}
	}
}

func (l *Local) Destroy() (err error) {
	err = l.Runner.Run()
	if err != nil {
		return err
	}

	return
}

type TerraformContainer interface {
	Run() (err error)
}

type LocalTerraformContainer struct {
	destroy bool
	hash    string
	module  *tfconfig.Module
	path    string
}

func NewLocalTerraformContainer(path string, destroy bool) (ltc LocalTerraformContainer, err error) {
	ltc.path = path
	ltc.destroy = destroy

	hash, err := util.PathHash(path)
	if err != nil {
		return ltc, fmt.Errorf("path error: %s", err)
	}
	ltc.hash = hash

	// parse the Terraform config
	module, diags := tfconfig.LoadModule(ltc.path)
	if diags.HasErrors() {
		return ltc, fmt.Errorf("error parsing terraform config: %s", err)
	}
	ltc.module = module

	return ltc, nil
}

func (ltc *LocalTerraformContainer) Run() (err error) {
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

func (ltc *LocalTerraformContainer) buildCmd() (buildCmd exec.Cmd) {
	args := ltc.buildArgs()
	return util.DockerBuildCommand(ltc.path, args)
}

func (ltc *LocalTerraformContainer) runCmd() (runCmd exec.Cmd) {
	// run terraform apply/destroy script inside container
	runArgs := ltc.runArgs(ltc.module.ModuleCalls)

	runCmd = util.DockerRunCommand(runArgs)

	return runCmd
}

func (ltc *LocalTerraformContainer) imageTag() (tag string) {
	return util.DockerImageTag(ltc.hash, "loc")
}

func (ltc *LocalTerraformContainer) rewriteModules(moduleCalls map[string]*tfconfig.ModuleCall) []string {
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

func (ltc *LocalTerraformContainer) renderApplySh(sedArgs []string, destroy bool) string {
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

func (ltc *LocalTerraformContainer) buildArgs() (buildArgs []string) {
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

	return buildArgs
}

func (ltc *LocalTerraformContainer) runArgs(moduleCalls map[string]*tfconfig.ModuleCall) (runArgs []string) {
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
	socketVolume := "/var/run/docker.sock:/var/run/docker.sock"

	runArgs = []string{
		"--rm",
		"--privileged",
		"--volume", stateVolume,
		"--volume", socketVolume,
		"--net", "host",
		ltc.imageTag(),
		"sh", "-c", applySh}

	return runArgs
}
