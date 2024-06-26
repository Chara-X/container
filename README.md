# Container runtime

```go
func ExampleContainer() {
	// container.Install("./ubuntu.tar")
	var con = container.New("ubuntu")
	defer con.Stop()
	con.Start()
	con.Connect("veth0", "172.18.0.20/16")
	con.Exec("sh")
	con.Commit("mc")
}
```
