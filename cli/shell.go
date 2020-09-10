package cli

import (
	"fmt"
)

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
