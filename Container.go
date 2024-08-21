package container

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/otiai10/copy"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type Container struct {
	*os.Process
	Pty  *os.File
	Args []string
}

func New(from string) *Container {
	var lowerdir, _ = os.ReadFile(from + "/lowerdir")
	var upperdir, _ = os.MkdirTemp("./", "*")
	var workdir, _ = os.MkdirTemp("./", "*")
	var target, _ = os.MkdirTemp("./", "*")
	var pty, tty, _ = pty.Open()
	var path, _ = exec.LookPath("container")
	var args = []string{path, string(lowerdir), upperdir, workdir, target}
	var proc, _ = os.StartProcess(path, args, &os.ProcAttr{
		Files: []*os.File{tty, tty, tty},
		Sys: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
			Setctty:    true, Setsid: true, Ctty: 0,
		},
	})
	go func() {
		proc.Wait()
		pty.Close()
		os.RemoveAll(upperdir)
		os.RemoveAll(workdir)
		os.RemoveAll(target)
	}()
	return &Container{proc, pty, args}
}
func (c *Container) Connect(name string) netlink.Link {
	var link = &netlink.Veth{LinkAttrs: netlink.LinkAttrs{Name: "eth0"}, PeerName: name}
	netlink.LinkAdd(link)
	netlink.LinkSetUp(link)
	var ns, _ = netns.GetFromPid(c.Pid)
	defer ns.Close()
	netlink.LinkSetNsFd(link, int(ns))
	return link
}
func (c *Container) Copy(from string, to string) {
	copy.Copy(from, filepath.Join(c.Args[2], to), copy.Options{AddPermission: 0777})
}
func (c *Container) Exec() {
	go io.Copy(os.Stdout, c.Pty)
	io.Copy(c.Pty, os.Stdin)
}
func (c *Container) Commit(to string) {
	var layers = []string{}
	for _, from := range append(strings.Split(c.Args[1], ":"), c.Args[2]) {
		var to = filepath.Join(to, "layers", filepath.Base(from))
		copy.Copy(from, to, copy.Options{AddPermission: 0777})
		layers = append(layers, to)
	}
	os.WriteFile(to+"/lowerdir", []byte(strings.Join(layers, ":")), 0777)
}
