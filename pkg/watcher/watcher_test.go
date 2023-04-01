package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/fsnotify.v1"
)

var cwd, _ = os.Getwd()
var fixturesPath = filepath.Join(cwd, "../", "../", "test_fixtures")

func TestLastEvent(t *testing.T) {
	le := lastEvent{}
	ts := time.Now()
	le.Set(ts)
	rts := le.Get()

	assert.Equal(t, ts, rts, nil)
}

func TestNewRepoWatcher(t *testing.T) {
	rw := NewRepoWatcher()

	assert.IsType(t, make(chan fsnotify.Event), rw.e, nil)
	assert.IsType(t, &lastEvent{}, rw.le, nil)
	assert.IsType(t, &applyLock{}, rw.al, nil)
}

func TestRepoWatcher(t *testing.T) {
	p := filepath.Join(fixturesPath, "multi-cloud")

	rw := NewRepoWatcher()
	rw.Start(p)
	defer rw.Stop()

	// change a file
	fp := filepath.Join(p, "test")
	file, err := os.Create(fp)
	if err != nil {
		t.Error(err)
	}
	file.Close()
	err = os.Remove(fp)
	if err != nil {
		t.Error(err)
	}

	e := <-rw.e
	assert.Equal(t, fp, e.Name, nil)
}

func TestRepoWatcherQueueTwoEvents(t *testing.T) {
	p := filepath.Join(fixturesPath, "multi-cloud")

	rw := NewRepoWatcher()
	rw.Start(p)
	defer rw.Stop()

	// make the first change
	fp := filepath.Join(p, "test")
	file, err := os.Create(fp)
	if err != nil {
		t.Error(err)
	}
	file.Close()
	err = os.Remove(fp)
	if err != nil {
		t.Error(err)
	}

	// make the second change
	fp = filepath.Join(p, "test2")
	file, err = os.Create(fp)
	if err != nil {
		t.Error(err)
	}
	file.Close()
	err = os.Remove(fp)
	if err != nil {
		t.Error(err)
	}

	e := <-rw.e
	assert.Equal(t, fp, e.Name, nil)
}

func TestRepoWatcherPathError(t *testing.T) {
	p := filepath.Join(fixturesPath, "may-not_exist")

	rw := NewRepoWatcher()
	rw.Start(p)
	defer rw.Stop()
}
