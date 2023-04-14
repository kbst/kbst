package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/kbst/kbst/pkg/export"
	"github.com/kbst/kbst/pkg/stack"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

type Repo struct {
	Catalog    map[string]util.Entry
	Framework  util.Entry
	Downloader util.Downloader
}

func (r Repo) Init(starter string, baseDomain string, namePrefix string, region string, envNames []string, baseCfg map[string]cty.Value, release string, gitRef string, path string) (err error) {
	var environments []stack.Environment

	for _, en := range envNames {
		isBaseKey := false
		if en == envNames[0] {
			isBaseKey = true
		}

		environments = append(environments, stack.Environment{
			Key:       en,
			IsBaseKey: isBaseKey,
		})
	}

	cj := util.CliJSON{}
	err = cj.Load(util.CachedDownloader{})
	if err != nil {
		log.Fatal(err)
	}

	filenames, err := r.extractArchive(starter, release, gitRef, path)
	if err != nil {
		return err
	}

	// determine unzip target path
	repoPath, err := filepath.Abs(filenames[0])
	if err != nil {
		return err
	}

	// remove .tf files included in starter archives
	err = r.cleanStarter(repoPath)
	if err != nil {
		return err
	}

	s := stack.NewStack(tfhcl.NewRoot(repoPath), cj)
	if err != nil {
		log.Fatal(err)
	}

	s.SetBaseDomain(cty.StringVal(baseDomain))
	s.Environments = environments

	switch starter {
	case "aks":
		_, err = s.AddCluster(namePrefix, "azurerm", region, "", stack.GenerateConfigurations(s.Environments, baseCfg))
	case "eks":
		_, err = s.AddCluster(namePrefix, "aws", region, "", stack.GenerateConfigurations(s.Environments, baseCfg))
	case "gke":
		_, err = s.AddCluster(namePrefix, "google", region, "", stack.GenerateConfigurations(s.Environments, baseCfg))
	default:
		return fmt.Errorf("unexpected error: starter: '%s' exists as archive, but is not implemented in CLI", starter)
	}
	if err != nil {
		return err
	}

	err = s.InitFiles(baseDomain)
	if err != nil {
		return err
	}

	// initialize git repository
	err = r.gitCommit(repoPath, fmt.Sprintf("Initialized from %s starter", strings.ToUpper(starter)))
	if err != nil {
		return err
	}

	return
}

func (r Repo) Import(es export.Stack, path string) (err error) {
	cj := util.CliJSON{}
	err = cj.Load(util.CachedDownloader{})
	if err != nil {
		log.Fatal(err)
	}

	providerToStarter := map[string]string{
		"aws":     "eks",
		"azurerm": "aks",
		"google":  "gke",
	}

	filenames, err := r.extractArchive(providerToStarter[es.Clusters[0].Provider], es.Clusters[0].Version, "", path)
	if err != nil {
		return err
	}

	// determine unzip target path
	repoPath, err := filepath.Abs(filenames[0])
	if err != nil {
		return err
	}

	// remove .tf files included in starter archives
	err = r.cleanStarter(repoPath)
	if err != nil {
		return err
	}

	s := stack.NewStack(tfhcl.NewRoot(repoPath), cj)
	if err != nil {
		log.Fatal(err)
	}

	s.SetBaseDomain(cty.StringVal(es.BaseDomain))

	s.Environments = es.StackEnvironments()

	for _, c := range es.Clusters {
		_, err = s.AddCluster(c.NamePrefix, c.Provider, c.Region, c.Version, es.StackConfigurations(c.Configurations))
		if err != nil {
			return err
		}
	}

	for _, np := range es.NodePools {
		_, err = s.AddNodePool(np.ClusterName, np.PoolName, es.StackConfigurations(np.Configurations))
		if err != nil {
			return err
		}
	}

	for _, svc := range es.Services {
		_, err = s.AddService(svc.ClusterName, svc.EntryName, svc.Version)
		if err != nil {
			return err
		}
	}

	err = s.InitFiles(es.BaseDomain)
	if err != nil {
		return err
	}

	// initialize git repository
	err = r.gitCommit(repoPath, fmt.Sprintf("Initialized from %s import", es.BaseDomain))
	if err != nil {
		return err
	}

	return
}

func (r Repo) downloadUrl(starter string, release string, gitRef string) (url string, err error) {
	if gitRef != "" {
		return fmt.Sprintf(
			"https://storage.googleapis.com/dev.quickstart.kubestack.com/kubestack-starter-%s-%s.zip",
			starter,
			gitRef,
		), nil
	}

	// determine version
	version, err := r.Framework.GetReleaseOrLatest(release)
	if err != nil {
		return url, err
	}

	url, ok := version.Archives[starter]
	if !ok {
		options := []string{}
		for k := range version.Archives {
			options = append(options, k)
		}
		sort.Strings(options)

		return url, fmt.Errorf(
			"'%s' is not a valid starter name, choose one of %v",
			starter,
			options,
		)
	}

	return url, nil
}

func (r Repo) extractArchive(starter string, release string, gitRef string, path string) (filenames []string, err error) {
	// download archive
	url, err := r.downloadUrl(starter, release, gitRef)
	if err != nil {
		return filenames, err
	}

	resp, err := r.Downloader.Download(url)
	if err != nil {
		return filenames, err
	}
	defer resp.Body.Close()

	// extract archive
	archive, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return filenames, err
	}

	filenames, err = util.Unzip(archive, path)
	if err != nil {
		return filenames, err
	}

	return filenames, nil
}

func (r Repo) cleanStarter(repoPath string) error {
	// replace .tf files in repoPath with generated files
	glob := filepath.Join(repoPath, "*.tf")
	contents, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Repo) gitCommit(p string, msg string) error {
	repo, err := git.PlainInit(p, false)
	if err != nil {
		return err
	}

	// make initial commit
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	copts := &git.CommitOptions{All: true}

	_, err = wt.Add(".")
	if err != nil {
		return err
	}

	_, err = wt.Commit(msg, copts)
	if err != nil {
		return err
	}

	return nil
}
