package container

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/otiai10/copy"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type Container struct {
	cmd  *exec.Cmd
	r, w *os.File
}

func New(from string) *Container {
	var lowerdir, _ = os.ReadFile(from + "/lowerdir")
	var upperdir, _ = os.MkdirTemp("./", "*")
	var workdir, _ = os.MkdirTemp("./", "*")
	var target, _ = os.MkdirTemp("./", "*")
	var r, cw, _ = os.Pipe()
	var cr, w, _ = os.Pipe()
	var cmd = exec.CommandContext(context.Background(), "container", string(lowerdir), upperdir, workdir, target)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.ExtraFiles = []*os.File{cr, cw}
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET}
	return &Container{cmd, r, w}
}
func (c *Container) Start() chan error {
	var err = make(chan error)
	if err := c.cmd.Start(); err != nil {
		panic(err)
	}
	go func() {
		err <- c.cmd.Wait()
	}()
	return err
}
func (c *Container) Connect(name string) netlink.Link {
	var link = &netlink.Veth{LinkAttrs: netlink.LinkAttrs{Name: "eth0"}, PeerName: name}
	netlink.LinkAdd(link)
	netlink.LinkSetUp(link)
	var ns, _ = netns.GetFromPid(c.cmd.Process.Pid)
	defer ns.Close()
	netlink.LinkSetNsFd(link, int(ns))
	return link
}
func (c *Container) Copy(from string, to string) {
	copy.Copy(from, filepath.Join(c.cmd.Args[2], to), copy.Options{AddPermission: 0777})
}
func (c *Container) Exec(cmd string) {
	c.w.WriteString(cmd + "\n")
	c.r.Read(make([]byte, 1))
}
func (c *Container) Commit(to string) {
	var layers = []string{}
	for _, from := range append(strings.Split(c.cmd.Args[1], ":"), c.cmd.Args[2]) {
		var to = filepath.Join(to, "layers", filepath.Base(from))
		copy.Copy(from, to, copy.Options{AddPermission: 0777})
		layers = append(layers, to)
	}
	os.WriteFile(to+"/lowerdir", []byte(strings.Join(layers, ":")), 0777)
}
func (c *Container) Stop() {
	c.cmd.Cancel()
	c.r.Close()
	c.w.Close()
	os.RemoveAll(c.cmd.Args[2])
	os.RemoveAll(c.cmd.Args[3])
	os.RemoveAll(c.cmd.Args[4])
}
