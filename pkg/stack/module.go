package stack

import (
	"github.com/kbst/kbst/pkg/tfhcl"
)

type Module struct {
	tfMod *tfhcl.Module
}

func (c *Module) Name() string {
	return c.tfMod.Name
}
