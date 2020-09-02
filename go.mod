module github.com/kbst/kbst

go 1.13

replace github.com/kbst/kbst/cli => ./cli

replace github.com/kbst/kbst/util => ./util

require (
	github.com/adrg/xdg v0.2.1
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/google/btree v1.0.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/kbst/kbst/cli v0.0.0-00010101000000-000000000000
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
)
