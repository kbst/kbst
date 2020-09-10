package cli

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func Cluster(argv []string) (err error) {
	usage := `
Usage:
  kbst cluster <name> add (aks | eks | gke | kind)

Options:
  -h, --help	Show this help.
`

	args, _ := docopt.ParseDoc(usage)
	fmt.Println(args)
	return
}
