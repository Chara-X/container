module github.com/Chara-X/container

go 1.22.2

replace github.com/Chara-X/util => ../util

replace github.com/Chara-X/netns => ../netns

require (
	github.com/Chara-X/netns v0.0.0-00010101000000-000000000000
	github.com/Chara-X/util v0.0.0-00010101000000-000000000000
	github.com/vishvananda/netlink v1.1.0
)

require golang.org/x/sync v0.3.0 // indirect

require (
	github.com/otiai10/copy v1.14.0
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
