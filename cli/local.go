package cli

import (
	"errors"
	"fmt"
	"log"

	"github.com/kbst/kbst/pkg/runner"
	"github.com/kbst/kbst/pkg/watcher"
)

type Local struct {
	Runner  runner.TerraformContainer
	Watcher watcher.Watcher
}

func (l *Local) Apply(path string, skipWatch bool) (err error) {
	// provision the development environment
	err = l.Runner.Run(false)
	if err != nil {
		return errors.New(fmt.Sprintf("provisioning local environment error: %s", err))
	}

	if skipWatch {
		return
	}

	// start watching for repository changes
	ev, err := l.Watcher.Start(path)
	if err != nil {
		return errors.New(fmt.Sprintf("error watching repository: %s", err))
	}
	defer l.Watcher.Stop()

	for {
		log.Println("#### Watching for changes")
		<-ev
		err = l.Runner.Run(false)
		if err != nil {
			return errors.New(fmt.Sprintf("updating local environment error: %s", err))
		}
	}
}

func (l *Local) Destroy() (err error) {
	err = l.Runner.Run(true)
	if err != nil {
		return err
	}

	return
}
