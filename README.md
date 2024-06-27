# Container runtime

```go
func ExampleContainer() {
	// container.Install("./ubuntu.tar")
	var con = container.New("./ubuntu")
	defer con.Stop()
	con.Start()
	// var l = con.Connect("veth0")
	// Config ...
	con.Copy("./bag", "/etc/bag")
	con.Exec("sh")
	con.Commit("./mc")
}
```
