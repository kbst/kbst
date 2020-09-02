package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/kbst/kbst/cli"
	"github.com/kbst/kbst/util"
)

func main() {
	usage := `Kubestack CLI

Usage:
  kbst <resource> [<arguments>...]

Available resources:
  repository
  cluster
  manifest
  shell
  help

Options:
  -h, --help   	 Show this help.
  -v, --version	 Show version.
`

	log.SetOutput(os.Stderr)
	log.SetPrefix("")
	log.SetFlags(0)

	parser := &docopt.Parser{OptionsFirst: true}
	args, _ := parser.ParseArgs(usage, nil, "kbst version v0.0.0")

	log.Printf("args: %s", args)

	resource := args["<resource>"].(string)
	arguments := args["<arguments>"].([]string)

	// If err, we can't show update notification
	cli, err := util.GetCli()
	if err == nil {
		// TODO: replace with version check and update notification
		for v := range cli.Versions {
			log.Println(v)
		}
	}

	err = switchKind(resource, arguments)
	if err != nil {
		log.Fatalln(err)
	}
}

func switchKind(res string, args []string) (err error) {
	argv := append([]string{res}, args...)
	log.Printf("argv: %s", argv)

	switch res {
	case "cluster":
		return cli.Cluster(argv)
	case "local":
		return cli.Local(argv)
	case "manifest":
		return cli.Manifest(argv)
	case "repository":
		return cli.Repository(argv)
	case "shell":
		return cli.Shell(argv)
	case "help":
		return switchKind(args[0], append(args[1:], "--help"))
	}

	return fmt.Errorf("'%s' is not a valid kbst resource. See 'kbst --help'", res)
}
