package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	var lowerdir, upperdir, workdir, target = os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	for _, v := range [][4]string{{"", "", "overlay", "lowerdir=" + lowerdir + ",upperdir=" + upperdir + ",workdir=" + workdir}, {"", "/proc", "proc", ""}, {"", "/tmp", "tmpfs", ""}, {"", "/dev", "tmpfs", ""}, {"", "/dev/pts", "devpts", ""}, {"", "/sys", "sysfs", ""}} {
		os.MkdirAll(target+v[1], 0777)
		if err := syscall.Mount(v[0], target+v[1], v[2], 0, v[3]); err != nil {
			panic(err)
		}
	}
	syscall.Chroot(target)
	os.Chdir("/")
	var cmd = exec.Command("sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Run()
}
