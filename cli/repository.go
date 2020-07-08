package cli

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func Repository(argv []string) (err error) {
	usage := `
Usage:
  kbst repository init [--path=path]

Options:
  -p, --path=path  Path to initialize the repository in [default: .].
  -h, --help	   Show this help.
`

	args, _ := docopt.ParseDoc(usage)
	fmt.Println(args)
	return
}
