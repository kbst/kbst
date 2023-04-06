package util

import (
	"os"
	"path/filepath"
)

var cwd, _ = os.Getwd()
var fixturesPath = filepath.Join(cwd, "../", "../", "test_fixtures")
