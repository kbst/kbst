package cli

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepoWatcher(t *testing.T) {
	tc := MockTerraformContainer{}
	rw := NewRepoWatcher(&tc)

	assert.IsType(t, &lastEvent{}, rw.le, nil)
	assert.IsType(t, &applyLock{}, rw.al, nil)
}

type MockRepositoryWatcher struct {
	throw bool
}

func (mrw MockRepositoryWatcher) Start(path string) {
	if mrw.throw {
		log.Fatalf("mock error")
	}

	return
}
