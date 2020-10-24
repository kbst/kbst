package cli

import (
	"fmt"
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
	Start(path string) (chan fsnotify.Event, error)
	Stop()
}

type repoWatcher struct {
	e  chan fsnotify.Event
	le *lastEvent
	al *applyLock
	w  *fsnotify.Watcher
}

func NewRepoWatcher() *repoWatcher {
	rw := repoWatcher{
		e:  make(chan fsnotify.Event),
		le: &lastEvent{},
		al: &applyLock{},
	}

	return &rw
}

func (rw *repoWatcher) Start(path string) (chan fsnotify.Event, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return rw.e, fmt.Errorf("watching filesystem failed: %s", err)
	}
	rw.w = watcher

	go rw.handleEvent()

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
			return rw.e, fmt.Errorf("watching '%s' failed: %s", fullPath, err)
		}
	}

	return rw.e, nil
}

func (rw *repoWatcher) Stop() {
	rw.w.Close()
}

func (rw *repoWatcher) handleEvent() {
	for {
		select {
		case e, ok := <-rw.w.Events:
			if !ok {
				return
			}

			ts := time.Now()
			rw.le.Set(ts)
			go rw.queueRun(ts, e)
		case err, ok := <-rw.w.Errors:
			if !ok {
				log.Printf("error watching for changes: %s", err)
				return
			}
		}
	}
}

func (rw *repoWatcher) queueRun(ts time.Time, e fsnotify.Event) {
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

	rw.e <- e
}
