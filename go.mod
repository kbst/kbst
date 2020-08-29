module github.com/kbst/kbst

go 1.13

replace github.com/kbst/kbst/cli => ./cli

require (
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20200806211835-c481b8bfa41e // indirect
	github.com/kbst/kbst/cli v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
)
