package stack

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"

	"golang.org/x/exp/maps"
)

func dockerfile(in []byte, clusters []Cluster) []byte {
	cps := make(map[string]bool)
	for _, c := range clusters {
		cps[c.cloudK8sPrefix()] = true
	}
	opts := maps.Keys(cps)

	re := regexp.MustCompile(`FROM (.*):(v\d+\.\d+\.\d+(?:-\w*\.\d+)?)(?:-)?(aks|eks|gke)?`)
	sc := bufio.NewScanner(bytes.NewReader(in))

	out := []byte{}
	for sc.Scan() {
		line := sc.Text()

		// FROM kubestack/framework:v0.18.0-beta.0 => ["kubestack/framework", "v0.18.0-beta.0"]
		// FROM kubestack/framework:v0.18.0-beta.0-eks => ["kubestack/framework", "v0.18.0-beta.0", "eks"]
		sm := re.FindStringSubmatch(line)

		// if the regex does not match
		// we leave the line unchanged
		if sm != nil {
			if len(opts) == 1 && (sm[3] == "" || sm[3] != opts[0]) {
				// current FROM __IS__ multi-cloud
				// currently defined clusters __ARE NOT__ multi-cloud
				// switch to single-provider image
				line = fmt.Sprintf("FROM %s:%s-%s", sm[1], sm[2], opts[0])
			}

			if len(opts) > 1 && sm[3] != "" {
				// current FROM __IS NOT__ multi-cloud
				// currently defined clusters __ARE__ multi-cloud
				// switch to multi-cloud image
				line = fmt.Sprintf("FROM %s:%s", sm[1], sm[2])
			}
		}

		// include every line changed or unchanged
		out = append(out, []byte(fmt.Sprintln(line))...)
	}

	return out
}
