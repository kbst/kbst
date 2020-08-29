package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"gopkg.in/fsnotify.v1"
)

func Local(argv []string) (err error) {
	usage := `
Usage:
  kbst local dev [--path=path]

Options:
  -p, --path=path  Path to initialize the repository in [default: .].
  -h, --help	   Show this help.
`

	args, _ := docopt.ParseDoc(usage)
	fmt.Println(args)

	if args["dev"] == true {
		repoWatch(args["--path"].(string))
	}

	return
}

type ApplyLock struct {
	mux sync.Mutex
}

type LastEvent struct {
	ts  time.Time
	mux sync.Mutex
}

func (l *LastEvent) Set(ts time.Time) {
	l.mux.Lock()
	l.ts = ts
	l.mux.Unlock()
}

func (l *LastEvent) Get() time.Time {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.ts
}

func repoWatch(path string) (err error) {
	applyLock := ApplyLock{}
	lastEvent := LastEvent{}

	// first apply to bring up dev env
	ts := time.Now()
	lastEvent.Set(ts)
	go handleChange(path, ts, &lastEvent, &applyLock)

	// then start watching
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					return
				}

				log.Println(event)

				ts := time.Now()
				lastEvent.Set(ts)
				go handleChange(path, ts, &lastEvent, &applyLock)
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

func handleChange(path string, ts time.Time, lastEvent *LastEvent, applyLock *ApplyLock) {
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

	// get current user id to set chown during docker build
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	imageTag := fmt.Sprintf("kbst:%d", ts.Unix())

	// build the docker image for this apply run
	buildCmd := exec.Command(
		"docker",
		"build",
		"--progress", "plain",
		"--file", "Dockerfile.loc",
		"--tag", imageTag,
		"--build-arg", fmt.Sprintf("UID=%s", u.Uid),
		"--build-arg", fmt.Sprintf("GID=%s", u.Gid),
		".")
	buildCmd.Env = []string{"DOCKER_BUILDKIT=1"}
	buildCmd.Dir = path
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout

	err = buildCmd.Run()
	if err != nil {
		log.Fatalln(err)
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
	tfStatePath := "/infra/terraform.tfstate.d"
	tfStateVolume := fmt.Sprintf("kbst-loc-terraform-state:%s", tfStatePath)
	dockerSocketVolume := "/var/run/docker.sock:/var/run/docker.sock"

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
	
	terraform destroy --auto-approve
	`, strings.Join(sedArgs, " "))
	runCmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--privileged",
		"--volume", tfStateVolume,
		"--volume", dockerSocketVolume,
		imageTag,
		"sh", "-c", applySh)
	runCmd.Stderr = os.Stderr
	runCmd.Stdout = os.Stdout

	err = runCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	return
}
