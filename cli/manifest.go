package cli

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func Manifest(argv []string) (err error) {
	usage := `
Usage:
  kbst manifest install <url> [--base-path=path] [--vendor-only]
  kbst manifest update <url> [--base-path=path]
  kbst manifest remove <name> [--base-path=path]

Options:
  -b=path, --base-path=path  Overwrite the base directory path [default: ./manifests/bases]. 
  -h, --help	             Show this help.
`

	args, _ := docopt.ParseDoc(usage)
	fmt.Println(args)
	return
}
