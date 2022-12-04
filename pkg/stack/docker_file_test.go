package stack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerfileUnchanged(t *testing.T) {
	in := []byte(`FROM aline:latest

ENV test=test\n
	test2=test2

RUN echo "Hello" && \n
	echo "World"
`)

	cls := []Cluster{}

	out := dockerfile(in, cls)

	assert.Equal(t, string(in), string(out))
}

func TestDockerfileAlwaysEndsInEmptyLine(t *testing.T) {
	in := []byte(`FROM aline:latest`)

	cls := []Cluster{}

	out := dockerfile(in, cls)

	assert.Equal(t, fmt.Sprintf("%s\n", in), string(out))
}

func TestDockerfileSingleToMulti(t *testing.T) {
	cls := []Cluster{
		{
			Provider: "azurerm",
		},
		{
			Provider: "google",
		},
	}

	// aks
	in := []byte("FROM kubestack/framework:v0.18.0-beta.0-aks")
	exp := []byte("FROM kubestack/framework:v0.18.0-beta.0\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))

	// eks
	in = []byte("FROM kubestack/framework:v0.18.0-beta.0-eks")
	exp = []byte("FROM kubestack/framework:v0.18.0-beta.0\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))

	// gke
	in = []byte("FROM kubestack/framework:v0.18.0-beta.0-gke")
	exp = []byte("FROM kubestack/framework:v0.18.0-beta.0\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))
}

func TestDockerfileMultiToSingleAKS(t *testing.T) {
	cls := []Cluster{
		{
			Provider: "azurerm",
		},
	}

	// aks
	in := []byte("FROM kubestack/framework:v0.18.0-beta.0")
	exp := []byte("FROM kubestack/framework:v0.18.0-beta.0-aks\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))
}

func TestDockerfileMultiToSingleEKS(t *testing.T) {
	cls := []Cluster{
		{
			Provider: "aws",
		},
	}

	// aks
	in := []byte("FROM kubestack/framework:v0.18.0-beta.0")
	exp := []byte("FROM kubestack/framework:v0.18.0-beta.0-eks\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))
}

func TestDockerfileMultiToSingleGKE(t *testing.T) {
	cls := []Cluster{
		{
			Provider: "google",
		},
	}

	// aks
	in := []byte("FROM kubestack/framework:v0.18.0-beta.0")
	exp := []byte("FROM kubestack/framework:v0.18.0-beta.0-gke\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))
}

func TestDockerfileCorrectSingle(t *testing.T) {
	cls := []Cluster{
		{
			Provider: "google",
		},
	}

	in := []byte("FROM kubestack/framework:v0.18.0-beta.0-INCORRECT")
	exp := []byte("FROM kubestack/framework:v0.18.0-beta.0-gke\n")

	assert.Equal(t, string(exp), string(dockerfile(in, cls)))
}
