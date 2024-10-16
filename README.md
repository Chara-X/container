# Container runtime

```go
func ExampleContainer() {
	// container.Install("./ubuntu.tar.gz", "./ubuntu")
	var con = container.New("./ubuntu")
	defer con.Kill()
	// var l = con.Connect("veth0")
	// Config conn ...
	con.Copy("./bag", "/etc/bag")
	con.Exec()
	con.Commit("./mc")
}
```

## Reference

[Containers the hard way: Gocker: A mini Docker written in Go](https://unixism.net/2020/06/containers-the-hard-way-gocker-a-mini-docker-written-in-go)

[Digging into Linux namespaces - part 1](https://blog.quarkslab.com/digging-into-linux-namespaces-part-1.html)

[Understanding Container Images, Part 3: Working with Overlays](https://blogs.cisco.com/developer/373-containerimages-03)

[How Container Networking Works: a Docker Bridge Network From Scratch](https://labs.iximiuz.com/tutorials/container-networking-from-scratch)
