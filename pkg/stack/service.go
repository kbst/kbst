package stack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
)

type Service struct {
	EntryName      string
	ClusterName    string
	Provider       string
	Version        string
	Configurations []Configuration
}

func (s *Service) ToHCL() map[string]*hclwrite.File {
	files := make(map[string]*hclwrite.File)

	f := hclwrite.NewEmptyFile()

	source := fmt.Sprintf("kbst.xyz/catalog/%s/%s", s.EntryName, s.Provider)
	version := strings.TrimPrefix(s.Version, "v")

	tfhcl.ModuleService(f, s.Name(), s.ClusterName, source, version, convertToTfhclConfiguration(s.Configurations))
	files[fmt.Sprintf("%s.tf", s.Name())] = f

	return files
}

func (s *Service) Name() string {
	return fmt.Sprintf("%s_%s_%s", s.ClusterName, "service", s.EntryName)
}
