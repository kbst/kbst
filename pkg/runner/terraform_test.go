package runner

import (
	"fmt"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/stretchr/testify/assert"
)

func TestNewLocalTerraformContainer(t *testing.T) {
	p := filepath.Join(fixturesPath, "multi-cloud")
	_, err := NewLocalTerraformContainer(p, false)

	assert.Equal(t, nil, err, nil)
}

func TestNewLocalTerraformContainerNoModule(t *testing.T) {
	p := filepath.Join(fixturesPath, "this-is-not-the-fixture-you-are-looking-for")
	_, err := NewLocalTerraformContainer(p, false)

	assert.Error(t, err, nil)
}

func TestBuldAndRunImageTagsMatch(t *testing.T) {
	p := filepath.Join(fixturesPath, "multi-cloud")
	ltc, _ := NewLocalTerraformContainer(p, false)
	h, _ := pathHash(p)
	expTag := fmt.Sprintf("kbst:%s-loc", h)

	buildCmd := ltc.buildCmd()
	runCmd := ltc.runCmd()

	assert.Contains(t, buildCmd.Args, expTag, nil)
	assert.Contains(t, runCmd.Args, expTag, nil)
}

var entries = map[string]string{
	"aks_zero": "github.com/kbst/terraform-kubestack//azurerm/cluster?ref=v0.10.0-beta.0",
	"eks_zero": "github.com/kbst/terraform-kubestack//aws/cluster?ref=v0.10.0-beta.0",
	"gke_zero": "github.com/kbst/terraform-kubestack//google/cluster?ref=v0.10.0-beta.0",
}

func getModuleCalls() map[string]*tfconfig.ModuleCall {
	mcs := make(map[string]*tfconfig.ModuleCall)
	for k, v := range entries {
		mcs[k] = &tfconfig.ModuleCall{Name: k, Source: v}
	}

	return mcs
}

func getExpSedArgs() []string {
	return []string{
		"-e s#github.com\\/kbst\\/terraform-kubestack\\/\\/azurerm\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/azurerm\\/cluster-local?ref=v0.10.0-beta.0#g",
		"-e s#github.com\\/kbst\\/terraform-kubestack\\/\\/aws\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/aws\\/cluster-local?ref=v0.10.0-beta.0#g",
		"-e s#github.com\\/kbst\\/terraform-kubestack\\/\\/google\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/google\\/cluster-local?ref=v0.10.0-beta.0#g",
	}
}

func TestGetSedArgs(t *testing.T) {
	ltc := localTerraformContainer{}
	sedArgs := ltc.rewriteModules(getModuleCalls())
	expSedArgs := getExpSedArgs()

	assert.Len(t, sedArgs, len(entries), nil)
	assert.ElementsMatch(t, expSedArgs, sedArgs, nil)
}

func getExpApplySh() string {
	return "\n\t#!/bin/sh\n\tset -e\n\t\n\t# disable eventual remote state\n\trm -f state.tf\n\t\n\t# replace cluster module sources with cluster-local implementation\n\tsed -i -e s#github.com\\/kbst\\/terraform-kubestack\\/\\/aws\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/aws\\/cluster-local?ref=v0.10.0-beta.0#g -e s#github.com\\/kbst\\/terraform-kubestack\\/\\/azurerm\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/azurerm\\/cluster-local?ref=v0.10.0-beta.0#g -e s#github.com\\/kbst\\/terraform-kubestack\\/\\/google\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/google\\/cluster-local?ref=v0.10.0-beta.0#g *.tf\n\t\n\tterraform init\n\t\n\tterraform workspace new loc || true\n\tterraform workspace select loc\n\t\n\tterraform apply --auto-approve\n\t"
}

func TestGetApplySh(t *testing.T) {
	ltc := localTerraformContainer{}
	applySh := ltc.renderApplySh(ltc.rewriteModules(getModuleCalls()), false)

	assert.Equal(t, getExpApplySh(), applySh, nil)
}

func getExpApplyShDestroy() string {
	return "\n\t#!/bin/sh\n\tset -e\n\t\n\t# disable eventual remote state\n\trm -f state.tf\n\t\n\t# replace cluster module sources with cluster-local implementation\n\tsed -i -e s#github.com\\/kbst\\/terraform-kubestack\\/\\/aws\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/aws\\/cluster-local?ref=v0.10.0-beta.0#g -e s#github.com\\/kbst\\/terraform-kubestack\\/\\/azurerm\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/azurerm\\/cluster-local?ref=v0.10.0-beta.0#g -e s#github.com\\/kbst\\/terraform-kubestack\\/\\/google\\/cluster?ref=v0.10.0-beta.0#github.com\\/kbst\\/terraform-kubestack\\/\\/google\\/cluster-local?ref=v0.10.0-beta.0#g *.tf\n\t\n\tterraform init\n\t\n\tterraform workspace new loc || true\n\tterraform workspace select loc\n\t\n\tterraform destroy --auto-approve\n\t"
}

func TestGetApplyShDestroy(t *testing.T) {
	ltc := localTerraformContainer{}
	applySh := ltc.renderApplySh(ltc.rewriteModules(getModuleCalls()), true)

	assert.Equal(t, getExpApplyShDestroy(), applySh, nil)
}

func TestBuildArgs(t *testing.T) {
	u, _ := user.Current()

	ltc, _ := NewLocalTerraformContainer("testpath", false)
	ba := ltc.buildArgs()

	expFile := []string{"--file", "Dockerfile.loc"}
	expTag := []string{"--tag", ltc.imageTag()}
	expBuildArg := []string{"--build-arg", fmt.Sprintf("UID=%s", u.Uid), "--build-arg", fmt.Sprintf("GID=%s", u.Gid)}

	assert.Subset(t, ba, expFile, nil)
	assert.Subset(t, ba, expTag, nil)
	assert.Subset(t, ba, expBuildArg, nil)
}

func TestRunArgs(t *testing.T) {
	p := filepath.Join(fixturesPath, "multi-cloud")
	h, _ := pathHash(p)
	ltc, _ := NewLocalTerraformContainer(p, false)
	ra := ltc.runArgs(getModuleCalls())

	expStateVolume := []string{"--volume", fmt.Sprintf("kbst-loc-terraform-state-%s:/infra/terraform.tfstate.d", h)}
	expSocketVolume := []string{"--volume", "/var/run/docker.sock:/var/run/docker.sock"}
	expImageTag := []string{fmt.Sprintf("kbst:%s-loc", h)}

	assert.Subset(t, ra, expStateVolume, nil)
	assert.Subset(t, ra, expSocketVolume, nil)
	assert.Subset(t, ra, expImageTag, nil)
}
