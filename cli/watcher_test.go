package cli

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRepoWatcher(t *testing.T) {
	rw := RepoWatcher{}

	assert.IsType(t, &lastEvent{}, rw.le, nil)
	assert.IsType(t, &applyLock{}, rw.al, nil)
}

type MockRepositoryWatcher struct {
	throw bool
}

func (mrw MockRepositoryWatcher) Start(path string) (run chan time.Time) {
	if mrw.throw {
		log.Fatalf("mock error")
	}

	return run
}
