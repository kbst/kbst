module github.com/kbst/kbst

go 1.13

replace github.com/kbst/kbst/cli => ./cli

replace github.com/kbst/kbst/util => ./util

require (
	github.com/adrg/xdg v0.2.1
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/hashicorp/terraform-config-inspect v0.0.0-20200806211835-c481b8bfa41e
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/spf13/cobra v1.0.0
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
	gopkg.in/fsnotify.v1 v1.4.7
)
