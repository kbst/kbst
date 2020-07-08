package cli

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
)

func Shell(argv []string) (err error) {
	usage := `
Usage:
  kbst shell [--path=path]

Options:
  -p, --path=path  Path to initialize the repository in [default: .].
  -h, --help	   Show this help.
`

	args, _ := docopt.ParseDoc(usage)
	log.Println(args)

	if args["shell"] == true {
		getCmd()
	}

	return
}

func getCmd() {
	fmt.Printf(
		"# to execute the commands use\n# eval \"$(kbst shell)\" \n\n%s\n%s\n",
		getBuildCmd(),
		getRunCmd(),
	)
	return
}

func getBuildCmd() string {
	return "docker build -t kbst:shell ."
}

func getRunCmd() string {
	return "docker run --rm -ti -v `pwd`:/infra -v /var/run/docker.sock:/var/run/docker.sock kbst:shell"
}
