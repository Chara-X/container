module github.com/Chara-X/container

go 1.23.1

replace github.com/Chara-X/util => ../util

require (
	github.com/Chara-X/util v0.0.0-00010101000000-000000000000
	github.com/creack/pty v1.1.23
	github.com/vishvananda/netlink v1.1.0
)

require golang.org/x/sync v0.3.0 // indirect

require (
	github.com/otiai10/copy v1.14.0
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
