package stack

import (
	"os"
	"path"
)

var cwd, _ = os.Getwd()
var fixturesPath = path.Join(cwd, "test_fixtures")
