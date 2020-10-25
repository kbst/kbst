package cli

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/kbst/kbst/pkg/watcher"
	"github.com/stretchr/testify/assert"
)

type MockTerraformContainer struct {
	runCount    int
	runCountMux sync.Mutex
	throw       bool
}

func (mtc *MockTerraformContainer) Run() (err error) {
	mtc.runCount++

	if mtc.throw {
		return errors.New("mock error")
	}

	return nil
}

func (mtc *MockTerraformContainer) Count() {
	return
}

func TestLocalApply(t *testing.T) {
	mtc := &MockTerraformContainer{}
	rw := watcher.NewRepoWatcher()

	local := Local{Runner: mtc, Watcher: rw}
	p := filepath.Join(fixturesPath, "multi-cloud")

	// start a watch
	go local.Apply(p, false)

	assert.Eventually(t, func() bool { return mtc.runCount == 1 }, 500*time.Millisecond, 50*time.Millisecond, "expected mtc.runCount == 1")
}

func TestLocalApplyProvisionError(t *testing.T) {
	mtc := &MockTerraformContainer{throw: true}

	local := Local{Runner: mtc}
	p := filepath.Join(fixturesPath, "multi-cloud")
	err := local.Apply(p, false)

	assert.Error(t, err, nil)
}

func TestLocalApplySkipWatch(t *testing.T) {
	mtc := &MockTerraformContainer{}

	local := Local{Runner: mtc}
	p := filepath.Join(fixturesPath, "multi-cloud")
	err := local.Apply(p, true)

	assert.Equal(t, nil, err, nil)
	assert.Equal(t, 1, mtc.runCount, nil)
}

func TestLocalDestroy(t *testing.T) {
	mtc := &MockTerraformContainer{}

	local := Local{Runner: mtc}
	err := local.Destroy()

	assert.Equal(t, nil, err, nil)
}

func TestLocalDestroyError(t *testing.T) {
	mtc := &MockTerraformContainer{throw: true}

	local := Local{Runner: mtc}
	err := local.Destroy()

	assert.Error(t, err, nil)
}
