package cli

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kbst/kbst/util"
	"gopkg.in/fsnotify.v1"
)

type applyLock struct {
	mux sync.Mutex
}

type lastEvent struct {
	ts  time.Time
	mux sync.Mutex
}

func (l *lastEvent) Set(ts time.Time) {
	l.mux.Lock()
	l.ts = ts
	l.mux.Unlock()
}

func (l *lastEvent) Get() time.Time {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.ts
}

func DevApply(path string, skipWatch bool) (err error) {
	applyLock := applyLock{}
	lastEvent := lastEvent{}

	// first apply to bring up dev env
	ts := time.Now()
	lastEvent.Set(ts)
	runLocalTerraformContainer(path, false, ts, &lastEvent, &applyLock, skipWatch)

	if skipWatch {
		return
	}

	// then start watching
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("test %s", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}

				ts := time.Now()
				lastEvent.Set(ts)
				go runLocalTerraformContainer(path, false, ts, &lastEvent, &applyLock, skipWatch)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	basePath := filepath.Dir(path)
	watchTargets := []string{
		".",
		"manifests/bases",
		"manifests/overlays/apps",
		"manifests/overlays/ops",
		"manifests/overlays/loc",
	}
	for i := range watchTargets {
		fullPath := filepath.Join(basePath, watchTargets[i])
		err = watcher.Add(fullPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
	return
}

func DevDestroy(path string) (err error) {
	applyLock := applyLock{}
	lastEvent := lastEvent{}

	// first apply to bring up dev env
	ts := time.Now()
	lastEvent.Set(ts)
	runLocalTerraformContainer(path, true, ts, &lastEvent, &applyLock, true)

	return
}

func runLocalTerraformContainer(path string, destroy bool, ts time.Time, lastEvent *lastEvent, applyLock *applyLock, skipWatch bool) {
	// postpone executing slightly
	time.Sleep(200 * time.Millisecond)

	// check if while we were sleeping another fs event queued an apply
	if ts != lastEvent.Get() {
		// cancel apply
		return
	}

	// even if we're the latest queued apply
	// we need to wait for a potential previous apply to finish
	applyLock.mux.Lock()
	defer applyLock.mux.Unlock()

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}

	// get current user id to set chown during docker build
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	imageTag := util.DockerImageTag(absPath, "loc")

	// build the docker image for this apply run
	buildArgs := []string{
		"--file", "Dockerfile.loc",
		"--tag", imageTag,
		"--build-arg", fmt.Sprintf("UID=%s", u.Uid),
		"--build-arg", fmt.Sprintf("GID=%s", u.Gid),
		"."}

	buildCmd := util.DockerBuildCommand(absPath, buildArgs)
	err = buildCmd.Run()
	if err != nil {
		log.Fatalf("docker build error: %s", err)
	}

	// parse the Terraform config
	module, _ := tfconfig.LoadModule(filepath.Dir(path))

	// prepare list of all module sources that need to be rewritten
	sedArgs := []string{}
	for _, value := range module.ModuleCalls {
		o := strings.Replace(value.Source, "/", "\\/", -1)
		r := strings.Replace(value.Source, "/cluster?", "/cluster-local?", 1)
		r = strings.Replace(r, "/", "\\/", -1)
		arg := fmt.Sprintf("-e s#%s#%s#g", o, r)
		sedArgs = append(sedArgs, arg)
	}

	// prepare volumes
	tfStatePathHash := util.PathHash(absPath)
	tfStatePath := "/infra/terraform.tfstate.d"
	tfStateVolume := fmt.Sprintf("kbst-loc-terraform-state-%s:%s", tfStatePathHash, tfStatePath)
	dockerSocketVolume := "/var/run/docker.sock:/var/run/docker.sock"

	tfCommand := "apply"
	if destroy {
		tfCommand = "destroy"
	}
	applySh := fmt.Sprintf(`
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

	runArgs := []string{
		"--rm",
		"--privileged",
		"--volume", tfStateVolume,
		"--volume", dockerSocketVolume,
		"--net", "host",
		imageTag,
		"sh", "-c", applySh}

	runCmd := util.DockerRunCommand(runArgs)
	err = runCmd.Run()
	if err != nil {
		log.Printf("docker run error: %s", err)
	}

	if skipWatch == false {
		log.Println("#### Watching for changes")
	}
	return
}
