package aks

import (
	"github.com/docopt/docopt-go"
)

func Parse() docopt.Opts {
	usage := `Add AKS cluster
	
	Usage:
	kbst cluster <name> add aks
	kbst cluster <name> add aks -h | --help
	
	Options:
	-h --help     Show this screen.`

	arguments, _ := docopt.ParseDoc(usage)
	return arguments
}
