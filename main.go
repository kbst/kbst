package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/kbst/kbst/cli"
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

	err := switchKind(resource, arguments)
	if err != nil {
		log.Fatalln(err)
	}
}

func switchKind(res string, args []string) (err error) {
	argv := append([]string{res}, args...)
	log.Printf("argv: %s", argv)

	switch res {
	case "repository":
		return cli.Repository(argv)
	case "cluster":
		return cli.Cluster(argv)
	case "manifest":
		return cli.Manifest(argv)
	case "shell":
		return cli.Shell(argv)
	case "help":
		return switchKind(args[0], append(args[1:], "--help"))
	}

	return fmt.Errorf("'%s' is not a valid kbst resource. See 'kbst --help'", res)
}
