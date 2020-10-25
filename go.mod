module github.com/kbst/kbst

go 1.13

replace github.com/kbst/kbst/cli => ./cli

replace github.com/kbst/kbst/util => ./util

replace github.com/kbst/kbst/pkg => ./pkg

require (
	github.com/adrg/xdg v0.2.1
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/hashicorp/terraform-config-inspect v0.0.0-20200806211835-c481b8bfa41e
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/mod v0.0.0-20190513183733-4bf6d317e70e
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/yaml.v2 v2.3.0
	sigs.k8s.io/kustomize/api v0.6.0
	sigs.k8s.io/yaml v1.2.0
)
