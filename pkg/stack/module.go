package stack

import (
	"github.com/kbst/kbst/pkg/tfhcl"
)

type Module struct {
	mod *tfhcl.Module
}

func (c *Module) Name() string {
	return c.mod.Name
}
