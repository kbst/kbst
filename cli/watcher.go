package cli

import (
	"log"
	"path/filepath"
	"sync"
	"time"

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

type Watcher interface {
	Start(path string)
}

type RepoWatcher struct {
	tc TerraformContainer
	le *lastEvent
	al *applyLock
	w  *fsnotify.Watcher
}

func NewRepoWatcher(tc TerraformContainer) (rw RepoWatcher) {
	rw.tc = tc
	rw.le = &lastEvent{}
	rw.al = &applyLock{}

	return rw
}

func (rw RepoWatcher) Start(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("watching filesystem failed: %s", err)
	}
	defer watcher.Close()
	rw.w = watcher

	done := make(chan bool)
	go rw.handleEvent(done)

	watchTargets := []string{
		".",
		"manifests/bases",
		"manifests/overlays/apps",
		"manifests/overlays/ops",
		"manifests/overlays/loc",
	}
	for _, t := range watchTargets {
		fullPath := filepath.Join(path, t)
		err = rw.w.Add(fullPath)
		if err != nil {
			log.Fatalf("watching '%s' failed: %s", fullPath, err)
		}
	}

	<-done
	return
}

func (rw RepoWatcher) handleEvent(done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			select {
			case _, ok := <-rw.w.Events:
				if !ok {
					return
				}

				ts := time.Now()
				rw.le.Set(ts)
				go rw.queueRun(ts)
			case err, ok := <-rw.w.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}
}

func (rw RepoWatcher) queueRun(ts time.Time) {
	// postpone run slightly
	time.Sleep(200 * time.Millisecond)

	// check if while we were sleeping another fs event queued an apply
	if ts != rw.le.Get() {
		// cancel apply
		return
	}

	// even if we're the latest queued apply
	// we need to wait for a potential previous apply to finish
	rw.al.mux.Lock()
	defer rw.al.mux.Unlock()

	err := rw.tc.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("#### Watching for changes")
}
